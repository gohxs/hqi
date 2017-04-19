hqi 
==========================
hexa query interface



Introduction
---------------------------------
Work In Progress,   
hqi is meant to interface a flat table query into a common solution and also to easly mock the interface for testing units



Usage
---------------------------

to create query interface we use `q := hqi.NewQuery( driver )`
Using slicedrv
Lets define a struct:

```go
type Model struct {
	Field1 string,
	Field2 int
}
```

Next we instantiate our slice:
```go
collection := []Model{}
```

Now we get a query interface using slice driver passing our collection to driver:
```go
q := hqi.NewQuery(&slicedrv.Driver{&collection})
```

And we can start:
```go
q.Insert(Model{"f1","f2"},Model{"f1_1","f2_2"})   // Inserts the two elements into slice

var result []Model{}
q.Find(Model{"f1","f2"}).List(&result)
q.Find(`{"Field1":"f1"}`,Model{"f1_1","f2_2"}).List(&result)   // Find where Field1 is f1 OR field1 "f1_1" AND field2 "f2_2"
```

[slice example](/examples/slice/main.go)




#### TODO

* Create driver capabilitie struct so if we are registering a subDocument driver will identify if it can process
* Finder features with operators >= <= < > !=
* Update operation 
* Matcher with Greater,Smaller checks, Zero fields
* Create a cursor to fetch parcels of data


#### DONE

* Schemer (information to pass to driver when user wants a scheme) (staging)
* Insert operation (staging)
* Delete (staging)



#### Operations (planning)

oper   | function               | equivalence
-------|------------------------|------------
NOT    | Inverse match          | !
GT     | Greater Than           |  >
GTE    | Greater equal or equal | >=
LT     | Less than              | <
LTE    | Less than or equal     | <






Matcher executor
--------------------------
* Match
* Sort
* Range
* Retrieve, Count,Delete,Update (CRUD here)


Stage based query
-------------------------

On stage based means, that if we perform a 
query. means we should have only Find option/Match  
After find we can have SORT, LIMIT,RESULT  
after SORT: LIMIT,RESULT  
after RANGE: RESULT  

* Stage1 return Find/Schema/Insert (Match)
* Stage2 sort, range, results 
* Stage3 All/Delete/Count executors


#### QuerySampler

rules for Find, if we have a model:

```go
type Person struct {
	Name string
	Age int
}
```

we could create a "cloned" Sampler as:

```go
type PersonSampler struct {
	Name interface{}
	Age interface{}
}

q.FindS(PersonSampler{Age:0}).  // Usually 0 is nothing and won't match
  List(&res)
// Or also:
q.FindS(PersonSampler{Age:hqi.Greater(10)}).
  List(&res)

```




-----

- [DRIVER](/drv/README.md)


#### Changes:

* 19-04-2017: Update implemented in mgo and slice


