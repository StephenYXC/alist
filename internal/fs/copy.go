package fs

import (
	"context"
	"fmt"
	stdpath "path"

	"github.com/alist-org/alist/v3/pkg/task"
	"github.com/alist-org/alist/v3/pkg/utils"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/operations"
	"github.com/pkg/errors"
)

var CopyTaskManager = task.NewTaskManager()

// Copy if in an account, call move method
// if not, add copy task
func Copy(ctx context.Context, account driver.Driver, srcPath, dstPath string) (bool, error) {
	srcAccount, srcActualPath, err := operations.GetAccountAndActualPath(srcPath)
	if err != nil {
		return false, errors.WithMessage(err, "failed get src account")
	}
	dstAccount, dstActualPath, err := operations.GetAccountAndActualPath(dstPath)
	if err != nil {
		return false, errors.WithMessage(err, "failed get dst account")
	}
	// copy if in an account, just call driver.Copy
	if srcAccount.GetAccount() == dstAccount.GetAccount() {
		return false, operations.Copy(ctx, account, srcActualPath, dstActualPath)
	}
	// not in an account
	CopyTaskManager.Submit(
		fmt.Sprintf("copy [%s](%s) to [%s](%s)", srcAccount.GetAccount().VirtualPath, srcActualPath, dstAccount.GetAccount().VirtualPath, dstActualPath),
		func(task *task.Task) error {
			return CopyBetween2Accounts(task, srcAccount, dstAccount, srcActualPath, dstActualPath)
		})
	return true, nil
}

func CopyBetween2Accounts(t *task.Task, srcAccount, dstAccount driver.Driver, srcPath, dstPath string) error {
	t.SetStatus("getting src object")
	srcObj, err := operations.Get(t.Ctx, srcAccount, srcPath)
	if err != nil {
		return errors.WithMessagef(err, "failed get src [%s] file", srcPath)
	}
	if srcObj.IsDir() {
		t.SetStatus("src object is dir, listing objs")
		objs, err := operations.List(t.Ctx, srcAccount, srcPath)
		if err != nil {
			return errors.WithMessagef(err, "failed list src [%s] objs", srcPath)
		}
		for _, obj := range objs {
			if utils.IsCanceled(t.Ctx) {
				return nil
			}
			srcObjPath := stdpath.Join(srcPath, obj.GetName())
			dstObjPath := stdpath.Join(dstPath, obj.GetName())
			CopyTaskManager.Submit(
				fmt.Sprintf("copy [%s](%s) to [%s](%s)", srcAccount.GetAccount().VirtualPath, srcObjPath, dstAccount.GetAccount().VirtualPath, dstObjPath),
				func(t *task.Task) error {
					return CopyBetween2Accounts(t, srcAccount, dstAccount, srcObjPath, dstObjPath)
				})
		}
	} else {
		CopyTaskManager.Submit(
			fmt.Sprintf("copy [%s](%s) to [%s](%s)", srcAccount.GetAccount().VirtualPath, srcPath, dstAccount.GetAccount().VirtualPath, dstPath),
			func(t *task.Task) error {
				return CopyFileBetween2Accounts(t, srcAccount, dstAccount, srcPath, dstPath)
			})
	}
	return nil
}

func CopyFileBetween2Accounts(t *task.Task, srcAccount, dstAccount driver.Driver, srcPath, dstPath string) error {
	srcFile, err := operations.Get(t.Ctx, srcAccount, srcPath)
	if err != nil {
		return errors.WithMessagef(err, "failed get src [%s] file", srcPath)
	}
	link, err := operations.Link(t.Ctx, srcAccount, srcPath, model.LinkArgs{})
	if err != nil {
		return errors.WithMessagef(err, "failed get [%s] link", srcPath)
	}
	stream, err := getFileStreamFromLink(srcFile, link)
	if err != nil {
		return errors.WithMessagef(err, "failed get [%s] stream", srcPath)
	}
	return operations.Put(t.Ctx, dstAccount, dstPath, stream, t.SetProgress)
}
