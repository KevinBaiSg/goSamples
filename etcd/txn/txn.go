package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	. "github.com/KevinBaiSg/goSamples/etcd/common"
	"github.com/Rhymond/go-money"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

func main() {
	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("filepath directory error ", err)
		return
	}
	viper.AddConfigPath(dir)

	c, err := NewClient()
	if err != nil {
		log.Fatal("Error: NewClient Failed:", err)
		return
	}

	kv := clientv3.NewKV(c)
	ctx := context.Background()

	// init
	putResponse, err := kv.Put(ctx, "kevin/amount", string(Int64ToBytes(10000)))
	if err != nil {
		log.Fatal("Error: Put Failed:", err)
	}
	if putResponse.PrevKv == nil {
		log.Print("Success: Put Ok; PrevKv KeyValue empty")
	} else {
		log.Print("Success: Put Ok; PrevKv KeyValue:", BytesToInt64(putResponse.PrevKv.Value))
	}

	putResponse, err = kv.Put(ctx, "tony/amount", string(Int64ToBytes(0)))
	if err != nil {
		log.Fatal("Error: Put Failed:", err)
	}
	if putResponse.PrevKv == nil {
		log.Print("Success: Put Ok; PrevKv KeyValue empty")
	} else {
		log.Print("Success: Put Ok; PrevKv KeyValue:", BytesToInt64(putResponse.PrevKv.Value))
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 10 ; i++ {
		wg.Add(1)
		go func(index int) { // 为了尽可能的出现冲突
			txnXfer(index, c, "kevin", "tony", 100)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func txnXfer(index int, etcd *clientv3.Client, from, to string, amount int64) error {
	for {
		from, to , result, err := doTxnXfer(etcd, from + "/amount", to + "/amount", amount)
		if err != nil {
			log.Printf("index: %d; Error: txnXfer Failed: %e", index, err)
		} else if result == false {
			log.Printf("index: %d; Failed: kevin: %s to tony: %s", index, from, to)
		} else {
			log.Printf("index: %d; Success: kevin: %s; tony: %s", index, from, to)
		}
	}
}

func doTxnXfer(etcd *clientv3.Client, from, to string, amount int64) (string, string, bool, error) {
	getResp, err := etcd.Txn(context.TODO()).
		Then(clientv3.OpGet(from), clientv3.OpGet(to)).
		Commit()
	if err != nil {
		return "", "", false, err
	}
	fromKV := getResp.Responses[0].GetResponseRange().Kvs[0]
	toKV := getResp.Responses[1].GetResponseRange().Kvs[0]

	fromV := money.New(BytesToInt64(fromKV.Value), "CNY")
	toV := money.New(BytesToInt64(toKV.Value), "CNY")
	amountV := money.New(amount, "CNY")
	greaterThan, err := fromV.GreaterThan(amountV)
	if err != nil {
		return "", "", false, err
	}
	if !greaterThan {
		return "", "", false, fmt.Errorf("insufficient value")
	}

	fromV, err = fromV.Subtract(amountV)
	if err != nil {
		return "", "", false, err
	}
	toV, err = toV.Add(amountV)
	if err != nil {
		return "", "", false, err
	}

	txn := etcd.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.ModRevision(from), "=", fromKV.ModRevision),
		clientv3.Compare(clientv3.ModRevision(to), "=", toKV.ModRevision))
	txn = txn.Then(
		clientv3.OpPut(from, string(Int64ToBytes(fromV.Amount()))),
		clientv3.OpPut(to, string(Int64ToBytes(toV.Amount()))))
	txnResponse, err := txn.Commit()
	if err != nil {
		return "", "", false, err
	}

	return fromV.Display(), toV.Display(), txnResponse.Succeeded, nil
}
