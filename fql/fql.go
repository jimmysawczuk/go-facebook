package fql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	_ = fmt.Sprintf
}

type FQLQuery struct {
	Query  string
	Params []interface{}

	Result []interface{}

	AccessToken string

	built_queries []string
}

type FQLError struct {
	Error string `json:"error_msg"`
	Code  int    `json:"error_code"`

	Query       string `json:"query"`
	AccessToken string `json:"access_token"`
}

type FQLResult struct {
	raw_result []interface{}
}

type FQLResultRow struct {
}

const fql_endpoint string = "https://api.facebook.com/method/fql.query"

func NewFQLQuery(query string, params ...interface{}) *FQLQuery {
	f := FQLQuery{
		Query:  query,
		Params: params,
	}

	return &f
}

func (this *FQLQuery) Exec() (err error) {

	err = this.build()
	if err != nil {
		return err
	}

	this.Result = make([]interface{}, 0)

	for _, query := range this.built_queries {
		resp, err := http.PostForm(fql_endpoint, url.Values{
			"format":       []string{"json"},
			"access_token": []string{this.AccessToken},
			"query":        []string{query},
		})

		fmt.Println(query)

		if err != nil {
			return err
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		error_container := FQLError{}
		err = json.Unmarshal(buf, &error_container)
		if err == nil {
			fmt.Println(err)
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
			for _, row := range res_arr {
				this.Result = append(this.Result, row)
			}
		}
	}

	return err
}

func (this *FQLQuery) build() error {
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
		param := this.Params[j]
		switch v {
		case `%s`:
			switch param.(type) {
			case string, fmt.Stringer:
				query = strings.Replace(query, `%s`, fmt.Sprintf(`'%s'`, param), 1)
			}
		case `%d`:
			switch param.(type) {
			case int, int32, int64:
				query = strings.Replace(query, `%d`, fmt.Sprintf(`'%d'`, param), 1)
			}
		case `%f`:
			switch param.(type) {
			case float32, float64:
				query = strings.Replace(query, `%f`, fmt.Sprintf(`'%f'`, param), 1)
			}
		}
		j++
	}

	fmt.Println(query)

	this.built_queries = []string{query}

	return nil
}

func (this *FQLResult) UnmarshalJSON(inc []byte) error {
	res := make([]interface{}, 0)
	err := unmarshalFQLResult(inc, &res)

	this.raw_result = res

	return err
}

func unmarshalFQLResult(inc []byte, res *[]interface{}) error {
	i := 0

	for i < len(inc) {
		switch inc[i] {
		case '[':
			pos := matching(inc, '[', ']', i)
			slice := inc[i : pos+1]
			unmarshalFQLResult(slice[1:len(slice)-1], res)
			i = pos + 1
		case '{':
			pos := matching(inc, '{', '}', i)
			slice := inc[i : pos+1]
			temp := make([]map[string]interface{}, 0)
			json.Unmarshal(slice, &temp)
			*res = append(*res, temp)
			i = pos + 1
		case ',':
			i++
		}
	}

	return nil
}

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
