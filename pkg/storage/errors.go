package storage

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// InternalError is generated when an error occurs in the storage package, i.e.,
// not from the underlying storage backend (e.g., etcd).
type InternalError struct {
	Reason string
}

func (e InternalError) Error() string {
	return e.Reason
}

func NewInternalErrorf(format string, a ...interface{}) InternalError {
	return InternalError{fmt.Sprintf(format, a...)}
}

// InvalidError is generated when an error caused by invalid API object occurs
// in the storage package.
type InvalidError struct {
	Errs field.ErrorList
}

func (e InvalidError) Error() string {
	return e.Errs.ToAggregate().Error()
}

func NewInvalidError(errors field.ErrorList) InvalidError {
	return InvalidError{errors}
}

const (
	ErrCodeKeyNotFound int = iota + 1
	ErrCodeKeyExists
	ErrCodeResourceVersionConflicts
	ErrCodeInvalidObj
	ErrCodeUnreachable
)

var errCodeToMessage = map[int]string{
	ErrCodeKeyNotFound:              "key not found",
	ErrCodeKeyExists:                "key exists",
	ErrCodeResourceVersionConflicts: "resource version conflicts",
	ErrCodeInvalidObj:               "invalid object",
	ErrCodeUnreachable:              "server unreachable",
}

type StorageError struct {
	Code               int
	Key                string
	ResourceVersion    int64
	AdditionalErrorMsg string
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("StorageError: %s, Code: %d, Key: %s, ResourceVersion: %d, AdditionalErrorMsg: %s",
		errCodeToMessage[e.Code], e.Code, e.Key, e.ResourceVersion, e.AdditionalErrorMsg)
}

func NewKeyNotFoundError(key string, rv int64) *StorageError {
	return &StorageError{
		Code:            ErrCodeKeyNotFound,
		Key:             key,
		ResourceVersion: rv,
	}
}

func NewKeyExistsError(key string, rv int64) *StorageError {
	return &StorageError{
		Code:            ErrCodeKeyExists,
		Key:             key,
		ResourceVersion: rv,
	}
}

func NewResourceVersionConflictsError(key string, rv int64) *StorageError {
	return &StorageError{
		Code:            ErrCodeResourceVersionConflicts,
		Key:             key,
		ResourceVersion: rv,
	}
}

func NewUnreachableError(key string, rv int64) *StorageError {
	return &StorageError{
		Code:            ErrCodeUnreachable,
		Key:             key,
		ResourceVersion: rv,
	}
}

func NewInvalidObjError(key, msg string) *StorageError {
	return &StorageError{
		Code:               ErrCodeInvalidObj,
		Key:                key,
		AdditionalErrorMsg: msg,
	}
}

var tooLargeResourceVersionCauseMsg = "Too large resource version"

// NewTooLargeResourceVersionError returns a timeout error with the given retrySeconds for a request for
// a minimum resource version that is larger than the largest currently available resource version for a requested resource.
func NewTooLargeResourceVersionError(minimumResourceVersion, currentRevision uint64, retrySeconds int) error {
	err := errors.NewTimeoutError(fmt.Sprintf("Too large resource version: %d, current: %d", minimumResourceVersion, currentRevision), retrySeconds)
	//err.ErrStatus.Details.Causes = []metav1.StatusCause{
	//	{
	//		Type:    metav1.CauseTypeResourceVersionTooLarge,
	//		Message: tooLargeResourceVersionCauseMsg,
	//	},
	//}
	return err
}
