package local

import (
	"context"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/operations"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/pkg/errors"
)

type Driver struct {
	model.Account
	Addition
}

func (d Driver) Config() driver.Config {
	return config
}

func (d *Driver) Init(ctx context.Context, account model.Account) error {
	d.Account = account
	err := utils.Json.UnmarshalFromString(d.Account.Addition, &d.Addition)
	if err != nil {
		return errors.Wrap(err, "error while unmarshal addition")
	}
	if !utils.Exists(d.RootFolder) {
		err = errors.Errorf("root folder %s not exists", d.RootFolder)
		d.SetStatus(err.Error())
	} else {
		d.SetStatus("OK")
	}
	operations.MustSaveDriverAccount(d)
	return err
}

func (d *Driver) Drop(ctx context.Context) error {
	return nil
}

func (d *Driver) GetAddition() driver.Additional {
	return d.Addition
}

func (d *Driver) List(ctx context.Context, dir model.Obj) ([]model.Obj, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Remove(ctx context.Context, obj model.Obj) error {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) Put(ctx context.Context, parentDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) Other(ctx context.Context, data interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

var _ driver.Driver = (*Driver)(nil)
