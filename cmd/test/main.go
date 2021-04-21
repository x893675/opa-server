package main

import (
	"context"
	"fmt"

	"github.com/x893675/opa-server/pkg/runtime"
	"github.com/x893675/opa-server/pkg/runtime/serializer/json"
	"github.com/x893675/opa-server/pkg/storage"
	"github.com/x893675/opa-server/pkg/storage/etcd3"
	"github.com/x893675/opa-server/pkg/storage/meta"
	"github.com/x893675/opa-server/pkg/storage/storagebackend"
	"github.com/x893675/opa-server/pkg/storage/storagebackend/factory"
)

type User struct {
	meta.ObjectMeta `json:",inline"`
	Password        string `json:"password"`
}

func (u *User) SetZeroValue() error {
	if u == nil {
		return fmt.Errorf("must be ptr")
	}
	u.Name = ""
	u.Password = ""
	return nil
}

func NewUser() runtime.Object {
	return &User{}
}

func main() {
	tc := storagebackend.TransportConfig{
		ServerList:    []string{"192.168.234.130:2379"},
		KeyFile:       "",
		CertFile:      "",
		TrustedCAFile: "",
	}
	codec := json.NewSerializerWithOptions(json.SerializerOptions{})
	c := storagebackend.NewDefaultConfig("/kubecaas.io", codec)
	c.Transport = tc
	etcdClient, err := factory.NewETCD3Client(c.Transport)
	if err != nil {
		panic(err)
	}
	store := etcd3.New(etcdClient, c.Codec, NewUser, c.Prefix, c.Paging, c.LeaseManagerConfig)

	user1 := &User{
		ObjectMeta: meta.ObjectMeta{
			Name:              "test3",
			CreationTimestamp: meta.Now(),
		},
		Password: "password1",
	}

	out := NewUser()

	if err := store.Create(context.TODO(), "test3", user1, out, 0); err != nil {
		panic(err)
	}

	fmt.Println(out)

	out2 := NewUser()
	if err := store.Get(context.TODO(), "test", storage.GetOptions{}, out2); err != nil {
		panic(err)
	}
	fmt.Println(out2)

	//store.List(context.TODO(),)

}

//func main() {
//	u := User{
//		ObjectMeta: meta.ObjectMeta{
//			Name: "haha",
//		},
//		Username: "wwww",
//	}
//	data, err := json.Marshal(u)
//	if err != nil {
//		panic(err)
//	}
//	obj := New()
//	t := reflect.TypeOf(obj).Elem()
//	obj2, _ := reflect.New(t).Interface().(runtime.Object)
//	if err := json.Unmarshal(data, obj2); err != nil {
//		panic(err)
//	}
//	fmt.Println("reflect error")
//}
