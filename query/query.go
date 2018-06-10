package query

import (
	"net/http"
	"github.com/ashriths/go-graph/system"
	"errors"
	"encoding/json"
)

type QueryParser interface {
	RetreiveQueryParams(request interface{}) (map[string]string, error)
	RetreiveQueryData(request interface{}) (map[string]string, error)
}

type HTTPQueryParser struct {
}

func NewHTTPQueryParser() *HTTPQueryParser {
	return &HTTPQueryParser{}
}

func (hqp *HTTPQueryParser) RetreiveQueryParams(request interface{}) (map[string]string, error) {
	//panic("implement me")
	Request := request.(*http.Request)
	params:= Request.URL.Query()
	if len(params) < 1 {
		system.Logln("Query contained no parameters")
		return make(map[string]string), errors.New("Query contained no parameters")
	} else {
		queryParams := make(map[string]string)
		for param := range params {
			queryParams[param] = params.Get(param)
		}
		return queryParams, nil
	}
}

func (hqp *HTTPQueryParser) RetreiveQueryData(request interface{}) (map[string]string, error) {
	//panic("implement me")
	Request := request.(*http.Request)
	var data map[string]string
	decoder := json.NewDecoder(Request.Body)
	err := decoder.Decode(&data)
	if err != nil {
		system.Logln("Failed to decode json")
		return make(map[string]string), errors.New("Failed to decode json")
	}

	return data, nil
}



