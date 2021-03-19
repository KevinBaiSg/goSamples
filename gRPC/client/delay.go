package main

import (
	"sync"
	"sync/atomic"
	"unsafe"

	pb "github.com/KevinBaiSg/goSamples/grpc/proto"
	"google.golang.org/grpc"
)

/*
	复用 client
*/
var (
	globalClientConn unsafe.Pointer
	lck              sync.Mutex
)

func GetClient(target string) (pb.RouteClient, error) {
	conn, err := GetConn(target)
	if err != nil {
		return (pb.RouteClient)(nil), err
	}
	return pb.NewRouteClient(conn), nil
}

func GetConn(target string) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConn) != nil { // double check
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	cli, err := newGRpcConn(target)
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConn, unsafe.Pointer(cli))
	return cli, nil
}

func newGRpcConn(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
