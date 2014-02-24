package config

import "os"
import "testing"

var iniContext = `
;comment one
#comment two
appname = beeapi
httpport = 8080
mysqlport = 3600
PI=3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
[demo]
key1="asta"
key2="xie"
CaseInsensitive=true
`

func TestIni(t *testing.T) {
	f, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(iniContext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	f.Close()

	defer os.Remove("testini.conf")
	iniConf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}

	if iniConf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}

}
