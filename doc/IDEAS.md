Sample ideas
-----------------------

Create Driver capabilities to identify driver for example:
Mgo,slice supports subDocuments/slices, SQL/Cassandra just flat tables



Sample query:
Considering 
```go
type Person struct {
	Name string
	Age int
	HairColor string
	Gender string
}
```
I want Females with dark hair age above 20
and if age > 30 blonde hair

Person 
  (Gender == "f" && Age > 20 && Age < 30 && HairColor == "dark") ||
  (Gender == "f" && Age >= 30 && HairColor == "blond")

Could also be:
Gender == "F" && ((Age >20 && Age < 30 && HAirColor == "dark") || (Age >= 30 && HairColor = "blond"))

```go
q.Find(PersonS{
	Gender:"f",
	HairColor:"dark",
	Age: hqi.AND(hqi.GT(20),hqi.LT(30)),
  },
PersonS{
	Gender:"f",
	HairColor:"dark",
	Age: hci.GT(30),
}).List(...)
```





#### Simplifying matcher:



```go

//Filter
q.Find(hqi.M{"Value>":10})
q.Find(hqi.GT("Value",10))
q.Find(hqi.M("Value",hqi.GT(10))
q.Find(hqi.O("Value",">",10))

//And operation

//Finding a Person older than  20 and below 40 study

q.Find( hqi.AND(hqi.F{"Age>":20},hqi.F{"Age<":40}) )
q.Find( hqi.M{"Age",hqi.AND(hqi.GT(20),hqi.LT(10))})
q.Find("Age < 10 AND Age > 20") // Needs a parser



```







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




