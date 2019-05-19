package tools

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

type Value struct {
	data interface{}
	kind reflect.Kind
}

func NewValue(obj interface{}) *Value {
	return &Value{
		data: obj,
		kind: reflect.ValueOf(obj).Kind(),
	}
}

func (this *Value) Int() int {
	if this.data == nil {
		return 0
	}
	switch this.kind {
	case reflect.Int:
		return this.data.(int)
	case reflect.Int64:
		return int(this.data.(int64))
	case reflect.Int32:
		return int(this.data.(int32))
	case reflect.Int8:
		return int(this.data.(int8))
	case reflect.Float64:
		d := decimal.NewFromFloat(this.data.(float64))
		return int(d.IntPart())
	case reflect.Float32:
		d := decimal.NewFromFloat32(this.data.(float32))
		return int(d.IntPart())
	case reflect.String:
		d, err := decimal.NewFromString(this.data.(string))
		if err != nil {
			logrus.Errorln(err)
			return 0
		}
		return int(d.IntPart())
	default:
		return 0
	}
}

func (this *Value) Int64() int64 {
	if this.data == nil {
		return 0
	}
	switch this.kind {
	case reflect.Int64:
		return this.data.(int64)
	default:
		return int64(this.Int())
	}
}

func (this *Value) Float64() float64 {
	if this.data == nil {
		return 0
	}
	switch this.kind {
	case reflect.Float64:
		return this.data.(float64)
	case reflect.Float32:
		return float64(this.data.(float32))
	default:
		return float64(this.Int64())
	}
}

func (this *Value) String() string {
	if this.data == nil {
		return ""
	}
	switch this.kind {
	case reflect.String:
		return this.data.(string)
	case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int8, reflect.Float32, reflect.Float64:
		d := decimal.NewFromFloat(this.Float64())
		return d.String()
	default:
		return ""
	}
}

func (this *Value) GetData() interface{} {
	return this.data
}

func (this *Value) Bind(obj interface{}) error {
	if this.data == nil {
		return nil
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("Argument is not a pointer at Value:Bind")
	}

	if this.kind.String() != reflect.TypeOf(obj).Kind().String() {
		return errors.New("The value type is not the argument type")
	}

	b, err := jsoniter.Marshal(this.data)
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(b, obj)
}

func (this *Value) StringSlice(separator string) []string {
	return strings.Split(this.String(), separator)
}

func (this *Value) IntSlice(separator string) []int {
	datas := make([]int, 0)
	for _, s := range this.StringSlice(separator) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		datas = append(datas, i)
	}
	return datas
}

func (this *Value) Boolean() bool {
	var ret bool
	switch this.String() {
	case "true", "True", "TRUE", "1":
		ret = true
	}
	return ret
}

func (this *Value) ParseVariable() (string, bool) {
	if strings.HasPrefix(this.String(), "${") && strings.HasSuffix(this.String(), "}") {
		return strings.TrimSuffix(strings.TrimPrefix(this.String(), "${"), "}"), true
	}

	return this.String(), false
}
