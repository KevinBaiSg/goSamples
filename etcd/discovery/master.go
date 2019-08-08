package discovery

import (
	"encoding/json"
	"log"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"

	. "github.com/KevinBaiSg/goSamples/etcd/common"
)

type Master struct {
	Path 		string
	Nodes 		map[string] *Node
	Client 		*clientv3.Client
}

//node is a client
type Node struct {
	State	bool
	Key		string
	Info    WorkerInfo
}

func NewMaster(watchPath string) (*Master,error) {
	client, err := NewClient()
	if err != nil {
		log.Fatal("Error: cannot new master client:", err)
		return nil, err
	}

	master := &Master {
		Path:	watchPath,
		Nodes:	make(map[string]*Node),
		Client: client,
	}

	//go master.WatchNodes()
	return master, err
}

func (m *Master) AddNode(key string, info *WorkerInfo) {
	node := &Node{
		State:	true,
		Key:	key,
		Info:	*info,
	}

	m.Nodes[node.Key] = node
}

func GetWorkerInfo(ev *clientv3.Event) *WorkerInfo {
	info := &WorkerInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		log.Fatal("Error: cannot Unmarshal WorkerInfo:", err)
	}
	return info
}

func (m *Master) WatchNodes()  {
	watchChan := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetWorkerInfo(ev)
				m.AddNode(string(ev.Kv.Key),info)
			case clientv3.EventTypeDelete:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				delete(m.Nodes, string(ev.Kv.Key))
			}
		}
	}
}
