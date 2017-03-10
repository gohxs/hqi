package slicedrv

import (
	"reflect"

	"github.com/gohxs/hqi"
)

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
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

		//log.Println("Field :", field, "exists")
		// Check for struct and do a submatch
		// Deep match
		if field2.Type().Kind() == reflect.Struct {
			if tval, ok := value.(map[string]interface{}); ok {
				return sMatch(tval, field2.Interface())
			}
			if tval, ok := value.(hqi.M); ok {
				return sMatch(tval, field2.Interface())
			}
			//
			//log.Println("It is a struct")
			continue
			// Same
		}
		//log.Println("Field2:", field2, field)
		//Check zero too
		if value != field2.Interface() {
			return false
		}
	}
	/*
		//val1 := reflect.ValueOf(sample)
		//typ1 := reflect.TypeOf(sample)
		for i := 0; i < val1.NumField(); i++ {
			//fieldName := typ1.Field(i).Name
			field1 := val1.Field(i)
			if isZero(field1.Interface()) {
				continue
			}
			field2 := val2.Field(i)
			if field1.Interface() != field2.Interface() {
				return false
			}
		}*/
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
