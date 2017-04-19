package slicedrv

import (
	"reflect"

	"github.com/gohxs/hqi"
)

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func objAssign(target, src interface{}) error {
	// Target must be pointer
	//targetTyp := reflect.TypeOf(target)
	/*if targetTyp.Kind() != reflect.Ptr {
		return errors.New("Target must be pointer")
	}*/
	targetVal := reflect.Indirect(reflect.ValueOf(target))

	switch pv := src.(type) {
	case hqi.M: // hqi.Map
		for k, v := range pv {
			// TODO: do something with k if it has a dot like sub struct
			fv := targetVal.FieldByName(k)
			if !fv.IsValid() {
				continue
			}
			fv.Set(reflect.ValueOf(v))
		}

	}
	return nil
}

//nonZero Matcher
// Extend to use map too,
// but if it is a map we don't use 0
func sMatch(sample hqi.M, obj2 interface{}) bool {
	// we can cache obj1
	val2 := reflect.ValueOf(obj2)
	for field, value := range sample {
		field2 := val2.FieldByName(field)

		if !field2.IsValid() { //unmatched field
			return false
		}
		if field2.Type().Kind() == reflect.Struct {
			if tval, ok := value.(map[string]interface{}); ok {
				return sMatch(tval, field2.Interface())
			}
			if tval, ok := value.(hqi.M); ok {
				return sMatch(tval, field2.Interface())
			}
			//
			continue
		}
		//Check zero too
		if value != field2.Interface() {
			return false
		}
	}
	return true
}

/// Huge func
func typeDiff(v1, v2 interface{}) int64 {
	// check same type?
	// Default return 0
	switch v1.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v1).Int() - reflect.ValueOf(v2).Int()
	case uint, uint8, uint16, uint32, uint64:
		i1, i2 := reflect.ValueOf(v1).Uint(), reflect.ValueOf(v2).Uint()
		if i1 >= i2 {
			return 1
		}
		if i1 < i2 {
			return -1
		}
	case float32:
		f1, f2 := v1.(float32), v2.(float32)
		if f1 > f2 {
			return 1
		}
		if f1 < f2 {
			return -1
		}
	case float64:
		f1, f2 := v1.(float64), v2.(float64)
		if f1 > f2 {
			return 1
		}
		if f1 < f2 {
			return -1
		}

	case string:
		s1, s2 := v1.(string), v2.(string)
		if s1 > s2 {
			return 1
		}
		if s1 < s2 {
			return -1
		}
	}

	return 0
}
