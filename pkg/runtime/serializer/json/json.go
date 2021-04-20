package json

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/x893675/opa-server/pkg/runtime"
)

// SerializerOptions holds the options which are used to configure a JSON/YAML serializer.
// example:
// (1) To configure a JSON serializer, set `Yaml` to `false`.
// (2) To configure a YAML serializer, set `Yaml` to `true`.
// (3) To configure a strict serializer that can return strictDecodingError, set `Strict` to `true`.
type SerializerOptions struct {
	// Yaml: configures the Serializer to work with JSON(false) or YAML(true).
	// When `Yaml` is enabled, this serializer only supports the subset of YAML that
	// matches JSON, and will error if constructs are used that do not serialize to JSON.
	Yaml bool

	// Pretty: configures a JSON enabled Serializer(`Yaml: false`) to produce human-readable output.
	// This option is silently ignored when `Yaml` is `true`.
	Pretty bool

	// Strict: configures the Serializer to return strictDecodingError's when duplicate fields are present decoding JSON or YAML.
	// Note that enabling this option is not as performant as the non-strict variant, and should not be used in fast paths.
	Strict bool
}

var _ runtime.Serializer = (*Serializer)(nil)

// Serializer handles encoding versioned objects into the proper JSON form
type Serializer struct {
	options    SerializerOptions
	identifier runtime.Identifier
}

func (s Serializer) Identifier() runtime.Identifier {
	return s.identifier
}

func (s Serializer) Encode(obj runtime.Object, w io.Writer) error {
	//TODO: add cacher
	//if co, ok := obj.(runtime.CacheableObject); ok {
	//	return co.CacheEncode(s.Identifier(), s.doEncode, w)
	//}
	return s.doEncode(obj, w)
}

func (s *Serializer) doEncode(obj runtime.Object, w io.Writer) error {
	//if s.options.Yaml {
	//	json, err := caseSensitiveJSONIterator.Marshal(obj)
	//	if err != nil {
	//		return err
	//	}
	//	data, err := yaml.JSONToYAML(json)
	//	if err != nil {
	//		return err
	//	}
	//	_, err = w.Write(data)
	//	return err
	//}

	//if s.options.Pretty {
	//	data, err := caseSensitiveJSONIterator.MarshalIndent(obj, "", "  ")
	//	if err != nil {
	//		return err
	//	}
	//	_, err = w.Write(data)
	//	return err
	//}
	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

func (s Serializer) Decode(data []byte, into runtime.Object) (runtime.Object, error) {
	//obj := reflect.New(reflect.TypeOf(into).Elem())
	//TODO: deepcopy ?
	if err := json.Unmarshal(data, into); err != nil {
		return nil, err
	}
	return into, nil
}

// NewSerializerWithOptions creates a JSON/YAML serializer that handles encoding versioned objects into the proper JSON/YAML
// form. If typer is not nil, the object has the group, version, and kind fields set. Options are copied into the Serializer
// and are immutable.
func NewSerializerWithOptions(options SerializerOptions) *Serializer {
	return &Serializer{
		options:    options,
		identifier: identifier(options),
	}
}

// identifier computes Identifier of Encoder based on the given options.
func identifier(options SerializerOptions) runtime.Identifier {
	result := map[string]string{
		"name":   "json",
		"yaml":   strconv.FormatBool(options.Yaml),
		"pretty": strconv.FormatBool(options.Pretty),
	}
	identifier, err := json.Marshal(result)
	if err != nil {
		// TODO: print error
		//klog.Fatalf("Failed marshaling identifier for json Serializer: %v", err)
	}
	return runtime.Identifier(identifier)
}
