package slicedrv

import (
	"reflect"

	"github.com/gohxs/hqi"
)

type Driver struct {
	// Personal info
	CollPtr interface{} // Should be a pointer to collection
}

func (d *Driver) Schema(obj interface{}) error {

	nSliceTyp := reflect.SliceOf(reflect.TypeOf(obj))

	//newSlice := reflect.MakeSlice(nSliceTyp, 0, 1)
	newSlice := reflect.New(nSliceTyp)

	d.CollPtr = newSlice.Interface() // Pointer

	return nil
}
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

func (d *Driver) Query(qp *hqi.QueryParam, res interface{}) error {
	ex := Executor{
		Coll:    d.CollPtr,
		collTyp: reflect.TypeOf(d.CollPtr).Elem(),
		collVal: reflect.Indirect(reflect.ValueOf(d.CollPtr)),
	}
	ex.Match(qp.Samples)
	ex.Sort(qp.Sort)
	ex.Range(qp.Skip, qp.Max)
	ex.Retrieve(res)

	//TODO to be implemented
	return nil
}

func (d *Driver) Count(qp *hqi.QueryParam) int {
	return -1
}

func (d *Driver) Delete(qp *hqi.QueryParam) error {
	// How will this work??

	//collTyp := reflect.TypeOf(d.CollPtr).Elem()
	collVal := reflect.Indirect(reflect.ValueOf(d.CollPtr))

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

func (d *Driver) Update(p *hqi.QueryParam, obj interface{}) error {
	return hqi.ErrNotImplemented
}

// Should be named differently?
/*func (d *Driver) Executor() hqi.Executor {
	// coll
	ex := Executor{
		Coll:    d.collPtr,
		collTyp: reflect.TypeOf(d.collPtr).Elem(),
		collVal: reflect.Indirect(reflect.ValueOf(d.collPtr)),
	}

	return &ex
}*/
