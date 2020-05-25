package astab

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type renderer struct {
	colwidth []int
	rows     [][]string
}

func (r *renderer) str(row, col int, val string) {
	if r.colwidth == nil {
		r.colwidth = make([]int, col+1)
	}
	for len(r.colwidth) <= col {
		r.colwidth = append(r.colwidth, 0)
	}
	if r.colwidth[col] < len(val) {
		r.colwidth[col] = len(val)
	}

	if r.rows == nil {
		r.rows = make([][]string, row)
	}
	for len(r.rows) <= row {
		r.rows = append(r.rows, []string{})
	}
	if r.rows[row] == nil {
		r.rows[row] = make([]string, col+1)
	}
	for len(r.rows[row]) <= col {
		r.rows[row] = append(r.rows[row], "")
	}
	r.rows[row][col] = val
}

func (r *renderer) write(w io.Writer) {
	for row := range r.rows {
		for col := range r.rows[row] {
			if col > 0 {
				fmt.Fprint(w, " ")
			}
			fmt.Fprintf(w, "%s%s",
				r.rows[row][col],
				strings.Repeat(" ", r.colwidth[col]-len(r.rows[row][col])))
		}
		fmt.Fprint(w, "\n")
	}
}

// Write takes a slice of structs, calculates the maximum string width of each
// exported struct field value, and renders a table to write to any io.Writer.
func Write(w io.Writer, slice interface{}) error {
	r := renderer{}
	exported := []int{}

	slcv := reflect.ValueOf(slice)
	if slcv.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got %s", slcv.Kind())
	}

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
		for j := 0; j < len(exported); j++ {
			r.str(0, j, el.Field(exported[j]).Name)
		}
	}
	for j := 0; j < slcv.Len(); j++ {
		for k := 0; k < len(exported); k++ {
			r.str(j+1, k, fmt.Sprint(slcv.Index(j).Field(exported[k]).Interface()))
		}
	}
	//fmt.Fprintf(w, "%#v\n", r)
	r.write(w)
	return nil
}
