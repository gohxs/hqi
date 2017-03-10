package mgodrv

import (
	"testing"

	"github.com/gohxs/hqi"

	mgo "gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

// Prepare driver for all tests
func PrepareDriver() *Driver {
	session, err := mgo.Dial("mongodb://admin:1q2w3e@localhost/mgo-test")
	if err != nil {
		panic(err)
	}
	db := session.DB("mgo-test")
	db.C("hqitest").DropCollection()
	driver := Driver{db.C("hqitest")}

	return &driver
}

func TestImpl(t *testing.T) {
	driver := PrepareDriver()
	tester.Test(t, driver)
}

func BenchmarkImpl(b *testing.B) {
	driver := PrepareDriver()
	tester.Benchmark(b, driver)
}

func BenchmarkQuery(b *testing.B) {
	driver := PrepareDriver()
	q := hqi.NewQuery(driver)
	tester.PrepareContent(b, q)
	b.Run("Query", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var res []tester.Model
			q.Find(hqi.M{"Name": "aaa"}).
				List(&res)
		}
	})
}

//Not so native
func BenchmarkNative(b *testing.B) {
	driver := PrepareDriver()
	q := hqi.NewQuery(driver)
	tester.PrepareContent(b, q)
	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var res []tester.Model
			driver.coll.
				Find(bson.M{"Name": "aaa"}).
				All(&res)
		}
	})
}
