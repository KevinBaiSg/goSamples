package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/KevinBaiSg/etcdSample/common"
	"github.com/coreos/etcd/clientv3"
	"log"
)

func main() {
	c, err := common.NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	kv := clientv3.NewKV(c)
	ctx := context.Background()

	response, err := kv.Put(ctx, "kevin", string(fromUInt64(1000)))
	if err != nil {
		log.Fatal("Error: Set Failed:", err)
	}
	log.Print("Success: Set Ok; Response:", response)

	response, err = kv.Put(ctx, "tony", string(fromUInt64(10)))
	if err != nil {
		log.Fatal("Error: Set Failed:", err)
	}
	log.Print("Success: Set Ok; Response:", response)

	txnXfer(c, "kevin", "tony", 300)
}

func txnXfer(etcd *clientv3.Client, from, to string, amount uint64) error {
	for {
		if ok, err := doTxnXfer(etcd, from, to, amount); err != nil {
			log.Fatal("Error: doTxnXfer: ", err)
			return err
		} else if ok {
			log.Print("Success: doTxnXfer:")
			//return nil
		}
	}
}

func toUInt64(v []byte) uint64 { x, _ := binary.Uvarint(v); return x }

func fromUInt64(v uint64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	return b[:binary.PutUvarint(b, v)]
}

func doTxnXfer(etcd *clientv3.Client, from, to string, amount uint64) (bool, error) {
	getResp, err := etcd.Txn(context.TODO()).
		Then(clientv3.OpGet(from), clientv3.OpGet(to)).
		Commit()
	if err != nil {
		return false, err
	}
	fromKV := getResp.Responses[0].GetResponseRange().Kvs[0]
	toKV := getResp.Responses[1].GetResponseRange().Kvs[0]
	fromV, toV := toUInt64(fromKV.Value), toUInt64(toKV.Value)
	if fromV < amount {
		return false, fmt.Errorf("insufficient value")
	}
	txn := etcd.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.ModRevision(from), "=", fromKV.ModRevision),
		clientv3.Compare(clientv3.ModRevision(to), "=", toKV.ModRevision))
	txn = txn.Then(
		clientv3.OpPut(from, string(fromUInt64(fromV - amount))),
		clientv3.OpPut(to, string(fromUInt64(toV + amount))))
	putResp, err := txn.Commit()
	if err != nil {
		return false, err
	}
	return putResp.Succeeded, nil
}