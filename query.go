package hqi

import (
	"encoding/json"
	"errors"
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
			cur = Struct2M(sample)

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
