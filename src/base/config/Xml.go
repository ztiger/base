package config

import (
	"errors"
	"github.com/beego/x2j"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

type XMLConfig struct {
}

type XMLConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

func (xmls *XMLConfig) Parse(fileName string) (ConfigContainer, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	x := &XMLConfigContainer{
		data: make(map[string]interface{}),
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	d, err := x2j.DocToMap(string(content))
	if err != nil {
		return nil, err
	}

	x.data = d["config"].(map[string]interface{})
	return x, nil
}

func (this *XMLConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(this.data[key].(string))
}

func (this *XMLConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(this.data[key].(string))
}

func (this *XMLConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(this.data[key].(string), 10, 64)
}

func (this *XMLConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(this.data[key].(string), 64)
}

func (this *XMLConfigContainer) String(key string) string {
	if v, ok := this.data[key].(string); ok {
		return v
	}
	return ""
}

func (this *XMLConfigContainer) Set(key, val string) error {
	this.Lock()
	defer this.Unlock()
	this.data[key] = val
	return nil
}

func (this *XMLConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := this.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("Not exist key")
}

func init() {
	Register(Xml, &XMLConfig{})
}
