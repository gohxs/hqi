package main

import (
	"database/sql"
	"log"

	mgo "gopkg.in/mgo.v2-unstable"

	"github.com/gohxs/hqi"
	"github.com/gohxs/hqi/drv/mgodrv"
	"github.com/gohxs/hqi/drv/slicedrv"
	"github.com/gohxs/hqi/drv/sqldrv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

/*type Details struct {
	Key   string
	Value string
}*/
type UserAccount struct {
	Kind string
	Data map[string]string
}
type User struct {
	Name        string
	Description string
	Account     []UserAccount
}

var coll = []User{}

func main() {

	// Init driver
	drv := driverMGO()
	//drv := driverSLICE()
	q := hqi.NewQuery(drv)

	q.Insert(&User{"admin", "Administration account", []UserAccount{{"email", map[string]string{"email": "admin@domain.tld", "pwd": "1q2w3e"}}}})
	q.Insert(&User{"user", "Regular user", []UserAccount{{"email", map[string]string{"email": "user@domain.tld", "pwd": "1q2w3e"}}}})

	// Lets say user login with email and passwd

	var email = "admin@domain.tld"
	var pwd = "1q2w3e"

	var res []User

	q.Find(hqi.M{
		"Account": hqi.M{
			"Kind": "email",
			"Data": hqi.M{
				"email": email,
				"pwd":   pwd,
			},
		},
	}).List(&res)
	if len(res) > 0 {
		log.Println("Result:", res)
	} else {
		log.Println("Invalid ogin")
	}
	//q.Insert(&Person{"Luis", 36, []Data{{"test", 1}, {"test", 2}}})
	//q.Insert(&Person{"Janaina", 28, []Data{{"test4", 3}, {"test4", 4}}})
	//q.Insert(&Person{"Luis", 36, Details{"h", "167"}})

	/*mgo.SetDebug(true)
	var aLogger *log.Logger
	aLogger = log.New(os.Stderr, "", log.LstdFlags)
	mgo.SetLogger(aLogger)*/
	// Native test
	//var res []Person
	//q.Find(hqi.M{"Name": "Luis"}).List(&res)
	//q.Find(hqi.M{"Data": hqi.M{"Kind": "test"}}).List(&res)
	//log.Println("Persons:", res)
	//q.Insert(&Person{"Luis", 36})
	/*q.Find(`{"Details":{"Key":"h"}}`).List(&res)
	q.Find(hqi.M{"Data": hqi.M{"Details": hqi.M{"Key": "h"}}}).List(&res)
	log.Println("Persons:", res)
	q.Find(`{"Data":{"Kind":"test","Details":{"Key":"h"}}}`).List(&res)*/

}

var (
	postgreDSN = "user=admin dbname=hqitest sslmode=disable"
	sqliteDSN  = "tmp.sqlite3"
)

func driverSQL(sqldriver string, dsn string) *sqldrv.Driver {

	db, err := sql.Open(sqldriver, dsn)
	if err != nil {
		return nil
	}
	_, err = db.Exec("DROP TABLE IF EXISTS hqitest")
	if err != nil {
		return nil
	}
	db.Exec("DROP TABLE IF EXISTS hqitest")
	//Schema this
	driver := &sqldrv.Driver{db, "hqitest"}

	q := hqi.NewQuery(driver)
	err = q.Schema(User{})
	if err != nil {
		panic(err)
	}
	return driver
}

func driverSLICE() *slicedrv.Driver {
	return &slicedrv.Driver{&coll}
}
func driverMGO() *mgodrv.Driver {
	session, err := mgo.Dial("mongodb://localhost/mgo-test")
	if err != nil {
		return nil
	}
	coll := session.DB("mgo-test").C("hqitest")
	coll.DropCollection()

	return &mgodrv.Driver{coll}

}
