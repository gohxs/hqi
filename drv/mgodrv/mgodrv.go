package mgodrv

import (
	"reflect"

	"github.com/gohxs/hqi"

	mgo "gopkg.in/mgo.v2-unstable"
)

type Driver struct {
	Coll *mgo.Collection
}

func (d *Driver) Schema(obj interface{}) error {
	return nil
}

func (d *Driver) Insert(objs ...interface{}) error {
	dObj := []interface{}{}
	for _, obj := range objs {
		// Unfortunetaly
		objVal := reflect.Indirect(reflect.ValueOf(obj))
		if reflect.TypeOf(obj).Kind() == reflect.Slice {
			for i := 0; i < objVal.Len(); i++ {
				dObj = append(dObj, objVal.Index(i).Interface())
			}
			continue
		}
		dObj = append(dObj, objVal.Interface())
	}
	return d.Coll.Insert(dObj...)
}

func (d *Driver) Query(qp *hqi.QueryParam, res interface{}) error {
	e := &executor{driver: d}

	e.match(qp.Samples)
	e.sort(qp.Sort)
	e.limit(qp.Skip, qp.Max)
	e.retrieve(res)
	return nil
}

func (d *Driver) Count(qp *hqi.QueryParam) int {
	return -1 // Not implemented
}

func (d *Driver) Delete(qp *hqi.QueryParam) error {
	e := &executor{driver: d}
	e.match(qp.Samples)
	return e.delete()
}

func (d *Driver) Update(qp *hqi.QueryParam, obj hqi.M) error {
	e := &executor{driver: d}
	e.match(qp.Samples)
	return e.update(obj)
}
