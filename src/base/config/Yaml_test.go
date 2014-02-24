package config

import (
	"os"
	"testing"
)

var yamlContext = `
"appname": beeapi
"httpport": 8080
"mysqlport": 3600
"PI": 3.1415976
"runmode": dev
"autorender": false
"copyrequestbody": true
`

func TestYaml(t *testing.T) {
	f, err := os.Create("testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(yamlContext)
	if err != nil {
		t.Fatal(err)
	}

	f.Close()
	defer os.Remove("testyaml.conf")
	yamlConf, err := NewConfig("yaml", "testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}

	if yamlConf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}

	if port, err := yamlConf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}

	if pi, err := yamlConf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}

}
