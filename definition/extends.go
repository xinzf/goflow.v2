package definition

import (
	"encoding/xml"
	"errors"
	"reflect"
)

type Extends struct {
	Text string `xml:",innerxml"`
}

func (static Extends) Bind(obj interface{}) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("The object is not a pointer")
	}

	return xml.Unmarshal([]byte(static.Text), obj)
}
