package etcd3

import (
	etcdrpc "go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"k8s.io/apimachinery/pkg/api/errors"
)

func interpretWatchError(err error) error {
	switch {
	case err == etcdrpc.ErrCompacted:
		return errors.NewResourceExpired("The resourceVersion for the provided watch is too old.")
	}
	return err
}

const (
	expired         string = "The resourceVersion for the provided list is too old."
	continueExpired string = "The provided continue parameter is too old " +
		"to display a consistent list result. You can start a new list without " +
		"the continue parameter."
	inconsistentContinue string = "The provided continue parameter is too old " +
		"to display a consistent list result. You can start a new list without " +
		"the continue parameter, or use the continue token in this response to " +
		"retrieve the remainder of the results. Continuing with the provided " +
		"token results in an inconsistent list - objects that were created, " +
		"modified, or deleted between the time the first chunk was returned " +
		"and now may show up in the list."
)

func interpretListError(err error, paging bool, continueKey, keyPrefix string) error {
	switch {
	case err == etcdrpc.ErrCompacted:
		if paging {
			return handleCompactedErrorForPaging(continueKey, keyPrefix)
		}
		return errors.NewResourceExpired(expired)
	}
	return err
}

func handleCompactedErrorForPaging(continueKey, keyPrefix string) error {
	// continueToken.ResoureVersion=-1 means that the apiserver can
	// continue the list at the latest resource version. We don't use rv=0
	// for this purpose to distinguish from a bad token that has empty rv.
	newToken, err := encodeContinue(continueKey, keyPrefix, -1)
	if err != nil {
		//utilruntime.HandleError(err)
		//TODO: log error
		return errors.NewResourceExpired(continueExpired)
	}
	statusError := errors.NewResourceExpired(inconsistentContinue)
	statusError.ErrStatus.ListMeta.Continue = newToken
	return statusError
}
