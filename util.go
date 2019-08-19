package xspider

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

//SaveJSON SaveJSON
func SaveJSON(v interface{}, filename string) error {
	b, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0666)
}
