package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/x893675/opa-server/pkg/watch"

	"github.com/x893675/opa-server/pkg/runtime"
	"github.com/x893675/opa-server/pkg/storage/meta"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

//TODO:
// Storage interfaces need to be separated into two groups; those that operate
// on collections and those that operate on individually named items.
// Collection interfaces:
// (Method: Current -> Proposed)
//    GET: Lister -> CollectionGetter
//    WATCH: Watcher -> CollectionWatcher
//    CREATE: Creater -> CollectionCreater
//    DELETE: (n/a) -> CollectionDeleter
//    UPDATE: (n/a) -> CollectionUpdater
//
// Single item interfaces:
// (Method: Current -> Proposed)
//    GET: Getter -> NamedGetter
//    WATCH: (n/a) -> NamedWatcher
//    CREATE: (n/a) -> NamedCreater
//    DELETE: Deleter -> NamedDeleter
//    UPDATE: Update -> NamedUpdater

// Storage is a generic interface for RESTful storage services.
// Resources which are exported to the RESTful API of apiserver need to implement this interface. It is expected
// that objects may implement any of the below interfaces.
type Storage interface {
	// New returns an empty object that can be used with Create and Update after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	New() runtime.Object
}

// KindProvider specifies a different kind for its API than for its internal storage.  This is necessary for external
// objects that are not compiled into the api server.  For such objects, there is no in-memory representation for
// the object, so they must be represented as generic objects (e.g. runtime.Unknown), but when we present the object as part of
// API discovery we want to present the specific kind, not the generic internal representation.
type KindProvider interface {
	Kind() string
}

// Lister is an object that can retrieve resources that match the provided field and label criteria.
type Lister interface {
	// NewList returns an empty object that can be used with the List call.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	NewList() runtime.Object
	// List selects resources in the storage which match to the selector. 'options' can be nil.
	List(ctx context.Context, options *meta.ListOptions) (runtime.Object, error)
}

// Getter is an object that can retrieve a named RESTful resource.
type Getter interface {
	// Get finds a resource in the storage by name and returns it.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	Get(ctx context.Context, name string, options *meta.GetOptions) (runtime.Object, error)
}

// GetterWithOptions is an object that retrieve a named RESTful resource and takes
// additional options on the get request. It allows a caller to also receive the
// subpath of the GET request.
type GetterWithOptions interface {
	// Get finds a resource in the storage by name and returns it.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	// The options object passed to it is of the same type returned by the NewGetOptions
	// method.
	// TODO: Pass metav1.GetOptions.
	Get(ctx context.Context, name string, options runtime.Object) (runtime.Object, error)

	// NewGetOptions returns an empty options object that will be used to pass
	// options to the Get method. It may return a bool and a string, if true, the
	// value of the request path below the object will be included as the named
	// string in the serialization of the runtime object. E.g., returning "path"
	// will convert the trailing request scheme value to "path" in the map[string][]string
	// passed to the converter.
	NewGetOptions() (runtime.Object, bool, string)
}

// GracefulDeleter knows how to pass deletion options to allow delayed deletion of a
// RESTful object.
type GracefulDeleter interface {
	// Delete finds a resource in the storage and deletes it.
	// The delete attempt is validated by the deleteValidation first.
	// If options are provided, the resource will attempt to honor them or return an invalid
	// request error.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	// Delete *may* return the object that was deleted, or a status object indicating additional
	// information about deletion.
	// It also returns a boolean which is set to true if the resource was instantly
	// deleted or false if it will be deleted asynchronously.
	Delete(ctx context.Context, name string, deleteValidation ValidateObjectFunc, options *meta.DeleteOptions) (runtime.Object, bool, error)
}

// MayReturnFullObjectDeleter may return deleted object (instead of a simple status) on deletion.
type MayReturnFullObjectDeleter interface {
	DeleteReturnsDeletedObject() bool
}

// CollectionDeleter is an object that can delete a collection
// of RESTful resources.
type CollectionDeleter interface {
	// DeleteCollection selects all resources in the storage matching given 'listOptions'
	// and deletes them. The delete attempt is validated by the deleteValidation first.
	// If 'options' are provided, the resource will attempt to honor them or return an
	// invalid request error.
	// DeleteCollection may not be atomic - i.e. it may delete some objects and still
	// return an error after it. On success, returns a list of deleted objects.
	DeleteCollection(ctx context.Context, deleteValidation ValidateObjectFunc, options *meta.DeleteOptions, listOptions *meta.ListOptions) (runtime.Object, error)
}

// Creater is an object that can create an instance of a RESTful object.
type Creater interface {
	// New returns an empty object that can be used with Create after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	New() runtime.Object

	// Create creates a new version of a resource.
	Create(ctx context.Context, obj runtime.Object, createValidation ValidateObjectFunc, options *meta.CreateOptions) (runtime.Object, error)
}

// NamedCreater is an object that can create an instance of a RESTful object using a name parameter.
type NamedCreater interface {
	// New returns an empty object that can be used with Create after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	New() runtime.Object

	// Create creates a new version of a resource. It expects a name parameter from the path.
	// This is needed for create operations on subresources which include the name of the parent
	// resource in the path.
	Create(ctx context.Context, name string, obj runtime.Object, createValidation ValidateObjectFunc, options *meta.CreateOptions) (runtime.Object, error)
}

// UpdatedObjectInfo provides information about an updated object to an Updater.
// It requires access to the old object in order to return the newly updated object.
type UpdatedObjectInfo interface {
	// Returns preconditions built from the updated object, if applicable.
	// May return nil, or a preconditions object containing nil fields,
	// if no preconditions can be determined from the updated object.
	Preconditions() *meta.Preconditions

	// UpdatedObject returns the updated object, given a context and old object.
	// The only time an empty oldObj should be passed in is if a "create on update" is occurring (there is no oldObj).
	UpdatedObject(ctx context.Context, oldObj runtime.Object) (newObj runtime.Object, err error)
}

// ValidateObjectFunc is a function to act on a given object. An error may be returned
// if the hook cannot be completed. A ValidateObjectFunc may NOT transform the provided
// object.
type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error

// ValidateAllObjectFunc is a "admit everything" instance of ValidateObjectFunc.
func ValidateAllObjectFunc(ctx context.Context, obj runtime.Object) error {
	return nil
}

// ValidateObjectUpdateFunc is a function to act on a given object and its predecessor.
// An error may be returned if the hook cannot be completed. An UpdateObjectFunc
// may NOT transform the provided object.
type ValidateObjectUpdateFunc func(ctx context.Context, obj, old runtime.Object) error

// ValidateAllObjectUpdateFunc is a "admit everything" instance of ValidateObjectUpdateFunc.
func ValidateAllObjectUpdateFunc(ctx context.Context, obj, old runtime.Object) error {
	return nil
}

// Updater is an object that can update an instance of a RESTful object.
type Updater interface {
	// New returns an empty object that can be used with Update after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	New() runtime.Object

	// Update finds a resource in the storage and updates it. Some implementations
	// may allow updates creates the object - they should set the created boolean
	// to true.
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo, createValidation ValidateObjectFunc, updateValidation ValidateObjectUpdateFunc, forceAllowCreate bool, options *meta.UpdateOptions) (runtime.Object, bool, error)
}

// CreaterUpdater is a storage object that must support both create and update.
// Go prevents embedded interfaces that implement the same method.
type CreaterUpdater interface {
	Creater
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo, createValidation ValidateObjectFunc, updateValidation ValidateObjectUpdateFunc, forceAllowCreate bool, options *meta.UpdateOptions) (runtime.Object, bool, error)
}

// CreaterUpdater must satisfy the Updater interface.
var _ Updater = CreaterUpdater(nil)

// Patcher is a storage object that supports both get and update.
type Patcher interface {
	Getter
	Updater
}

// Watcher should be implemented by all Storage objects that
// want to offer the ability to watch for changes through the watch api.
type Watcher interface {
	// 'label' selects on labels; 'field' selects on the object's fields. Not all fields
	// are supported; an error should be returned if 'field' tries to select on a field that
	// isn't supported. 'resourceVersion' allows for continuing/starting a watch at a
	// particular version.
	Watch(ctx context.Context, options *meta.ListOptions) (watch.Interface, error)
}

// StandardStorage is an interface covering the common verbs. Provided for testing whether a
// resource satisfies the normal storage methods. Use Storage when passing opaque storage objects.
type StandardStorage interface {
	Getter
	Lister
	CreaterUpdater
	GracefulDeleter
	CollectionDeleter
	Watcher
}

// Redirector know how to return a remote resource's location.
type Redirector interface {
	// ResourceLocation should return the remote location of the given resource, and an optional transport to use to request it, or an error.
	ResourceLocation(ctx context.Context, id string) (remoteLocation *url.URL, transport http.RoundTripper, err error)
}

// Responder abstracts the normal response behavior for a REST method and is passed to callers that
// may wish to handle the response directly in some cases, but delegate to the normal error or object
// behavior in other cases.
type Responder interface {
	// Object writes the provided object to the response. Invoking this method multiple times is undefined.
	Object(statusCode int, obj runtime.Object)
	// Error writes the provided error to the response. This method may only be invoked once.
	Error(err error)
}

// Connecter is a storage object that responds to a connection request.
type Connecter interface {
	// Connect returns an http.Handler that will handle the request/response for a given API invocation.
	// The provided responder may be used for common API responses. The responder will write both status
	// code and body, so the ServeHTTP method should exit after invoking the responder. The Handler will
	// be used for a single API request and then discarded. The Responder is guaranteed to write to the
	// same http.ResponseWriter passed to ServeHTTP.
	Connect(ctx context.Context, id string, options runtime.Object, r Responder) (http.Handler, error)

	// NewConnectOptions returns an empty options object that will be used to pass
	// options to the Connect method. If nil, then a nil options object is passed to
	// Connect. It may return a bool and a string. If true, the value of the request
	// path below the object will be included as the named string in the serialization
	// of the runtime object.
	NewConnectOptions() (runtime.Object, bool, string)

	// ConnectMethods returns the list of HTTP methods handled by Connect
	ConnectMethods() []string
}

// ResourceStreamer is an interface implemented by objects that prefer to be streamed from the server
// instead of decoded directly.
type ResourceStreamer interface {
	// InputStream should return an io.ReadCloser if the provided object supports streaming. The desired
	// api version and an accept header (may be empty) are passed to the call. If no error occurs,
	// the caller may return a flag indicating whether the result should be flushed as writes occur
	// and a content type string that indicates the type of the stream.
	// If a null stream is returned, a StatusNoContent response wil be generated.
	InputStream(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error)
}

// StorageMetadata is an optional interface that callers can implement to provide additional
// information about their Storage objects.
type StorageMetadata interface {
	// ProducesMIMETypes returns a list of the MIME types the specified HTTP verb (GET, POST, DELETE,
	// PATCH) can respond with.
	ProducesMIMETypes(verb string) []string

	// ProducesObject returns an object the specified HTTP verb respond with. It will overwrite storage object if
	// it is not nil. Only the type of the return object matters, the value will be ignored.
	ProducesObject(verb string) interface{}
}

// ResetFieldsStrategy is an optional interface that a storage object can
// implement if it wishes to provide the fields reset by its strategies.
type ResetFieldsStrategy interface {
	GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set
}

// CreateUpdateResetFieldsStrategy is a union of RESTCreateUpdateStrategy
// and ResetFieldsStrategy.
type CreateUpdateResetFieldsStrategy interface {
	RESTCreateUpdateStrategy
	ResetFieldsStrategy
}

// UpdateResetFieldsStrategy is a union of RESTUpdateStrategy
// and ResetFieldsStrategy.
type UpdateResetFieldsStrategy interface {
	RESTUpdateStrategy
	ResetFieldsStrategy
}
