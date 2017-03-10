

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


