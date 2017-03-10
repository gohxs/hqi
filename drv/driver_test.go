package drv_test

import (
	"database/sql"

	"os"
	"testing"

	"github.com/gohxs/hqi"
	"github.com/gohxs/hqi/drv/mgodrv"
	"github.com/gohxs/hqi/drv/slicedrv"
	"github.com/gohxs/hqi/drv/sqldrv"
	"github.com/gohxs/hqi/tester"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	mgo "gopkg.in/mgo.v2-unstable"
)

// Run full tests on several drivers

func TestMGO(t *testing.T) {
	getDriver := func() hqi.Driver {
		session, err := mgo.Dial("mongodb://localhost/mgo-test")
		if err != nil {
			t.Error("MGO connection fail")
			return nil
		}
		coll := session.DB("mgo-test").C("hqitest")
		coll.DropCollection()

		return &mgodrv.Driver{coll}
	}
	tester.Test(t, getDriver)
}
func BenchmarkMGO(b *testing.B) {
	getDriver := func() hqi.Driver {
		session, err := mgo.Dial("mongodb://admin:1q2w3e@localhost/mgo-test")
		if err != nil {
			b.Error("MGO connection fail")
			return nil
		}
		coll := session.DB("mgo-test").C("hqitest")
		coll.DropCollection()

		return &mgodrv.Driver{coll}
	}
	tester.Benchmark(b, getDriver)
}

func TestSLICE(t *testing.T) {
	getDriver := func() hqi.Driver {
		return &slicedrv.Driver{[]struct{}{}}
	}
	tester.Test(t, getDriver)
}
func BenchmarkSLICE(b *testing.B) {
	getDriver := func() hqi.Driver {
		return &slicedrv.Driver{[]struct{}{}}
	}
	tester.Benchmark(b, getDriver)

}

func TestSQLpg(t *testing.T) {
	getDriver := func() hqi.Driver {
		db, err := sql.Open("postgres", "user=admin dbname=hqitest sslmode=disable")
		if err != nil {
			t.Error("PG connection fail", err)
			return nil
		}
		_, err = db.Exec("DROP TABLE IF EXISTS hqitest")
		if err != nil {
			t.Error("PG FAIL")
			return nil
		}

		return &sqldrv.Driver{db, "hqitest"}
	}
	tester.Test(t, getDriver)
}
func BenchmarkSQLpg(b *testing.B) {
	getDriver := func() hqi.Driver {
		db, err := sql.Open("postgres", "user=admin dbname=hqitest sslmode=disable")
		if err != nil {
			b.Error("PG connection fail")
			return nil
		}
		_, err = db.Exec("DROP TABLE IF EXISTS hqitest")
		if err != nil {
			b.Error("PG FAIL")
			return nil
		}

		return &sqldrv.Driver{db, "hqitest"}
	}
	tester.Benchmark(b, getDriver)
}

func TestSQLlite(t *testing.T) {
	getDriver := func() hqi.Driver {
		db, err := sql.Open("sqlite3", "tmp.sqlite3")
		if err != nil {
			t.Error("Sqlite fail")
			return nil
		}
		_, err = db.Exec("DROP TABLE IF EXISTS hqitest")
		if err != nil {
			t.Error("Sqlite fail")
			return nil
		}

		return &sqldrv.Driver{db, "hqitest"}
	}
	tester.Test(t, getDriver)
	//os.Remove("tmp.sqlite3")
}

func BenchmarkSQLlite(b *testing.B) {
	getDriver := func() hqi.Driver {
		db, err := sql.Open("sqlite3", "tmp.sqlite3")
		if err != nil {
			b.Error("Sqlite fail")
			return nil
		}
		_, err = db.Exec("DROP TABLE IF EXISTS hqitest")
		if err != nil {
			b.Error("Sqlite fail")
			return nil
		}

		return &sqldrv.Driver{db, "hqitest"}
	}
	tester.Benchmark(b, getDriver)

	os.Remove("tmp.sqlite3")
}
