package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/beego/goyaml2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type YAMLConfig struct {
}

type YAMLConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

func ReadYmalReader(path string) (cnf map[string]interface{}, err error) {
	err = nil
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	err = nil
	buf, err := ioutil.ReadAll(f)
	if err != nil || len(buf) < 3 {
		return nil, err
	}

	if string(buf[0:1]) == "{" {
		log.Println("Look lile a Json,Try it")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			log.Println("It is Json map")
			return cnf, nil
		}
	}

	_map, _err := goyaml2.Read(bytes.NewBuffer(buf))
	if _err != nil {
		log.Println("Goyaml2 ERR >", string(buf), _err)
		return nil, _err
	}

	if _map == nil {
		log.Println("Goyaml2 output nil?Pls report bug\n" + string(buf))
	}

	if cnf, ok := _map.(map[string]interface{}); ok {
		return cnf, nil
	}

	return nil, errors.New("Not a map")
}

func (this *YAMLConfig) Parse(fileName string) (ConfigContainer, error) {
	y := &YAMLConfigContainer{
		data: make(map[string]interface{}),
	}

	cnf, err := ReadYmalReader(fileName)
	if err != nil {
		return nil, err
	}

	y.data = cnf
	return y, nil
}

func (this *YAMLConfigContainer) Bool(key string) (bool, error) {
	if v, ok := this.data[key].(bool); ok {
		return v, nil
	}
	return false, errors.New("Not bool value")
}

func (this *YAMLConfigContainer) Int(key string) (int, error) {
	if v, ok := this.data[key].(int64); ok {
		return int(v), nil
	}
	return 0, errors.New("Not int value")
}

func (this *YAMLConfigContainer) Int64(key string) (int64, error) {
	if v, ok := this.data[key].(int64); ok {
		return v, nil
	}
	return 0, errors.New("Not int64 value")
}

func (this *YAMLConfigContainer) Float(key string) (float64, error) {
	if v, ok := this.data[key].(float64); ok {
		return v, nil
	}
	return 0.0, errors.New("Not float64 value")
}

func (this *YAMLConfigContainer) String(key string) string {
	if v, ok := this.data[key].(string); ok {
		return v
	}
	return ""
}

func (this *YAMLConfigContainer) Set(key, val string) error {
	this.Lock()
	defer this.Unlock()
	this.data[key] = val
	return nil
}

func (this *YAMLConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := this.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("Not exist key")
}

func init() {
	Register(Yaml, &YAMLConfig{})
}
