package fs

import (
	"context"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/operations"
	"github.com/pkg/errors"
)

func MakeDir(ctx context.Context, account driver.Driver, path string) error {
	account, actualPath, err := operations.GetAccountAndActualPath(path)
	if err != nil {
		return errors.WithMessage(err, "failed get account")
	}
	return operations.MakeDir(ctx, account, actualPath)
}

func Move(ctx context.Context, account driver.Driver, srcPath, dstPath string) error {
	srcAccount, srcActualPath, err := operations.GetAccountAndActualPath(srcPath)
	if err != nil {
		return errors.WithMessage(err, "failed get src account")
	}
	dstAccount, dstActualPath, err := operations.GetAccountAndActualPath(srcPath)
	if err != nil {
		return errors.WithMessage(err, "failed get dst account")
	}
	if srcAccount.GetAccount() != dstAccount.GetAccount() {
		return errors.WithStack(ErrMoveBetwwenTwoAccounts)
	}
	return operations.Move(ctx, account, srcActualPath, dstActualPath)
}

func Rename(ctx context.Context, account driver.Driver, srcPath, dstName string) error {
	account, srcActualPath, err := operations.GetAccountAndActualPath(srcPath)
	if err != nil {
		return errors.WithMessage(err, "failed get account")
	}
	return operations.Rename(ctx, account, srcActualPath, dstName)
}

func Remove(ctx context.Context, account driver.Driver, path string) error {
	account, actualPath, err := operations.GetAccountAndActualPath(path)
	if err != nil {
		return errors.WithMessage(err, "failed get account")
	}
	return operations.Remove(ctx, account, actualPath)
}
