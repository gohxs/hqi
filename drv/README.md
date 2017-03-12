hqi Driver spec
========================

What we should know about the driver interface

#### Schema:
```go 
Schema(obj interface{}) error
```  
Schema should take a struct to define the collection 

----

#### Insert:
```go
Insert(obj ...interface{}) error
```
Inserts an object into collection

----

#### Query:
```go
Query(qp *QueryParam, res interface{}) error
```
Filters/Order/limit data based on QueryParam, and 'fill' param 2 
that should be a slice pointer

---

#### Delete:
```go
Delete(qp *QueryParam) error
```
Delete elements matched by filter

---

#### Update: (WIP)
```go
Update(qp *QueryParam, obj interface{}) error
```
Filtering by query param, takes second param to update all matches

---

#### Count:
```go
Count(qp *QueryParam) int // can error too
```

Counts the results matched by filter

---



### QueryParam Based
QueryParam is basic information on how to query/limit/organize data from the collection 

```go
type QueryParam struct {
	Samples  []M
	Sort     []Field
	Max      int
	Skip     int
}
```

Samples are composed by an slice of 'hqi.M'
Fields in the same sample should be handled as 'AND' other samples should be handled as an 'OR'

i.e:

```go
// Taking this struct as example
type User struct {
	Name string
	Level int
}

qry.Find(User{Name:"name",Level:1})
// SELECT * FROM 'collection' WHERE (Name = "name" AND Level = 1)

qry.Find(User{Name:"name1"},User{Name:"name2"})
 // SELECT * FORM 'collection' WHERE (Name = "name1" OR Name = "name2")

qry.Find(User{Name:"name1",Level:0}, User{Name:"name2"})
 // SELECT * FROM 'collection' WHERE (Name = "name1" AND Level = 0) OR (Name = "name2")
```



