package slicedrv

import (
	"reflect"

	"github.com/gohxs/hqi"
)

//Driver driver for slice
type Driver struct {
	// Personal info
	CollPtr interface{} // Should be a pointer to collection
}

//Schema create new schema based on obj struct
func (d *Driver) Schema(obj interface{}) error {

	nSliceTyp := reflect.SliceOf(reflect.TypeOf(obj))

	//newSlice := reflect.MakeSlice(nSliceTyp, 0, 1)
	newSlice := reflect.New(nSliceTyp)

	d.CollPtr = newSlice.Interface() // Pointer

	return nil
}

//Insert inert values to slice
func (d *Driver) Insert(objs ...interface{}) error {

	collVal := reflect.Indirect(reflect.ValueOf(d.CollPtr))
	for _, obj := range objs {
		objVal := reflect.Indirect(reflect.ValueOf(obj))
		objTyp := reflect.TypeOf(obj)
		if objTyp.Kind() == reflect.Slice {
			for i := 0; i < objVal.Len(); i++ {
				collVal = reflect.Append(collVal, objVal.Index(i))
			}
			continue
		}
		collVal = reflect.Append(collVal, objVal)
	}
	reflect.ValueOf(d.CollPtr).Elem().Set(collVal)
	//d.collPtr = collVal.Interface()
	//Add objs to thing
	return nil

}

// Query values from slice
func (d *Driver) Query(qp *hqi.QueryParam, res interface{}) error {
	ex := executor{
		Coll:    d.CollPtr,
		collTyp: reflect.TypeOf(d.CollPtr).Elem(),
		collVal: reflect.Indirect(reflect.ValueOf(d.CollPtr)),
	}
	ex.match(qp.Samples)
	ex.sort(qp.Sort)
	ex.limit(qp.Skip, qp.Max)
	ex.retrieve(res)

	return nil
}

//Count matched by samples (Not implemented)
func (d *Driver) Count(qp *hqi.QueryParam) int {
	return -1 // Not implemented
}

//Delete matched by samples
func (d *Driver) Delete(qp *hqi.QueryParam) error {
	// How will this work??

	//collTyp := reflect.TypeOf(d.CollPtr).Elem()
	collVal := reflect.Indirect(reflect.ValueOf(d.CollPtr))
	// Delete by indexes

	var resList = collVal // Back to begin so it won't mess with indexes
	for i := resList.Len() - 1; i >= 0; i-- {
		vv := resList.Index(i)
		v := vv.Interface() // maybe slow :/
		for _, sample := range qp.Samples {
			if !sMatch(sample, v) { //if match remove
				continue
			}
			resList = reflect.AppendSlice(resList.Slice(0, i),
				resList.Slice(i+1, resList.Len()))
		}
	}
	reflect.ValueOf(d.CollPtr).Elem().Set(resList)

	return nil
}

//Update elements with obj Fields based on samples
func (d *Driver) Update(qp *hqi.QueryParam, obj hqi.M) error {
	type assignval struct {
		target reflect.Value
		src    reflect.Value
	}
	// Match jjjk
	ex := executor{
		Coll:    d.CollPtr,
		collTyp: reflect.TypeOf(d.CollPtr).Elem(),
		collVal: reflect.Indirect(reflect.ValueOf(d.CollPtr)),
	}
	ex.match(qp.Samples) // Retrieve match

	//Work on resList
	valMap := []assignval{}

	for _, v := range ex.resList {
		elVal := ex.collVal.Index(v) // Original value
		newVal := reflect.New(elVal.Type())
		newVal.Elem().Set(elVal) // Set oldValue
		// Should not assign directly?
		err := objAssign(newVal.Interface(), obj)
		if err != nil {
			return err
		}
		valMap = append(valMap, assignval{elVal, newVal})
	}

	// If everything ok we start assigning
	for _, v := range valMap {
		v.target.Set(reflect.Indirect(v.src)) // Set value on the original coll with src
	}

	return nil
}
