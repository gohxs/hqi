package mgodrv

import (
	"fmt"
	"strings"

	"github.com/gohxs/hqi"

	mgo "gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

type Executor struct {
	driver *Driver
	mq     *mgo.Query

	// Query information
	filter    bson.M
	sort      []string
	skip, max int
}

// Convert hqiMap to bsonMap
func hm2bm(obj hqi.M) bson.M {
	var ret = bson.M{}
	for k, v := range obj { // Issue
		switch xt := v.(type) {
		case []hqi.M:
			barr := []bson.M{} // Array to be placed on and
			for _, v := range xt {
				bsub := hm2bm(v)
				bret := bson.M{}
				for k2, v2 := range bsub {
					bret[strings.ToLower(k)+"."+k2] = v2
				}
				barr = append(barr, bret)
			}
			ret["$and"] = barr
		case hqi.M: // subobject
			bsub := hm2bm(xt)
			for k2, v2 := range bsub {
				ret[strings.ToLower(k)+"."+k2] = v2
			}
		case map[string]interface{}: // this should be on Query
			bsub := hm2bm(hqi.M(xt))
			for k2, v2 := range bsub {
				ret[strings.ToLower(k)+"."+k2] = v2
			}
		default:
			ret[strings.ToLower(k)] = v
		}
		// If v is a hqi.M we should sub this
	}
	return ret
}

func (e *Executor) Match(samples []hqi.M) {
	if len(samples) == 0 {
		// filter = nil
		e.mq = e.driver.Coll.Find(nil)
		return
	}

	bsonSamples := []bson.M{}
	for _, smpl := range samples { // OR
		bSmpl := hm2bm(smpl)

		bsonSamples = append(bsonSamples, bSmpl)
	}

	//log.Println("Samples:", samples)
	//log.Println("BsonSamples:", bsonSamples)
	// If bsonSamples is 1 we pass directly in filter, else we  do a OR
	if len(bsonSamples) == 1 {
		e.filter = bsonSamples[0]
		//XXX
		//e.mq = e.driver.Coll.Find(bsonSamples[0])
	} else {
		e.filter = bson.M{"$or": bsonSamples}
		//XXX
		//e.mq = e.driver.Coll.Find(filter)
	}
	//Convert M to bson.M
	//log.Println("Coll filter:", e.filter)
}

func (e *Executor) Sort(fields []hqi.Field) {
	var sortfields []string
	for _, sf := range fields {
		if sf.Value == hqi.SortDesc {
			sortfields = append(sortfields, fmt.Sprintf("-%s", strings.ToLower(sf.Name)))
			continue
		}
		sortfields = append(sortfields, fmt.Sprintf("%s", strings.ToLower(sf.Name)))
	}
	if len(sortfields) > 0 {
		e.sort = sortfields
		//XXX e.mq = e.mq.Sort(sortfields...)
	}
}

func (e *Executor) Range(skip int, max int) {
	e.skip = skip
	e.max = max
	// Leave this to retriever
	/*if skip > 0 {
		e.mq = e.mq.Skip(skip)
	}
	if max > 0 {
		e.mq = e.mq.Limit(max)
	}*/
}

func (e *Executor) Retrieve(res interface{}) error {

	mq := e.driver.Coll.Find(e.filter)

	if len(e.sort) > 0 {
		mq = mq.Sort(e.sort...)
	}
	if e.skip > 0 {
		mq = mq.Skip(e.skip)
	}
	if e.max > 0 {
		mq = mq.Limit(e.max)
	}
	return mq.All(res)
	//err := e.mq.All(res)
	//if err != nil {
	// Set error
	//}
}

func (e *Executor) Delete() error {
	_, err := e.driver.Coll.RemoveAll(e.filter)
	return err
}
