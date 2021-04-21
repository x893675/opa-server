package rest

import (
	"context"
	"fmt"

	"github.com/x893675/opa-server/pkg/runtime"
	"github.com/x893675/opa-server/pkg/storage/meta"
	"github.com/x893675/opa-server/pkg/storage/names"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// RESTCreateStrategy defines the minimum validation, accepted input, and
// name generation behavior to create an object that follows Kubernetes
// API conventions.
type RESTCreateStrategy interface {
	// The name generator is used when the standard GenerateName field is set.
	// The NameGenerator will be invoked prior to validation.
	names.NameGenerator

	// PrepareForCreate is invoked on create before validation to normalize
	// the object.  For example: remove fields that are not to be persisted,
	// sort order-insensitive list fields, etc.  This should not remove fields
	// whose presence would be considered a validation error.
	//
	// Often implemented as a type check and an initailization or clearing of
	// status. Clear the status because status changes are internal. External
	// callers of an api (users) should not be setting an initial status on
	// newly created objects.
	PrepareForCreate(ctx context.Context, obj runtime.Object)
	// Validate returns an ErrorList with validation errors or nil.  Validate
	// is invoked after default fields in the object have been filled in
	// before the object is persisted.  This method should not mutate the
	// object.
	Validate(ctx context.Context, obj runtime.Object) field.ErrorList
	// Canonicalize allows an object to be mutated into a canonical form. This
	// ensures that code that operates on these objects can rely on the common
	// form for things like comparison.  Canonicalize is invoked after
	// validation has succeeded but before the object has been persisted.
	// This method may mutate the object. Often implemented as a type check or
	// empty method.
	Canonicalize(obj runtime.Object)
}

// BeforeCreate ensures that common operations for all resources are performed on creation. It only returns
// errors that can be converted to api.Status. It invokes PrepareForCreate, then GenerateName, then Validate.
// It returns nil if the object should be created.
func BeforeCreate(strategy RESTCreateStrategy, ctx context.Context, obj runtime.Object) error {
	objectMeta, kerr := objectMetaAndKind(obj)
	if kerr != nil {
		return kerr
	}

	//if strategy.NamespaceScoped() {
	//	if !ValidNamespace(ctx, objectMeta) {
	//		return errors.NewBadRequest("the namespace of the provided object does not match the namespace sent on the request")
	//	}
	//} else if len(objectMeta.GetNamespace()) > 0 {
	//	objectMeta.SetNamespace(metav1.NamespaceNone)
	//}
	objectMeta.SetDeletionTimestamp(nil)
	objectMeta.SetDeletionGracePeriodSeconds(nil)
	strategy.PrepareForCreate(ctx, obj)
	FillObjectMetaSystemFields(objectMeta)
	//if len(objectMeta.GetGenerateName()) > 0 && len(objectMeta.GetName()) == 0 {
	//	objectMeta.SetName(strategy.GenerateName(objectMeta.GetGenerateName()))
	//}

	// Ensure managedFields is not set unless the feature is enabled
	//if !utilfeature.DefaultFeatureGate.Enabled(features.ServerSideApply) {
	//	objectMeta.SetManagedFields(nil)
	//}

	// ClusterName is ignored and should not be saved
	//if len(objectMeta.GetClusterName()) > 0 {
	//	objectMeta.SetClusterName("")
	//}

	if errs := strategy.Validate(ctx, obj); len(errs) > 0 {
		//return errors.NewInvalid("", objectMeta.GetName(), errs)
		return fmt.Errorf("invalid object: %s, error is %s", objectMeta.GetName(), errs.ToAggregate().Error())
	}

	// Custom validation (including name validation) passed
	// Now run common validation on object meta
	// Do this *after* custom validation so that specific error messages are shown whenever possible
	// TODO: validate
	//if errs := genericvalidation.ValidateObjectMetaAccessor(objectMeta, strategy.NamespaceScoped(), path.ValidatePathSegmentName, field.NewPath("metadata")); len(errs) > 0 {
	//	return errors.NewInvalid(kind.GroupKind(), objectMeta.GetName(), errs)
	//}

	strategy.Canonicalize(obj)

	return nil
}

// CheckGeneratedNameError checks whether an error that occurred creating a resource is due
// to generation being unable to pick a valid name.
//func CheckGeneratedNameError(strategy RESTCreateStrategy, err error, obj runtime.Object) error {
//	if !errors.IsAlreadyExists(err) {
//		return err
//	}
//
//	objectMeta, kerr := objectMetaAndKind(strategy, obj)
//	if kerr != nil {
//		return kerr
//	}
//
//	if len(objectMeta.GetGenerateName()) == 0 {
//		return err
//	}
//
//	return errors.NewServerTimeoutForKind(kind.GroupKind(), "POST", 0)
//}

// objectMetaAndKind retrieves kind and ObjectMeta from a runtime object, or returns an error.
func objectMetaAndKind(obj runtime.Object) (meta.Object, error) {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return objectMeta, nil
}
