package rest

import (
	"github.com/google/uuid"
	"github.com/x893675/opa-server/pkg/storage/meta"
)

// FillObjectMetaSystemFields populates fields that are managed by the system on ObjectMeta.
func FillObjectMetaSystemFields(m meta.Object) {
	m.SetCreationTimestamp(meta.Now())
	m.SetUID(meta.UID(uuid.New().String()))
}
