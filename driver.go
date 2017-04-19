package hqi

//Driver creates executors and initialize what is needed
type Driver interface {
	// Schema
	Schema(obj interface{}) error
	Insert(obj ...interface{}) error

	// Queryier or Finder
	Query(qp *QueryParam, res interface{}) error
	Delete(qp *QueryParam) error
	Update(qp *QueryParam, data M) error

	Count(qp *QueryParam) int // can error too

	// New: update all matches in query to obj
	//Executor() Executor // Creates executor per query
}
