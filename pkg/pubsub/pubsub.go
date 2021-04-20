package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/x893675/opa-server/pkg/signal"
	"github.com/x893675/opa-server/pkg/storage/storagebackend"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
)

// The short keepalive timeout and interval have been chosen to aggressively
// detect a failed etcd server without introducing much overhead.
const keepaliveTime = 30 * time.Second
const keepaliveTimeout = 10 * time.Second

// dialTimeout is the timeout for failing to establish a connection.
// It is set to 20 seconds as times shorter than that will cause TLS connections to fail
// on heavily loaded arm64 CPUs (issue #64649)
const dialTimeout = 20 * time.Second

func newETCD3Client(c storagebackend.TransportConfig) (*clientv3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      c.CertFile,
		KeyFile:       c.KeyFile,
		TrustedCAFile: c.TrustedCAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	// NOTE: Client relies on nil tlsConfig
	// for non-secure connections, update the implicit variable
	if len(c.CertFile) == 0 && len(c.KeyFile) == 0 && len(c.TrustedCAFile) == 0 {
		tlsConfig = nil
	}
	//networkContext := egressselector.Etcd.AsNetworkContext()
	//var egressDialer utilnet.DialFunc
	//if c.EgressLookup != nil {
	//	egressDialer, err = c.EgressLookup(networkContext)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	dialOptions := []grpc.DialOption{
		grpc.WithBlock(), // block until the underlying connection is up
		//grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
		//grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
	}
	//if egressDialer != nil {
	//	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
	//		u, err := url.Parse(addr)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return egressDialer(ctx, "tcp", u.Host)
	//	}
	//	dialOptions = append(dialOptions, grpc.WithContextDialer(dialer))
	//}
	cfg := clientv3.Config{
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    keepaliveTime,
		DialKeepAliveTimeout: keepaliveTimeout,
		DialOptions:          dialOptions,
		Endpoints:            c.ServerList,
		TLS:                  tlsConfig,
	}

	return clientv3.New(cfg)
}

func NewClient(endpoints []string) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}
	etcdclient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error:cannot connec to etcd:", err)
	}
	return etcdclient
}

type Data struct {
	Timestamp       *time.Time
	ResourceVersion int
}

func main() {
	tr := storagebackend.TransportConfig{
		ServerList:    []string{"192.168.234.130:2379"},
		KeyFile:       "",
		CertFile:      "",
		TrustedCAFile: "",
	}
	//c, err := newETCD3Client(tr)
	//if err != nil {
	//	panic(err)
	//}
	//defer c.Close()
	//resp, err := c.MemberList(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp.Members)

	producer, err := newETCD3Client(tr)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	consumer, err := newETCD3Client(tr)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	key := "/pubsub/hello"

	go func() {

		for {
			newData := time.Now().String()
			//txnResp, err := producer.KV.Txn(context.TODO()).If(
			//	notFound(key),
			//).Then(
			//	clientv3.OpPut(key, newData),
			//).Commit()
			//txnResp, err := producer.KV.Txn(context.TODO()).If(
			//	clientv3.Compare(clientv3.ModRevision(key), "=", 0),
			//).Then(
			//	clientv3.OpPut(key, newData),
			//).Commit()
			txnResp, err := producer.KV.Put(context.TODO(), key, newData, clientv3.WithPrevKV())
			if err != nil {
				fmt.Println("pub error:", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}
			fmt.Printf("%v\n", txnResp)
			fmt.Printf("+%v\n", *txnResp.PrevKv)
			//if !txnResp.Succeeded {
			//	fmt.Println("pub error: key exist")
			//	time.Sleep(2 * time.Second)
			//	continue
			//}
			time.Sleep(2 * time.Second)
		}
	}()
	stopCh := signal.SetupSignalHandler()
	<-stopCh
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}

//func main() {
//	endpoints := []string{"http://192.168.234.130:2379"}
//	producer := NewClient(endpoints)
//	consumer := NewClient(endpoints)
//
//	clientv3.NewWatcher()
//	watcherOptions := &clientv3.WatcherOptions{
//		Recursive:  false,
//		AfterIndex: 0,
//	}
//	etcdclt := consumer.Watcher("say", watcherOptions)
//
//	n := 1
//	go func() {
//		for {
//			producer.Set(context.Background(), "say", strconv.Itoa(n), &client.SetOptions{})
//			n++
//			time.Sleep(2 * time.Second)
//		}
//	}()
//
//	for {
//		resp, err := etcdclt.Next(context.TODO())
//		if err != nil {
//			fmt.Println("Error,consumer err:", err.Error())
//			panic(err)
//		}
//		if resp.Node.Dir {
//			continue
//		}
//		fmt.Printf("[%s] %s %s\n", resp.Action, resp.Node.Key, resp.Node.Value)
//	}
//}
