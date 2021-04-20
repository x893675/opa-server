package storage

import (
	"github.com/x893675/opa-server/pkg/runtime"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

// SelectionPredicate is used to represent the way to select objects from api storage.
type SelectionPredicate struct {
	Label               labels.Selector
	Field               fields.Selector
	GetAttrs            AttrFunc
	IndexLabels         []string
	IndexFields         []string
	Limit               int64
	Continue            string
	AllowWatchBookmarks bool
}

// AttrFunc returns label and field sets and the uninitialized flag for List or Watch to match.
// In any failure to parse given object, it returns error.
type AttrFunc func(obj runtime.Object) (labels.Set, fields.Set, error)

// Matches returns true if the given object's labels and fields (as
// returned by s.GetAttrs) match s.Label and s.Field. An error is
// returned if s.GetAttrs fails.
func (s *SelectionPredicate) Matches(obj runtime.Object) (bool, error) {
	if s.Empty() {
		return true, nil
	}
	labels, fields, err := s.GetAttrs(obj)
	if err != nil {
		return false, err
	}
	matched := s.Label.Matches(labels)
	if matched && s.Field != nil {
		matched = matched && s.Field.Matches(fields)
	}
	return matched, nil
}

// MatchesObjectAttributes returns true if the given labels and fields
// match s.Label and s.Field.
func (s *SelectionPredicate) MatchesObjectAttributes(l labels.Set, f fields.Set) bool {
	if s.Label.Empty() && s.Field.Empty() {
		return true
	}
	matched := s.Label.Matches(l)
	if matched && s.Field != nil {
		matched = (matched && s.Field.Matches(f))
	}
	return matched
}

// MatchesSingle will return (name, true) if and only if s.Field matches on the object's
// name.
func (s *SelectionPredicate) MatchesSingle() (string, bool) {
	if len(s.Continue) > 0 {
		return "", false
	}
	// TODO: should be namespace.name
	if name, ok := s.Field.RequiresExactMatch("metadata.name"); ok {
		return name, true
	}
	return "", false
}

// Empty returns true if the predicate performs no filtering.
func (s *SelectionPredicate) Empty() bool {
	return s.Label.Empty() && s.Field.Empty()
}

// For any index defined by IndexFields, if a matcher can match only (a subset)
// of objects that return <value> for a given index, a pair (<index name>, <value>)
// wil be returned.
func (s *SelectionPredicate) MatcherIndex() []MatchValue {
	var result []MatchValue
	for _, field := range s.IndexFields {
		if value, ok := s.Field.RequiresExactMatch(field); ok {
			result = append(result, MatchValue{IndexName: FieldIndex(field), Value: value})
		}
	}
	for _, label := range s.IndexLabels {
		if value, ok := s.Label.RequiresExactMatch(label); ok {
			result = append(result, MatchValue{IndexName: LabelIndex(label), Value: value})
		}
	}
	return result
}

// LabelIndex add prefix for label index.
func LabelIndex(label string) string {
	return "l:" + label
}

// FiledIndex add prefix for field index.
func FieldIndex(field string) string {
	return "f:" + field
}
