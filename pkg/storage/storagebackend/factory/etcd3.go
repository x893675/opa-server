package factory

import (
	"time"

	"github.com/x893675/opa-server/pkg/storage/storagebackend"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
)

const (
	// The short keepalive timeout and interval have been chosen to aggressively
	// detect a failed etcd server without introducing much overhead.
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second

	// dialTimeout is the timeout for failing to establish a connection.
	// It is set to 20 seconds as times shorter than that will cause TLS connections to fail
	// on heavily loaded arm64 CPUs (issue #64649)
	dialTimeout = 20 * time.Second

	dbMetricsMonitorJitter = 0.5
)

func NewETCD3Client(c storagebackend.TransportConfig) (*clientv3.Client, error) {
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
