package config

import (
	"fmt"
)

var (
	Ini  string = "ini"
	Json string = "json"
	Xml  string = "xml"
	Yaml string = "yaml"
)

/**
 * ConfigContainer 定义怎样设置和获取数据
 * @author abram
 */
type ConfigContainer interface {
	Set(key, val string) error
	String(key string) string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DIY(key string) (interface{}, error)
}

/**
 * Config 是个适配器，定义解释配置文件的接口
 * @author abram
 */
type Config interface {
	Parse(key string) (ConfigContainer, error)
}

/**
 * 适配器列表
 * @author abram
 */
var adapters = make(map[string]Config)

/**
 * 注册解释配置文件的适配器
 * @author abram
 * @param name 适配器的名字
 * @param adapter 适配器
 */
func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config:Register adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("config:Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

/**
 * 解释配置文件
 * @author abram
 * @param 适配器名字
 * @param 配置文件的路径
 * @return ConfigContainer,error
 */
func NewConfig(adapterName, fileName string) (ConfigContainer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adapterName %q (forgotten import?)", adapterName)
	}
	return adapter.Parse(fileName)
}
