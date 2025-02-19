package bootstrap

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alist-org/alist/v3/cmd/args"
	"github.com/alist-org/alist/v3/conf"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

func InitConfig() {
	log.Infof("reading config file: %s", args.Config)
	if !utils.Exists(args.Config) {
		log.Infof("config file not exists, creating default config file")
		_, err := utils.CreateNestedFile(args.Config)
		if err != nil {
			log.Fatalf("failed to create config file")
		}
		conf.Conf = conf.DefaultConfig()
		if !utils.WriteToJson(args.Config, conf.Conf) {
			log.Fatalf("failed to create default config file")
		}
	} else {
		configBytes, err := ioutil.ReadFile(args.Config)
		if err != nil {
			log.Fatalf("reading config file error:%s", err.Error())
		}
		conf.Conf = conf.DefaultConfig()
		err = utils.Json.Unmarshal(configBytes, conf.Conf)
		if err != nil {
			log.Fatalf("load config error: %s", err.Error())
		}
		log.Debugf("config:%+v", conf.Conf)
		// update config.json struct
		confBody, err := utils.Json.MarshalIndent(conf.Conf, "", "  ")
		if err != nil {
			log.Fatalf("marshal config error:%s", err.Error())
		}
		err = ioutil.WriteFile(args.Config, confBody, 0777)
		if err != nil {
			log.Fatalf("update config struct error: %s", err.Error())
		}
	}
	if !conf.Conf.Force {
		confFromEnv()
	}
	err := os.RemoveAll(filepath.Join(conf.Conf.TempDir))
	if err != nil {
		log.Errorln("failed delete temp file:", err)
	}
	err = os.MkdirAll(conf.Conf.TempDir, 0700)
	if err != nil {
		log.Fatalf("create temp dir error: %s", err.Error())
	}
	log.Debugf("config: %+v", conf.Conf)
}

func confFromEnv() {
	prefix := "ALIST_"
	if args.NoPrefix {
		prefix = ""
	}
	log.Infof("load config from env with prefix: %s", prefix)
	if err := env.Parse(conf.Conf, env.Options{
		Prefix: prefix,
	}); err != nil {
		log.Fatalf("load config from env error: %s", err.Error())
	}
}
