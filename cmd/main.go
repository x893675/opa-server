package main

import (
	"context"

	"github.com/open-policy-agent/opa/runtime"
)

func main() {
	ctx := context.Background()
	addr := []string{":8181"}
	rt, err := runtime.NewRuntime(ctx, runtime.Params{
		Addrs: &addr,
	})
	if err != nil {
		panic(err)
	}
	if err := rt.Serve(ctx); err != nil {
		panic(err)
	}
}
