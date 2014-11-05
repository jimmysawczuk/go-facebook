package fql

import (
	"fmt"
	"strings"
)

type Param struct {
	data interface{}
}

type ParamSet []Param

type FQLParam interface {
	Escape() string
	Len() int
}

func NewParam(in interface{}) Param {
	p := Param{in}
	return p
}

func NewParamSet(in interface{}) ParamSet {
	p := ParamSet{}
	switch in.(type) {
	case []int:
		for _, v := range in.([]int) {
			p = append(p, NewParam(v))
		}
	case []int32:
		for _, v := range in.([]int32) {
			p = append(p, NewParam(v))
		}
	case []int64:
		for _, v := range in.([]int64) {
			p = append(p, NewParam(v))
		}
	case []float32:
		for _, v := range in.([]float32) {
			p = append(p, NewParam(v))
		}
	case []float64:
		for _, v := range in.([]float64) {
			p = append(p, NewParam(v))
		}
	case []string:
		for _, v := range in.([]string) {
			p = append(p, NewParam(v))
		}
	}

	return p
}

func (this Param) String() string {
	switch this.data.(type) {
	case int, int32, int64:
		return fmt.Sprintf(`%d`, this.data)
	case float32, float64:
		return fmt.Sprintf(`%f`, this.data)
	case string, fmt.Stringer:
		return fmt.Sprintf(`%s`, this.data)
	default:
		return ""
	}
}

func (this ParamSet) String() []string {
	tbr := []string{}
	for _, v := range this {
		tbr = append(tbr, v.String())
	}
	return tbr
}

func (this Param) Escape() string {
	return fmt.Sprintf(`'%s'`, this.String())
}

func (this Param) Len() int {
	return 1
}

func (this ParamSet) Escape() string {
	val := this.String()
	for i, v := range val {
		val[i] = fmt.Sprintf(`'%s'`, v)
	}

	return strings.Join(val, ", ")
}

func (this ParamSet) Len() int {
	return len(this)
}
