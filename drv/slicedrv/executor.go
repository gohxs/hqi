package slicedrv

import (
	"reflect"
	"sort"

	"github.com/gohxs/hqi"
)

// Exeuctor Driver way
type Executor struct {
	Coll interface{} // Collection

	collVal reflect.Value
	collTyp reflect.Type
	resList reflect.Value
}

//Match matcher implementation
func (e *Executor) Match(samples []hqi.M) {
	if len(samples) == 0 {
		e.resList = e.collVal
		return
	}
	e.resList = reflect.MakeSlice(e.collTyp, 0, 1)
	for i := 0; i < e.collVal.Len(); i++ {
		vv := e.collVal.Index(i)
		v := vv.Interface() // maybe slow :/

		for _, sample := range samples {
			//log.Println("Matching sample", sample)
			if !sMatch(sample, v) {
				continue
			}
			e.resList = reflect.Append(e.resList, vv)
		}
	}
}

//Sort implements Sorter
func (e *Executor) Sort(fields []hqi.Field) {
	if fields == nil {
		return
	}
	if e.resList == e.collVal { // Prevents manipulation on original list
		reflect.Copy(e.resList, e.collVal)
		/*e.resList = reflect.MakeSlice(e.collTyp, 0, 1)
		for i := 0; i < e.collVal.Len(); i++ {
			e.resList = reflect.Append(e.resList, e.collVal.Index(i))
		}*/
	}
	// if smaller should be true
	sort.Slice(e.resList.Interface(), func(i, j int) bool {
		for _, sf := range fields {
			sortType := sf.Value
			f1 := e.resList.Index(i).FieldByName(sf.Name)
			f2 := e.resList.Index(j).FieldByName(sf.Name)
			if !f1.IsValid() || !f2.IsValid() {
				return false
			}
			ret := typeDiff(f1.Interface(), f2.Interface())
			if ret != 0 {
				// if smaller it will be true, then check the sortType (if it is  Desc it will inverse)
				return (ret < 0) == (sortType == hqi.SortAsc)
			}
		}
		return false // just return something
	})
}

//Range implements Ranger
func (e *Executor) Range(skip, max int) {
	mlen := e.resList.Len()
	if skip > mlen {
		e.resList = reflect.MakeSlice(e.collTyp, 0, 0) // empty
		return
	}
	cMax := mlen - skip         // Computed max
	if max != 0 && max < cMax { //if qd.Max is before list content
		cMax = max // Max the list qd.Max
	}
	e.resList = e.resList.Slice(skip, skip+cMax)

}

//Retrieve implements the retriever method for executor
func (e *Executor) Retrieve(res interface{}) {
	/*if kind == hqi.ResultCount {
		reflect.ValueOf(res).Elem().Set(reflect.ValueOf(e.resList.Len()))
	}*/
	reflect.ValueOf(res).Elem().Set(e.resList)
	// Maybe others
}
