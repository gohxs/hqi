package sqldrv

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gohxs/hqi"
)

type Driver struct {
	DB        *sql.DB
	TableName string
}

var (
	TypeMap = map[string]string{
		"int": "integer", "int8": "integer", "int16": "integer", "int32": "integer", "int64": "integer",
		"string": "text",
	}
)

func (d *Driver) Schema(obj interface{}) error {
	var qry bytes.Buffer
	qry.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", d.TableName))
	elemTyp := reflect.TypeOf(obj)
	for i := 0; i < elemTyp.NumField(); i++ {
		if i != 0 {
			qry.WriteString(", ")
		}
		qry.WriteString(fmt.Sprintf("%s", elemTyp.Field(i).Name))
		sqlType := TypeMap[elemTyp.Field(i).Type.Name()]
		qry.WriteString(fmt.Sprintf(" %s", sqlType))
	}
	qry.WriteString(");")

	//log.Printf("CREATE stmt:\n%s\n", qry.String())
	_, err := d.DB.Exec(qry.String())
	if err != nil {
		return err
	}

	return nil
}
func (d *Driver) Insert(objs ...interface{}) error {

	var qry bytes.Buffer
	var qryParam = []interface{}{}
	var objCount = 0

	qry.WriteString("INSERT INTO ")
	qry.WriteString(d.TableName)
	qry.WriteString(" VALUES\n")
	getObj := func(elemVal reflect.Value) {
		if objCount != 0 {
			qry.WriteString(",\n")
		}
		qry.WriteString("(")
		//elemVal := reflect.ValueOf(obj)
		for i := 0; i < elemVal.NumField(); i++ {
			if i != 0 {
				qry.WriteString(", ")
			}
			qryParam = append(qryParam, elemVal.Field(i).Interface())
			qry.WriteString(fmt.Sprintf("$%d", len(qryParam)))
		}
		qry.WriteString(")")
		// Execute here
		objCount++
	}

	// Create Insert stmt
	for _, obj := range objs {
		objTyp := reflect.TypeOf(obj)
		objVal := reflect.Indirect(reflect.ValueOf(obj))
		if objTyp.Kind() == reflect.Slice {
			for i := 0; i < objVal.Len(); i++ {
				getObj(objVal.Index(i))
			}
			continue
		}
		getObj(objVal)
	}
	qry.WriteString(";")
	//log.Printf("INSERT:\n%s %v\n", qry.String(), qryParam)
	//log.Println("  Param:", qryParam)
	_, err := d.DB.Exec(qry.String(), qryParam...)
	if err != nil {
		return fmt.Errorf("%s\n%s\n%v", err, qry.String(), qryParam)
	}
	return nil
}

func (d *Driver) Query(qp *hqi.QueryParam, res interface{}) error {
	e := Executor{driver: d}
	e.Where(qp.Samples)
	e.Sort(qp.Sort)
	e.Limit(qp.Skip, qp.Max)

	return e.Retrieve(res)
}

func (d *Driver) Count(qp *hqi.QueryParam) int {
	return -1
}
func (d *Driver) Delete(qp *hqi.QueryParam) error {
	e := Executor{driver: d}
	e.Where(qp.Samples)
	return e.Delete()
}
func (d *Driver) Update(qp *hqi.QueryParam, obj interface{}) error {
	return hqi.ErrNotImplemented
}

/*
func (d *Driver) Executor() hqi.Executor {
	ex := Executor{driver: d}
	return &ex
}*/
