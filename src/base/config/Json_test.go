package config

import (
	"os"
	"testing"
)

var jsonContext = `{
	"appname":"beeapi",
	"httpport":8080,
	"mysqlport":3600,
	"PI":3.1415976,
	"runmode":"dev",
	"autorender":false,
	"copyrequestbody":true,
	"database":{
		"host":"host",
		"port":"port",
		"database":"database",
		"username":"username",
		"password":"password",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info"
		}
	}
}`

func TestJson(t *testing.T) {
	f, err := os.Create("testjson.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(jsonContext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	f.Close()
	defer os.Remove("testjson.conf")

	jsonConf, err := NewConfig("json", "testjson.conf")
	if err != nil {
		t.Fatal(err)
	}

	if jsonConf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
}
