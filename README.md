hqi 
==========================
hexa query interface


----
Stage based criterian
It will return the Proper stage for a flat query


Introduction
---------------------------------
hqi means to interface a flat table query into a common solution


- [IDEAS](/doc/IDEAS.md)
- [CHANGELOG](/doc/CHANGELOG.md)




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



### Finder

- OR
	```go
// will perform obj1 OR obj2
q.Find(obj1,obj2)   
```
- AND
	```go
// will perform  f1 = 1 AND f2 = 2
q.Find(hqi.M{"f1":1,"f2":2})    
// Will delete like (name = 'aaa' AND value = 5) OR (name = 'bbb')
q.Find(hqi.M{"name":"aaa","value":5},hqi.M{"name":"bbb"}).Delete()
```

- Require complex filter
- Easy to add

#### Operations

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


Sample ideas
-----------------------
#### Different way to perform queries:

```go
	q.Select(&res).      // This could be a struct{Name string}
		Where(samples...).
		Sort("field").
		Limit(1).  // ...int
		Exec() // We know that will execute here
```
this way we can change parts independently and output for different structs?:
```go
type Model {
	Name string
	Value int
}
type ModelName {
	Name string
}

qry := q.Select("Name")
qry.Where({"Test":"1"}).Exec()

qry.Where({"Test":"2"}).Exec()
```

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
q.FindS(PersonSampler{Age:gocql.Greater(10)}).
  List(&res)

```





#### Internal operation (done)

Samples could be stored in map[string]interface{} to define fields, executors  
would read from these fields and do the operations,   
Right now there is no way for a searcher to compary something to 0   

#### Executor vs Query

was experimenting to change The executor back to  one function to make the implementation less complex
as it now user must to create the driver with some funcitons and executor with others
Single Implementation in one struct would be better


#### Possibilities

Imagine a common DB struct with several data tables

and a Frontend that we would be able to execute standard SQL

```go
db := hdb.PrepareDB();
db.Collection("User",mgodrv.Driver{d.DB("dbtest").C("user")})
db.Collection("Things",slicedrv.Driver{[]MySlice{}})
db.Collection("orders",sqldrv.Driver{sql.Open("postgres","dsn")})
db.Collection("invoices",restdrv.Driver{"http://api.domain.tld"})

rows,err := db.Query("select * from User")
rows.Next()

```




03-02-2017
----
### Added
- Implemented Delete in drivers

### Changed
- Chaged Range to Limit
- Package is now named as hqi, changed l from language to i interface (hexasoftware|High|hyper) query interface
- Sort no longer requires two fields Prefix is - would also be > greater to smaller?  before:
  ```go
//Before:
q.Find().SortAsc("field1","field2").SortDesc("field3").SortAsc("field4")
// Now:
q.Find().Sort("field1","field2","-field3","field4")
```


