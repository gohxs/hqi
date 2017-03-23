package main

import (
	"fmt"

	"github.com/gohxs/hqi"
	"github.com/gohxs/hqi/drv/slicedrv"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func main() {
	coll := []User{}

	q := hqi.NewQuery(&slicedrv.Driver{&coll})

	q.Insert(User{"first1", "last1", "email1@domain.tld"})
	q.Insert(User{"first2", "last2", "email2@domain.tld"}, User{"first3", "last3", "email3@domain.tld"})
	fmt.Println("collection:", coll)

	res := []User{}
	q.Find(`{"FirstName":"first1"}`, `{"FirstName":"first2"}`).List(&res)
	fmt.Println("Find result:", res)

}
