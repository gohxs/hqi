package hqi

//Param builder

//SecondStage Second stage where dev can Sort or do the rest
type (
	FirstStage interface {
		Where(samples ...interface{})
		SecondStage
	}
	SecondStage interface { // Or first
		Sort(...string) ThirdStage
		ThirdStage
	}
	//ThirdStage Third stage to slice data
	ThirdStage interface {
		Skip(n int) ThirdStage
		Max(n int) ThirdStage
		Limit(skip, max int) FinalStage
		FinalStage
	}
	//FinalStage  the stage to fetch data
	FinalStage interface {
		//One(res interface{}) error  // Cursor instead?
		List(res interface{}) error // cursor instead?
		Delete() error
		Count() int
	}

	//ExecFunc the main handler func
	//ExecFunc func(qd *BuilderData, res interface{}) error

	// QueryParam flat Query information
	QueryParam struct {
		Samples []M
		Sort    []Field
		Max     int
		Skip    int
	}
)

/*
func (bd *BuilderData) Execute(e Executor, res interface{}) error {
	e.Match(bd.Samples)
	e.Sort(bd.Sort)
	e.Range(bd.Skip, bd.Max)
	e.Retrieve(res)
	return nil
}*/

// Executor the main executor
type Builder struct {
	driver Driver
	//Executor Factory
	//Executor DriverFactory
	data QueryParam

	// Hide this
	//Driver   Driver
}

func (b *Builder) Sort(fields ...string) ThirdStage {
	for _, v := range fields {
		if v[0] == byte('-') {
			b.data.Sort = append(b.data.Sort, Field{Name: v[1:], Value: SortDesc})
			continue
		}
		b.data.Sort = append(b.data.Sort, Field{Name: v, Value: SortAsc})
	}
	return &*b // Return copy
}

// Ranger or Skipper,
func (b *Builder) Skip(n int) ThirdStage {
	b.data.Skip = n
	return &*b
}

// Maxer
func (b *Builder) Max(n int) ThirdStage {
	b.data.Max = n
	return &*b
}

func (b *Builder) Limit(fi, li int) FinalStage {
	b.data.Skip = fi
	if li != 0 {
		b.data.Max = li - fi
	}
	return &*b // builder copy?
}

func (b *Builder) List(res interface{}) error {
	return b.driver.Query(&b.data, res)

	/*e := b.driver.Executor()
	return b.data.Execute(e, res) /**/

	// Exec func
}
func (b *Builder) Delete() error {
	return b.driver.Delete(&b.data)
}
func (b *Builder) Count() int {
	return b.driver.Count(&b.data)
}

/*
func (b *Builder) Count() int {
	b.data.ResultKind = ResultCount
	var count int

	e := b.driver.Executor()

	err := b.data.Execute(e, &count)
	if err != nil {
		return -1
	}
	return count
}*/
