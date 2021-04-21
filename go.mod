module github.com/x893675/opa-server

go 1.15

require (
	github.com/google/gofuzz v1.1.0
	github.com/google/uuid v1.1.2
	github.com/open-policy-agent/opa v0.27.1
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489
	google.golang.org/grpc v1.27.1
	k8s.io/apimachinery v0.21.0
	k8s.io/apiserver v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/structured-merge-diff/v4 v4.1.0
)

replace (
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 // ae9734ed278b is the SHA for git tag v3.4.13
	google.golang.org/grpc => google.golang.org/grpc v1.27.1
)
