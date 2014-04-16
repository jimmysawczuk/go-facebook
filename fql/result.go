package fql

import (
// "fmt"
)

type Result struct {
	raw_result []interface{}
	Fields     map[string]int
	Rows       []Row
}

type Row []interface{}

func (this Result) Len() int {
	return len(this.Rows)
}

func (this *Result) Unpack() {
	this.Fields = make(map[string]int)
	for i := range this.raw_result {
		if raw_row, ok := this.raw_result[i].(map[string]interface{}); ok {
			if len(this.Fields) == 0 {
				j := 0
				for k := range raw_row {
					this.Fields[k] = j
					j++
				}
			}

			row := []interface{}{}
			for _, v := range raw_row {
				row = append(row, v)
			}
			this.Rows = append(this.Rows, row)
		}
	}

	this.raw_result = []interface{}{}
}
