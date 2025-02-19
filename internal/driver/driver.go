package driver

import (
	"context"

	"github.com/alist-org/alist/v3/internal/model"
)

type Driver interface {
	Meta
	Reader
	Writer
	Other
}

type Meta interface {
	Config() Config
	// Init If already initialized, drop first
	// need to unmarshal string to addition first
	Init(ctx context.Context, account model.Account) error
	Drop(ctx context.Context) error
	// GetAccount just get raw account
	GetAccount() model.Account
	GetAddition() Additional
}

type Other interface {
	Other(ctx context.Context, data interface{}) (interface{}, error)
}

type Reader interface {
	// List files in the path
	// if identify files by path, need to set ID with path,like path.Join(dir.GetID(), obj.GetName())
	// if identify files by id, need to set ID with corresponding id
	List(ctx context.Context, dir model.Obj) ([]model.Obj, error)
	// Link get url/filepath/reader of file
	Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error)
	//Get(ctx context.Context, path string) (FileInfo, error) // maybe not need
}

type Writer interface {
	// MakeDir make a folder named `dirName` in `parentDir`
	MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error
	// Move `srcObject` to `dstDir`
	Move(ctx context.Context, srcObj, dstDir model.Obj) error
	// Rename rename `srcObject` to `newName`
	Rename(ctx context.Context, srcObj model.Obj, newName string) error
	// Copy `srcObject` to `dstDir`
	Copy(ctx context.Context, srcObj, dstDir model.Obj) error
	// Remove remove `object`
	Remove(ctx context.Context, obj model.Obj) error
	// Put upload `stream` to `parentDir`
	Put(ctx context.Context, parentDir model.Obj, stream model.FileStreamer, up UpdateProgress) error
}

type UpdateProgress func(percentage int)
