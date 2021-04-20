package meta

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

// GetOptions is the standard query options to the standard REST get call.
type GetOptions struct {
	// When specified:
	// - if unset, then the result is returned from remote storage based on quorum-read flag;
	// - if it's 0, then we simply return what we currently have in cache, no guarantee;
	// - if set to non zero, then the result is at least as fresh as given rv.
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,1,opt,name=resourceVersion"`
	// +k8s:deprecated=includeUninitialized,protobuf=2
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	// A selector based on labels
	LabelSelector labels.Selector
	// A selector based on fields
	FieldSelector fields.Selector
	// If true, watch for changes to this list
	Watch bool
	// allowWatchBookmarks requests watch events with type "BOOKMARK".
	// Servers that do not implement bookmarks may ignore this flag and
	// bookmarks are sent at the server's discretion. Clients should not
	// assume bookmarks are returned at any specific interval, nor may they
	// assume the server will send any BOOKMARK event during a session.
	// If this is not a watch, this field is ignored.
	// If the feature gate WatchBookmarks is not enabled in apiserver,
	// this field is ignored.
	AllowWatchBookmarks bool
	// When specified with a watch call, shows changes that occur after that particular version of a resource.
	// Defaults to changes from the beginning of history.
	// When specified for list:
	// - if unset, then the result is returned from remote storage based on quorum-read flag;
	// - if it's 0, then we simply return what we currently have in cache, no guarantee;
	// - if set to non zero, then the result is at least as fresh as given rv.
	ResourceVersion string
	// Timeout for the list/watch call.
	TimeoutSeconds *int64
	// Limit specifies the maximum number of results to return from the server. The server may
	// not support this field on all resource types, but if it does and more results remain it
	// will set the continue field on the returned list object.
	Limit int64
	// Continue is a token returned by the server that lets a client retrieve chunks of results
	// from the server by specifying limit. The server may reject requests for continuation tokens
	// it does not recognize and will return a 410 error if the token can no longer be used because
	// it has expired.
	Continue string
}

// CreateOptions may be provided when creating an API object.
type CreateOptions struct {
	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
	// +k8s:deprecated=includeUninitialized,protobuf=2
}

// UpdateOptions may be provided when updating an API object.
// All fields in UpdateOptions should also be present in PatchOptions.
type UpdateOptions struct {
	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty" protobuf:"bytes,1,rep,name=dryRun"`
}

// DeleteOptions may be provided when deleting an API object.
type DeleteOptions struct {
	// The duration in seconds before the object should be deleted. Value must be non-negative integer.
	// The value zero indicates delete immediately. If this value is nil, the default grace period for the
	// specified type will be used.
	// Defaults to a per object value if not specified. zero means delete immediately.
	// +optional
	GracePeriodSeconds *int64 `json:"gracePeriodSeconds,omitempty" protobuf:"varint,1,opt,name=gracePeriodSeconds"`

	// Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status will be
	// returned.
	// +k8s:conversion-gen=false
	// +optional
	Preconditions *Preconditions `json:"preconditions,omitempty" protobuf:"bytes,2,opt,name=preconditions"`

	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty" protobuf:"bytes,5,rep,name=dryRun"`
}

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
type Preconditions struct {
	// Specifies the target UID.
	// +optional
	UID *UID `json:"uid,omitempty" protobuf:"bytes,1,opt,name=uid,casttype=k8s.io/apimachinery/pkg/types.UID"`
	// Specifies the target ResourceVersion
	// +optional
	ResourceVersion *string `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
}
