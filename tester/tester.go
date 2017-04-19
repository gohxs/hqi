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
		{"aaa", 6},
		{"aaa", 5},
		{"bbb", 4},
		{"bbb", 3},
		{"ccc", 2},
		{"ccc", 1},
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
		e.CheckEQ(fmt.Sprint(res), "[{aaa 6} {aaa 5} {aaa 6} {aaa 5}]")
	})

	t.Run("Skip&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Skip(4).Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 2} {ccc 1}]")
	})

	t.Run("Sort(name)&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Sort("-Name", "Value").Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 1} {ccc 1}]")
	})

	t.Run("Sort(-Name,-Value)&Max", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		q.Find().Sort("-Name", "-Value").Max(2).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 2} {ccc 2}]")
	})
	t.Run("Remove", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		var err error

		err = q.Find(hqi.M{"Name": "ccc", "Value": 2}).Delete()
		e.CheckEQ(err, nil)

		q.Find(hqi.M{"Name": "ccc"}).List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{ccc 1} {ccc 1}]")
	})
	t.Run("RemoveOR", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		var err error

		// Remove all bbb 6 or all aaa
		err = q.Find(hqi.M{"Name": "bbb", "Value": 4}, hqi.M{"Name": "aaa"}, hqi.M{"Name": "ccc", "Value": 2}).Delete()
		e.CheckEQ(err, nil)

		q.Find().Sort("Name", "Value").List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{bbb 3} {bbb 3} {ccc 1} {ccc 1}]")
	})

	t.Run("Update", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var err error
		e := &ErrChecker{t}

		err = q.Find(hqi.M{"Name": "aaa"}).Update(hqi.M{"Value": 20}) // Find all aaa and set value to 20
		if err == hqi.ErrNotImplemented {
			t.Skip("Not implemented")
			return
		}
		e.CheckEQ(err, nil)

		var res []Model
		err = q.Find().List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{aaa 20} {aaa 20} {bbb 4} {bbb 3} {ccc 2} {ccc 1} {aaa 20} {aaa 20} {bbb 4} {bbb 3} {ccc 2} {ccc 1}]")
	})
	t.Run("UpdateOR", func(t *testing.T) {
		q := PrepareHQI(t, dc)
		var res []Model
		e := &ErrChecker{t}
		var err error

		// Update all bbb 6 or all aaa to [xxx 999]
		err = q.Find(hqi.M{"Name": "bbb", "Value": 4}, hqi.M{"Name": "aaa"}, hqi.M{"Name": "ccc", "Value": 2}).
			Update(hqi.M{"Name": "xxx", "Value": 99})
		if err == hqi.ErrNotImplemented {
			t.Skip("Not implemented")
			return
		}

		e.CheckEQ(err, nil)

		q.Find().Sort("Name", "Value").List(&res)
		e.CheckEQ(fmt.Sprint(res), "[{bbb 3} {bbb 3} {ccc 1} {ccc 1} {xxx 99} {xxx 99} {xxx 99} {xxx 99} {xxx 99} {xxx 99} {xxx 99} {xxx 99}]")
	})
}

// Benchmark suite
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
