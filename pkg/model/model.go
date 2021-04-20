package model

import "github.com/x893675/opa-server/pkg/storage/meta"

type User struct {
	meta.ObjectMeta `json:",inline"`
	Username        string `json:"username"`
}
