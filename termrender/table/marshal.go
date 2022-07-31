package table

import "fmt"

// Marshal wraps an Encoder and handles any panics that may occur from invalid
// input. Any panic returns an error.
func Marshal(input interface{}) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()

	enc := &Encoder{}
	err = enc.Encode(input)
	if err != nil {
		return
	}

	return enc.Bytes(), nil
}
