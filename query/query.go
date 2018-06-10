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
	RetrieveParamByName(request interface{}, param string) (string, error)
	RetrieveDataByName(request interface{}, data string) (string, error)
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

func (hqp *HTTPQueryParser) RetrieveParamByName(request interface{}, param string) (string, error) {
	//panic("implement me")
	params, err := hqp.RetreiveQueryParams(request)
	if err != nil {
		system.Logln("Failed to parse query params")
		return "", errors.New("Failed to parse query params")
	}
	paramvalue, ok := params[param]
	if !ok {
		system.Logln("Query param ", param, " not found in request")
		return "", errors.New("Query param not found in request")
	}
	return paramvalue, nil
}

func (hqp *HTTPQueryParser) RetrieveDataByName(request interface{}, data string) (string, error) {
	//panic("implement me")
	datavalues, err := hqp.RetreiveQueryData(request)
	if err != nil {
		system.Logln("Failed to parse query data")
		return "", errors.New("Failed to parse query data")
	}
	datavalue, ok := datavalues[data]
	if !ok {
		system.Logln("Data key ", data, " not found in request")
		return "", errors.New("Data key not found in request")
	}
	return datavalue, nil
}


