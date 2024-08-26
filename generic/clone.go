package generic

import (
	"errors"
	"reflect"
)

func CloneValueTo(dst interface{}, src interface{}) {
	y := reflect.ValueOf(dst)
	if y.Kind() != reflect.Ptr {
		panic(errors.New("invalid dst type"))
	}
	starY := y.Elem()
	x := reflect.ValueOf(src)
	starY.Set(x)
}

func CloneValue(src interface{}) interface{} {
	x := reflect.ValueOf(src)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		return starY.Interface()
	} else {
		return x.Interface()
	}
}

func MakePointerTo(obj interface{}) interface{} {
	val := reflect.ValueOf(obj)
	vp := reflect.New(val.Type())
	vp.Elem().Set(val)
	return vp.Interface()
}
