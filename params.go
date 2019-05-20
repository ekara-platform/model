package model

import (
	"errors"
	"reflect"
)

//Parameters represents the parameters coming from a descriptor
type Parameters map[string]interface{}

func createParameters(src map[string]interface{}) (Parameters, error) {
	err := mustHaveStringKeys(src)
	if err != nil {
		return Parameters{}, err
	}
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst, nil
}

func (r Parameters) inherit(parent Parameters) (Parameters, error) {
	dst, err := createParameters(map[string]interface{}{})
	if err != nil {
		return Parameters{}, err
	}
	merge(dst, parent)
	merge(dst, r)
	return dst, nil
}

func mustHaveStringKeys(m map[string]interface{}) error {
	for _, v := range m {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Map {
			if vv.Type().Key().Kind() != reflect.String {
				return errors.New("parameters should only have string keys")
			}
			// check deeper...
			return mustHaveStringKeys(v.(map[string]interface{}))
		}
	}
	return nil
}

func merge(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Map {
			// The value is a map so we try to go deeper if they have the same key type
			// Otherwise we overwrite the destination map with the source one
			vd := reflect.ValueOf(dst[k])
			if vd.Kind() != reflect.Map || vd.Type().Key() != vv.Type().Key() {
				dst[k] = make(map[string]interface{})
			}
			merge(dst[k].(map[string]interface{}), v.(map[string]interface{}))
		} else if vv.Kind() == reflect.Slice {
			// The value is a slice so we try to concatenate if they have the same element type
			// Otherwise we overwrite the destination slice with the source one
			vd := reflect.ValueOf(dst[k])
			if vd.Kind() != reflect.Slice || vd.Type().Elem() != vv.Type().Elem() {
				dst[k] = reflect.MakeSlice(reflect.SliceOf(vv.Type().Elem()), 0, vv.Len()).Interface()
				vd = reflect.ValueOf(dst[k])
			}
			dst[k] = reflect.AppendSlice(vd, vv).Interface()
		} else {
			if v != nil {
				dst[k] = v
			}
		}
	}
}
