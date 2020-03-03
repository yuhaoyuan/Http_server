package rpc

import (
	"errors"
	"github.com/yuhaoyuan/RPC_server/corn"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

// Pool 基本方法
type Pool interface {
	Get() (*corn.Client, error)

	Put(*corn.Client) error

	Close(*corn.Client) error

	Release()

	Len() int
}
