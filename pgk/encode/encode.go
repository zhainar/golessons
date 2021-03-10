package encode

import (
	"bytes"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Encode postgresql errors on Windows from Windows1251 to Utf-8
func encodeError(err error) error {
	// --- Encoding: convert s from UTF-8 to ShiftJIS
	// declare a bytes.Buffer b and an encoder which will write into this buffer
	var b bytes.Buffer

	wInUTF8 := transform.NewWriter(&b, charmap.Windows1251.NewDecoder())

	// encode our string
	wInUTF8.Write([]byte(err.Error()))
	wInUTF8.Close()

	return errors.New(b.String())
}
