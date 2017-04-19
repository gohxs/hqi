package slicedrv

import (
	"reflect"
	"sort"

	"github.com/gohxs/hqi"
)

type resList struct {
	elemIndex int
}

// executor Driver way
type executor struct {
	Coll interface{} // Collection

	collVal reflect.Value
	collTyp reflect.Type
	//resList reflect.Value

	skip, max int

	startI, endI int
	resList      []int
}

//Match matcher implementation
func (e *executor) match(samples []hqi.M) {
	if len(samples) == 0 {
		e.resList = nil
		return
	}
	newResList := []int{}
	for i := 0; i < e.collVal.Len(); i++ {
		vv := e.collVal.Index(i)
		v := vv.Interface() // maybe slow :/
		for _, sample := range samples {
			//log.Println("Matching sample", sample)
			if !sMatch(sample, v) {
				continue
			}
			newResList = append(newResList, i) // append Index
		}
	}
	if newResList != nil {
		e.resList = newResList
	}
}

func (e *executor) sort(fields []hqi.Field) {
	if fields == nil {
		return
	}
	if e.resList == nil { // Build list
		for i := 0; i < e.collVal.Len(); i++ {
			e.resList = append(e.resList, i) // All indexes
		}
	}

	// forEach
	sort.Slice(e.resList, func(i, j int) bool {
		for _, sf := range fields {
			sortType := sf.Value
			f1 := e.collVal.Index(e.resList[i]).FieldByName(sf.Name)
			f2 := e.collVal.Index(e.resList[j]).FieldByName(sf.Name)
			if !f1.IsValid() || !f2.IsValid() {
				return false
			}
			ret := typeDiff(f1.Interface(), f2.Interface())
			if ret != 0 {
				// if smaller it will be true, then check the sortType (if it is  Desc it will inverse)
				return (ret < 0) == (sortType == hqi.SortAsc)
			}
		}
		return true // just return something
	})
}

func (e *executor) limit(skip, max int) {

	mlen := e.collVal.Len()
	// StartI endI
	if e.resList != nil { // carefull nil is empty too?
		mlen = len(e.resList)
	}
	last := mlen // Total

	// If max is setted we recalculate
	if max != 0 && (max+skip) < last {
		last = max + skip // Max the list qd.Max
	}

	e.startI = skip
	e.endI = last

}

// Create reflected list and retrieve
func (e *executor) retrieve(res interface{}) { // with limit
	resData := reflect.MakeSlice(e.collTyp, 0, 1)

	if e.resList == nil {
		//resData = e.collVal.Slice(e.startI, e.endI) // We should not allow user to edit a collection object directly
		for i := e.startI; i < e.endI; i++ {
			v := e.collVal.Index(i)
			resData = reflect.Append(resData, v)
		}
	} else {
		for i := e.startI; i < e.endI; i++ {
			v := e.collVal.Index(e.resList[i])
			resData = reflect.Append(resData, v)
		}
	}

	reflect.ValueOf(res).Elem().Set(resData) // is this possible?
}
