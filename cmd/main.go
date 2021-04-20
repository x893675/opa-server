package main

import (
	"context"
	"fmt"
	"time"

	"github.com/open-policy-agent/opa/runtime"
	"github.com/x893675/opa-server/pkg/signal"
)

func main() {
	stopCh := signal.SetupSignalHandler()
	ctx := context.TODO()
	addr := []string{":8181"}
	rt, err := runtime.NewRuntime(ctx, runtime.Params{
		Addrs:                  &addr,
		GracefulShutdownPeriod: 1,
	})
	if err != nil {
		panic(err)
	}

	errChan := make(chan error, 1)

	go func() {
		errChan <- rt.Serve(ctx)
	}()

	select {
	case err := <-errChan:
		fmt.Println("start opa server error:", err.Error())
	case <-stopCh:
		// wait opa server graceful shutdown ...
		time.Sleep(5 * time.Second)
	}
}
