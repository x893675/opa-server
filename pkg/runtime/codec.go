package runtime

import "bytes"

// Decode is a convenience wrapper for decoding data into an Object.
//func Decode(d Decoder, data []byte) (Object, error) {
//	obj, _, err := d.Decode(data, nil, nil)
//	return obj, err
//}

// Encode is a convenience wrapper for encoding to a []byte from an Encoder
func Encode(e Encoder, obj Object) ([]byte, error) {
	// TODO: reuse buffer
	buf := &bytes.Buffer{}
	if err := e.Encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode is a convenience wrapper for decoding data into an Object.
func Decode(d Decoder, data []byte) (Object, error) {
	obj, err := d.Decode(data, nil)
	return obj, err
}
