package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Parser struct {
	//data map[string]string
	body map[string]interface{}
	//form map[string]string
}

func New() *Parser {
	return &Parser{
		//data: make(map[string]string),
		body: make(map[string]interface{}),
		//form: make(map[string]string),
	}
}

func (p *Parser) Parse(r *http.Request) *Parser {

	json.NewDecoder(r.Body).Decode(&p.body)

	r.ParseForm()
	for k, v := range r.Form {
		if len(v) != 0 {
			p.body[k] = v[0]
		}
	}

	for k, v := range mux.Vars(r) {
		p.body[k] = v
	}

	fmt.Printf("%+v\n", p)

	return p
}

func (p *Parser) GetInt(key string) (int, error) {
	i := p.body[key]
	switch v := i.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case string:
		val, e := strconv.ParseInt(v, 10, 64)
		return int(val), e
	case float64:
		return int(v), nil
	default:
		return 0, errors.New("type error")
	}
}

func (p *Parser) GetInt64(key string) (int64, error) {
	i := p.body[key]
	switch v := i.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case float64:
		return int64(v), nil
	default:
		return 0, errors.New("type error")
	}
}

func (p *Parser) GetString(key string) string {
	return p.body[key].(string)
}
