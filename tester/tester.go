package tester

import (
	"fmt"
	"testing"

	"github.com/gohxs/hqi"
)

// Something to test implementation can be used in test packages for drivers
type Model struct {
	Name  string
	Value int
}

var (
	Data = []Model{
		{"aaa", 1},
		{"aaa", 2},
		{"bbb", 3},
		{"bbb", 4},
		{"ccc", 5},
		{"ccc", 6},
	}
)

func PrepareHQI(t Testing, dc func() hqi.Driver) hqi.Query {
	q := hqi.NewQuery(dc())
	//e := &ErrChecker{t}

	var err error
	err = q.Schema(Model{})
	if err != nil {
		panic(err)
	}
	err = q.Insert(Data)
	if err != nil {
		panic(err)
	}

	err = q.Insert(Data)
	if err != nil {
		panic(err)
	}
	// Needs drop delete All
	//e.MCheckEQ("Creating schema", q.Schema(Model{}), nil)
	// Double data
	//	e.MCheckEQ(fmt.Sprint("Inserting data ", Data), q.Insert(Data), nil)
	//	e.MCheckEQ(fmt.Sprint("Inserting data AGAIN", Data), q.Insert(Data), nil)

	return q
}

//Tester  testing implementation on drivers
func Test(t *testing.T, dc func() hqi.Driver) {
	{
		q := hqi.NewQuery(dc())
		e := &ErrChecker{t}
		// Needs drop delete All
		e.MCheckEQ("Creating schema", q.Schema(Model{}), nil)
		// Double data
		e.MCheckEQ(fmt.Sprint("Inserting data ", Data), q.Insert(Data), nil)
		e.MCheckEQ(fmt.Sprint("Inserting data AGAIN", Data), q.Insert(Data), nil)
	}

	// Initialize data
	t.Run("Match", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find(`{"Name":"aaa"}`).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{aaa 1} {aaa 2} {aaa 1} {aaa 2}]")
	})

	t.Run("Skip&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Skip(4).Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 5} {ccc 6}]")
	})

	t.Run("Sort(name)&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Sort("-Name", "Value").Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 5} {ccc 5}]")
	})

	t.Run("Sort(-Name,-Value)&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Sort("-Name", "-Value").Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 6} {ccc 6}]")
	})
	t.Run("Remove", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		var err error

		err = q.Find(hqi.M{"Name": "ccc", "Value": 6}).Delete()
		e.CheckEQ(err, nil)

		q.Find(hqi.M{"Name": "ccc"}).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 5} {ccc 5}]")
	})
	t.Run("RemoveOR", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		var err error

		// Remove all bbb 6 or all aaa
		err = q.Find(hqi.M{"Name": "bbb", "Value": 4}, hqi.M{"Name": "aaa"}, hqi.M{"Name": "ccc", "Value": 6}).Delete()
		e.CheckEQ(err, nil)

		q.Find().Sort("Name", "Value").List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{bbb 3} {bbb 3} {ccc 5} {ccc 5}]")

	})
}

// Benchmark
func Benchmark(b *testing.B, getDriver func() hqi.Driver) {
	q := PrepareHQI(b, getDriver)
	b.Run("Match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var res []Model
			q.Find(hqi.M{"Name": "aaa"}).
				List(&res)
		}
	})
	b.Run("Insert&Delete", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q.Insert(Model{Name: "zzz", Value: 99})
			q.Find(hqi.M{"Name": "zzz"}).Delete()

		}
	})

}
