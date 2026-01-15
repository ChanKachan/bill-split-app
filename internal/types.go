package internal

import (
	postgre "bill-split/pkg/postgreWrapper"
	"context"
)

type UserInfo struct {
	UserId int
}

type CtxQuery struct {
	Ctx      context.Context
	Tx       postgre.TxPG
	UserInfo UserInfo
}
