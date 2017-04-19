package main

import (
	"fmt"
	"log"

	"github.com/gohxs/hqi"
	"github.com/gohxs/hqi/drv/slicedrv"
)

type User struct {
	Name   string
	Number int
}

func main() {
	coll := []User{}

	q := hqi.NewQuery(&slicedrv.Driver{&coll})

	q.Insert(User{"a", 3}, User{"b", 2}, User{"c", 1})
	q.Insert(User{"a", 3}, User{"b", 2}, User{"c", 1})
	fmt.Println("collection:", coll)
	res := []User{}

	// Limit
	q.Find().Sort("-Name").List(&res)
	log.Println("Result for Sort 1:", res)

	q.Find().Skip(1).List(&res)

	q.Find().Max(1).List(&res)
	log.Println("Result for Max 1:", res)

	//find
	sqry := q.Find(`{"Name":"a"}`, `{"Name":"b"}`)

	sqry.List(&res)
	log.Println("Name:a or Name:b", res)
	sqry.Limit(1, 0).List(&res)
	log.Println("Limit(1,0):", res)

	sqry.Limit(0, 2).List(&res)
	log.Println("Limit(0,2):", res)
	sqry.Update(hqi.M{"Number": 99}) // Why not?

	sqry.List(&res)
	log.Println("After Update Name:a or Name:b", res)

	fmt.Println("collection:", coll)

}
