package hqi

import (
	"reflect"
	"strings"
)

// Export this
// Convert struct to M
func Struct2M(obj interface{}) M {

	if v, ok := obj.(M); ok { // already an M pass through
		return v
	}
	var ret = M{}
	//objTyp := reflect.TypeOf(obj)
	objVal := reflect.ValueOf(obj)
	for i := 0; i < objVal.Type().NumField(); i++ {
		fieldTyp := objVal.Type().Field(i)
		value := objVal.Field(i)
		valI := value.Interface()

		fName := fieldTyp.Name
		omitEmpty := false

		// PARSE struct TAGS
		tagStr, ok := fieldTyp.Tag.Lookup("hqi")
		if ok {
			opts := strings.Split(tagStr, ",")
			if opts[0] != "" {
				fName = opts[0]
			}
			if len(opts) > 1 && opts[1] == "omitempty" {
				omitEmpty = true
			}
		}

		// Check nil or zero if omitEmpty
		if valI == nil || (isZero(valI) && omitEmpty) {
			continue
		}
		valKind := reflect.TypeOf(valI).Kind()
		switch valKind {
		case reflect.Slice:
			var s = []M{} // new slice
			for si := 0; si < value.Len(); si++ {
				s = append(s, Struct2M(value.Index(si).Interface()))
			}
			ret[fName] = s
		case reflect.Map:
			var m = M{}
			for _, k := range value.MapKeys() {
				m[k.String()] = value.MapIndex(k).Interface()
			}
			ret[fName] = m
		case reflect.Struct:
			ret[fName] = Struct2M(valI) // recursive
		default:
			ret[fName] = valI

		}
	}
	return ret
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
