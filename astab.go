package astab

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type renderer struct {
	colwidth []int
	rows     [][]*string
}

func (r *renderer) str(row, col int, val *string) {
	if r.colwidth[col] < len(*val) {
		r.colwidth[col] = len(*val)
	}
	if r.rows[row] == nil {
		r.rows[row] = make([]*string, len(r.colwidth))
	}
	r.rows[row][col] = val
}

func (r *renderer) write(w io.Writer) {
	for row := range r.rows {
		for col := range r.rows[row] {
			if col > 0 {
				w.Write([]byte(" "))
			}
			w.Write([]byte(*r.rows[row][col]))
			w.Write([]byte(strings.Repeat(" ", r.colwidth[col]-len(*r.rows[row][col]))))
		}
		w.Write([]byte("\n"))
	}
}

// Write takes a slice of structs, calculates the maximum string width of each
// exported struct field value, and renders a table to write to any io.Writer.
func Write(w io.Writer, slice interface{}) error {
	slcv := reflect.ValueOf(slice)
	if slcv.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got %s", slcv.Kind())
	}
	r := renderer{rows: make([][]*string, 1+slcv.Len())}
	exported := []int{}

	{
		el := slcv.Type().Elem()
		if el.Kind() != reflect.Struct {
			return fmt.Errorf("expected slice of struct, got slice of %s", el.Kind())
		}

		for j := 0; j < el.NumField(); j++ {
			if el.Field(j).PkgPath == "" {
				exported = append(exported, j)
			}
		}
		if len(exported) == 0 {
			return fmt.Errorf("struct %q does not have exported fields", el)
		}
		r.colwidth = make([]int, len(exported))
		for j := 0; j < len(exported); j++ {
			s := el.Field(exported[j]).Name
			r.str(0, j, &s)
		}
	}
	for j := 0; j < slcv.Len(); j++ {
		for k := 0; k < len(exported); k++ {
			s := fmt.Sprint(slcv.Index(j).Field(exported[k]).Interface())
			r.str(j+1, k, &s)
		}
	}
	r.write(w)
	return nil
}
