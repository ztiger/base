package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type JsonConfig struct {
}

type JsonConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

func (this *JsonConfig) Parse(fileName string) (ConfigContainer, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	x := &JsonConfigContainer{
		data: make(map[string]interface{}),
		//sync.RWMutex{}
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &x.data)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (this *JsonConfigContainer) getdata(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	if len(key) == 0 {
		return nil
	}

	sectionKey := strings.Split(key, "::")
	if len(sectionKey) >= 2 {
		cruval, ok := this.data[sectionKey[0]]
		if !ok {
			return nil
		}

		for _, key := range sectionKey[1:] {
			if v, ok := cruval.(map[string]interface{}); !ok {
				return nil
			} else if cruval, ok = v[key]; !ok {
				return nil
			}
		}
		return cruval
	} else {
		if v, ok := this.data[key]; ok {
			return v
		}
	}
	return nil
}

func (this *JsonConfigContainer) Bool(key string) (bool, error) {
	val := this.getdata(key)
	if val != nil {
		if v, ok := val.(bool); ok {
			return v, nil
		} else {
			return false, errors.New("Not bool value")
		}
	} else {
		return false, errors.New("Not exist key: " + key)
	}
	return false, nil
}

func (this *JsonConfigContainer) Int(key string) (int, error) {
	val := this.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		} else {
			return 0, errors.New("Not int value")
		}
	} else {
		return 0, errors.New("Not exist key:" + key)
	}
	return 0, nil
}

func (this *JsonConfigContainer) Int64(key string) (int64, error) {
	val := this.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		} else {
			return 0, errors.New("Not int64 value")
		}
	} else {
		return 0, errors.New("Not exist key:" + key)
	}
	return 0, nil
}

func (this *JsonConfigContainer) Float(key string) (float64, error) {
	val := this.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		} else {
			return 0.0, errors.New("Not float64 value")
		}
	} else {
		return 0.0, errors.New("Not exist key:" + key)
	}
	return 0.0, nil
}

func (this *JsonConfigContainer) String(key string) string {
	val := this.getdata(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v
		} else {
			return ""
		}
	} else {
		return ""
	}
	return ""
}

func (this *JsonConfigContainer) Set(key, val string) error {
	this.Lock()
	defer this.Unlock()
	this.data[key] = val
	return nil
}

func (this *JsonConfigContainer) DIY(key string) (v interface{}, err error) {
	val := this.getdata(key)
	if val != nil {
		return val, nil
	} else {
		return nil, errors.New("Not exist key.")
	}
	return nil, nil

}

func init() {
	Register(Json, &JsonConfig{})
}
