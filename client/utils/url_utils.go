package utils

import (
	"github.com/gorilla/schema"
	"net/url"
)

var encoder = schema.NewEncoder()

func ToKeyValuePairs(model interface{}) (map[string][]string, error) {
	queryParams := url.Values{}
	err := encoder.Encode(model, queryParams)
	return queryParams, err

}

func ProcessAsQuery(qValues url.Values, paramMap *map[string][]string) url.Values {
	for k, vList := range *paramMap {
		for _, v := range vList {
			qValues.Add(k, v)
		}
	}
	return qValues
}
