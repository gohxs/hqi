package hqi

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
)

const (

	//ResultCount get the number of entries instead of Content
	ResultCount = iota
	//ResultOne get one result
	ResultOne
	//ResultList get a list of results (even one)
	ResultList
	//SortAsc sort smaller to bigger
	SortAsc = iota
	//SortDesc sort bigger to smaller
	SortDesc
)

//Field common field type
type (
	Field struct {
		Name  string
		Value interface{}
	}
	M     map[string]interface{}
	Query struct {
		driver Driver
	}
)

func NewQuery(driver Driver) Query {
	if driver == nil {
		panic("Driver is nil")
	}
	return Query{driver}
}

//Find Initiates a query builder
func (q *Query) Find(samples ...interface{}) SecondStage {
	// Convert samples to hql.M
	samplesMap := []M{}
	for _, sample := range samples {
		var cur M
		switch t := sample.(type) {
		case M:
			cur = t
		case map[string]interface{}:
			cur = t
		case string:
			cur = M{}
			json.Unmarshal([]byte(t), &cur)
		default:
			//TODO: do own conversion
			cur = s2m(sample)

			/*cur = M{}

			data, _ := json.Marshal(sample)
			json.Unmarshal(data, &cur)*/
		}
		samplesMap = append(samplesMap, cur)
	}

	//Convert samples to map[string]interface{}

	b := Builder{data: QueryParam{Samples: samplesMap}, driver: q.driver}
	return &b
}

func (q *Query) Schema(sample interface{}) error {
	return q.driver.Schema(sample)
}
func (q *Query) Insert(objs ...interface{}) error {
	return q.driver.Insert(objs...)
}

// Convert struct to M
func s2m(obj interface{}) M {
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
				s = append(s, s2m(value.Index(si).Interface()))
			}
			ret[fName] = s
		case reflect.Map:
			var m = M{}
			for _, k := range value.MapKeys() {
				m[k.String()] = value.MapIndex(k).Interface()
			}
			ret[fName] = m
		case reflect.Struct:
			ret[fName] = s2m(valI) // recursive
		default:
			ret[fName] = valI

		}
	}
	return ret
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
