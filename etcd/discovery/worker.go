package discovery

import (
	"encoding/json"
	"errors"
	"log"

	. "github.com/KevinBaiSg/goSamples/etcd/common"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

// workerInfo is the service register information to etcd
type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
}

type Worker struct {
	Name		string
	Info    	WorkerInfo
	stop		chan error
	leaseID     clientv3.LeaseID
	client		*clientv3.Client
}

func NewWorker(name string, info WorkerInfo) (*Worker, error) {
	client, err := NewClient()
	if err != nil {
		log.Fatal("Error: cannot new worker client:", err)
		return nil, err
	}

	return &Worker {
		Name:		name,
		Info:		info,
		stop:		make (chan error),
		client: 	client,
	}, err
}

func (s *Worker)  Start() error {

	ch, err := s.keepAlive()
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		select {
		case err := <- s.stop:
			s.revoke()
			return err
		case <- s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
			}
		}
	}
}

func (s *Worker) Stop()  {
	s.stop <- nil
}

func (s *Worker) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error){

	info := &s.Info

	key := "services/" + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal("Error: cannot creates a new lease: ", err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal("Error: client cannot put a key-value pair into etcd: ", err)
		return nil, err
	}
	s.leaseID = resp.ID

	return  s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Worker) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseID)
	if err != nil {
		log.Println("Error: client cannot revokes the given lease: ", err)
	}
	log.Printf("servide: %s stop\n", s.Name)
	return err
}