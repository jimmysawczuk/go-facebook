// Package fql is a package designed to facilitate easy creation and execution of FQL queries. Still a work in progress.
package fql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Query struct {
	Query      string
	raw_params []interface{}
	Params     []FQLParam
	Options    QueryOptions

	Result Result

	AccessToken string

	built_queries []string
}

type Error struct {
	Message string `json:"error_msg"`
	Code    int    `json:"error_code"`

	Query       string `json:"query"`
	AccessToken string `json:"access_token"`
}

type QueryOptions struct {
	MaxParams int
}

const (
	fql_endpoint string = "https://api.facebook.com/method/fql.query"
)

var (
	MaxParams int = 1
)

// Builds a query based on the format string and supplied parameters.
//
// Proper format parameters are:
// - %s - a single string
// - %d - a single integer
// - %f - a single floating point number
// - %S - an array of strings (which are concatenated in the query)
// - %D - an array of integers
// - %F - an array of floating point numbers
//
func NewQuery(query string, params ...interface{}) *Query {
	f := Query{
		Query:      query,
		raw_params: params,
		Params:     []FQLParam{},
		Options: QueryOptions{
			MaxParams: MaxParams,
		},
	}

	for _, v := range f.raw_params {
		switch v.(type) {
		case int, int32, int64, float32, float64, string:
			f.Params = append(f.Params, NewParam(v))
		case []int, []int32, []int64, []float32, []float64, []string:
			f.Params = append(f.Params, NewParamSet(v))
		}
	}

	return &f
}

// Executes the query, with the result being stored in this.Result.
func (this *Query) Exec() (err error) {

	err = this.build()
	if err != nil {
		return err
	}

	this.Result = Result{}

	for _, query := range this.built_queries {
		if query == "" {
			continue
		}

		resp, err := http.PostForm(fql_endpoint, url.Values{
			"format":       []string{"json"},
			"access_token": []string{this.AccessToken},
			"query":        []string{query},
		})

		if err != nil {
			return err
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		error_container := Error{}
		err = json.Unmarshal(buf, &error_container)
		if err == nil {
			return error_container
		} else {
			switch err.(type) {
			case *json.UnmarshalTypeError:
			// nothing
			default:
				return err
			}
		}

		var res interface{}
		err = json.Unmarshal(buf, &res)
		if err != nil {
			return err
		}

		if res_arr, ok := res.([]interface{}); ok {
			for _, _row := range res_arr {
				this.Result.raw_result = append(this.Result.raw_result, _row)
			}
		}
	}

	this.Result.Unpack()

	return err
}

func (this *Query) build() error {
	tokens := make(map[int]string)
	for i := 0; i < len(this.Query)-1; i++ {
		v := this.Query[i]
		if v == '%' {
			switch this.Query[i+1] {
			case 's', 'd', 'f', 'S', 'D', 'F':
				tokens[i] = "%" + string(this.Query[i+1])
			}
		}
	}

	if len(tokens) != len(this.Params) {
		return fmt.Errorf("Token/parameter mismatch")
	}

	query := this.Query

	j := 0
	for _, v := range tokens {
		switch v {
		case `%s`, `%d`, `%f`:
			param := this.Params[j]
			query = strings.Replace(query, v, param.Escape(), 1)
		}
		j++
	}

	this.built_queries = []string{query}

	j = 0
	for _, v := range tokens {
		switch v {
		case `%S`, `%D`, `%F`:
			param := this.Params[j]

			if param.Len() > this.Options.MaxParams && this.Options.MaxParams > 0 {
				len_built_queries := len(this.built_queries)
				for i := 0; i < len_built_queries; i++ {
					new_queries := []string{}
					for m := 0; m < param.Len(); m += this.Options.MaxParams {
						sliced_query := param.(ParamSet)[m : m+this.Options.MaxParams].Escape()
						new_query := strings.Replace(query, v, sliced_query, 1)
						new_queries = append(new_queries, new_query)
					}
					this.built_queries[i] = ""
					this.built_queries = append(this.built_queries, new_queries...)
				}
			}
		}
		j++
	}

	return nil
}

// func (this *Result) UnmarshalJSON(inc []byte) error {
// 	res := make([]interface{}, 0)
// 	err := unmarshalFQLResult(inc, &res)

// 	this.raw_result = res

// 	return err
// }

// func unmarshalFQLResult(inc []byte, res *[]interface{}) error {
// 	i := 0

// 	for i < len(inc) {
// 		switch inc[i] {
// 		case '[':
// 			pos := matching(inc, '[', ']', i)
// 			slice := inc[i : pos+1]
// 			unmarshalFQLResult(slice[1:len(slice)-1], res)
// 			i = pos + 1
// 		case '{':
// 			pos := matching(inc, '{', '}', i)
// 			slice := inc[i : pos+1]
// 			temp := make([]map[string]interface{}, 0)
// 			json.Unmarshal(slice, &temp)
// 			*res = append(*res, temp)
// 			i = pos + 1
// 		case ',':
// 			i++
// 		}
// 	}

// 	return nil
// }

func matching(inc []byte, open, close byte, start_at int) int {
	depth := 0

	for i := start_at; i < len(inc); i++ {
		if inc[i] == open {
			depth++
		}

		if inc[i] == close {
			depth--
		}

		if depth == 0 {
			return i
		}
	}

	return -1
}

func (this Error) Error() string {
	return fmt.Sprintf("FQL Error %d: %s", this.Code, this.Message)
}
