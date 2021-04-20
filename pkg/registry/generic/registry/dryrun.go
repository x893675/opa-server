package registry

import (
	"github.com/x893675/opa-server/pkg/runtime"
	"k8s.io/apiserver/pkg/storage"
)

type DryRunnableStorage struct {
	Storage storage.Interface
	Codec   runtime.Codec
}
