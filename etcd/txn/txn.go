package main

import (
	"context"
	"fmt"
	"github.com/KevinBaiSg/goSamples/etcd/common"
	"github.com/coreos/etcd/clientv3"
	"log"
	"strconv"
)

func main() {
	c, err := common.NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	kv := clientv3.NewKV(c)
	ctx := context.Background()

	response, err := kv.Put(ctx, "kevin", "1000")
	if err != nil {
		log.Fatal("Error: Set Failed:", err)
	}
	log.Print("Success: Set Ok; Response:", response)

	response, err = kv.Put(ctx, "tony", "0")
	if err != nil {
		log.Fatal("Error: Set Failed:", err)
	}
	log.Print("Success: Set Ok; Response:", response)

	txnXfer(c, "kevin", "tony", 1)
}

func txnXfer(etcd *clientv3.Client, from, to string, amount uint64) error {
	for {
		if from, to , err := doTxnXfer(etcd, from, to, amount); err != nil {
			log.Fatal("Error: txnXfer Failed:", err)
		} else {
			log.Printf("Success: kevin: %s; tony: %s", from, to)
		}
	}
}

func toInt64(v string) (uint64, error) {
	return strconv.ParseUint(v, 10, 64)
}

func fromUInt64(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func doTxnXfer(etcd *clientv3.Client, from, to string, amount uint64) (string, string, error) {
	getResp, err := etcd.Txn(context.TODO()).
		Then(clientv3.OpGet(from), clientv3.OpGet(to)).
		Commit()
	if err != nil {
		return "", "", err
	}
	fromKV := getResp.Responses[0].GetResponseRange().Kvs[0]
	toKV := getResp.Responses[1].GetResponseRange().Kvs[0]
	fromV, err := toInt64(string(fromKV.Value))

	if err != nil {
		return "", "", err
	}
	toV, err := toInt64(string(toKV.Value))
	if err != nil {
		return "", "", err
	}

	if fromV < amount {
		return "", "", fmt.Errorf("insufficient value")
	}
	txn := etcd.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.ModRevision(from), "=", fromKV.ModRevision),
		clientv3.Compare(clientv3.ModRevision(to), "=", toKV.ModRevision))
	txn = txn.Then(
		clientv3.OpPut(from, string(fromUInt64(fromV - amount))),
		clientv3.OpPut(to, string(fromUInt64(toV + amount))))
	_, err = txn.Commit()
	if err != nil {
		return "", "", err
	}
	return string(fromKV.Value), string(toKV.Value), nil
}