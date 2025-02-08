package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type UploadContent struct {
	Name  string
	Value interface{}
}

type UploadForm struct {
	FileName string
	Contents []*UploadContent
}

func ConvertFileObjectToReader(obj interface{}) (data io.Reader, err error) {
	switch obj.(type) {
	case string:
		return strings.NewReader(obj.(string)), nil

	case []byte:
		return bytes.NewReader(obj.([]byte)), nil

	case io.Reader:
		return obj.(io.Reader), nil

	default:
		return nil, errors.New("not support file handle data")
	}

}
