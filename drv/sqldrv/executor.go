// Package sqlite test a different implementation
package sqldrv

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/gohxs/hqi"
)

//Executor sqlite executor
type Executor struct {
	driver      *Driver
	whereClause string
	orderClause string
	limitClause string
	fieldVar    []interface{}
}

//Match matcher implementation
func (e *Executor) Where(samples []hqi.M) {
	var qry bytes.Buffer
	if len(samples) == 0 {
		return
	}
	qry.WriteString("WHERE ")
	for sampleI, sample := range samples {
		if sampleI != 0 {
			qry.WriteString(" OR ")
		}
		qry.WriteString("(")
		c := 0
		for field, value := range sample {
			if c != 0 {
				qry.WriteString(" AND ")
			}
			e.fieldVar = append(e.fieldVar, value)

			var op = "="
			fname := field
			// we can check a suffix here:
			suffix := field[len(field)-1]
			if suffix == '>' { // Check for others
				op = string(suffix)
				fname = strings.TrimSpace(field[:len(field)-1])
			}

			qry.WriteString(fmt.Sprintf("%s %s $%d", fname, op, len(e.fieldVar)))
			c++
		}
		qry.WriteString(")")
	}
	//log.Println("Where:", qry.String())

	e.whereClause = qry.String()
}

//Sort sorter implementation
func (e *Executor) Sort(fields []hqi.Field) {
	var qry bytes.Buffer
	// Sorter
	for i, sort := range fields {
		if i == 0 {
			qry.WriteString("ORDER BY")
		} else {
			qry.WriteString(",")
		}
		if sort.Value == hqi.SortAsc {
			qry.WriteString(fmt.Sprintf(" %s ASC", sort.Name))
		} else {
			qry.WriteString(fmt.Sprintf(" %s DESC", sort.Name))
		}
	}
	e.orderClause = qry.String()
}

//Range ranger implementation
func (e *Executor) Limit(skip, max int) {
	var qry bytes.Buffer

	if max > 0 {
		qry.WriteString(fmt.Sprintf("LIMIT %d", max))
	}
	if skip > 0 {
		qry.WriteString(fmt.Sprintf(" OFFSET %d", skip))
	}

	e.limitClause = qry.String()
}

// Retrieve implementation
func (e *Executor) Retrieve(res interface{}) error {
	// ignore kind?

	// Build Select
	var qry bytes.Buffer
	qry.WriteString(fmt.Sprintf("SELECT"))
	elemTyp := reflect.TypeOf(res).Elem().Elem()
	for i := 0; i < elemTyp.NumField(); i++ {
		if i != 0 {
			qry.WriteString(",")
		}
		qry.WriteString(fmt.Sprintf(" %s", elemTyp.Field(i).Name))
	}
	qry.WriteString(fmt.Sprintf(" FROM %s", e.driver.TableName))

	if e.whereClause != "" {
		qry.WriteString("\n")
		qry.WriteString(e.whereClause)
	}
	if e.orderClause != "" {
		qry.WriteString("\n")
		qry.WriteString(e.orderClause)
	}
	if e.limitClause != "" {
		qry.WriteString("\n")
		qry.WriteString(e.limitClause)
	}

	// Retriever
	qryStr := qry.String()
	//log.Printf("SQL:\n%s\n", qryStr)
	qrows, err := e.driver.DB.Query(qryStr, e.fieldVar...)
	if err != nil {
		return err
	}
	//Build Fields to scan
	// Pointer of struct
	resType := reflect.TypeOf(res).Elem()
	resList := reflect.MakeSlice(resType, 0, 1)
	for qrows.Next() {
		var fields []interface{}
		sliceElem := reflect.New(elemTyp) // new Struct
		for i := 0; i < elemTyp.NumField(); i++ {
			// Create new Field
			fieldPtr := sliceElem.Elem().Field(i).Addr()
			fields = append(fields, fieldPtr.Interface())
		}
		qrows.Scan(fields...)
		//Copy struct Content
		resList = reflect.Append(resList, sliceElem.Elem())

		//Scan some how
	}
	reflect.ValueOf(res).Elem().Set(resList)

	return nil
}

func (e *Executor) Delete() error {
	var qry bytes.Buffer

	// We have where clause
	fmt.Fprintf(&qry, "DELETE FROM %s", e.driver.TableName)
	//qry.WriteString("DELETE FROM " + e.driver.TableName)
	if e.whereClause != "" {
		qry.WriteString("\n")
		qry.WriteString(e.whereClause)
	}

	qryStr := qry.String()

	//log.Println("Query:", qryStr, e.fieldVar)

	_, err := e.driver.DB.Exec(qryStr, e.fieldVar...)
	if err != nil {
		return err
	}

	return nil
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

/*func CreateExecutor(db *sql.DB, tableName string) hqi.ExecFunc {
	return func(qd *hqi.BuilderData, res interface{}) error {
		e := &Executor{db: db, tableName: tableName}

		return qd.Execute(e, res)
	}
}*/
