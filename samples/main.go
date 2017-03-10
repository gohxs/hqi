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
