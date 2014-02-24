package config

import (
	"os"
	"testing"
)

var xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>beeapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
</config>
`

func TestXML(t *testing.T) {
	file, err := os.Create("testXml.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.WriteString(xmlcontext)

	if err != nil {
		file.Close()
		t.Fatal(err)
	}

	file.Close()
	defer os.Remove("testXml.conf")
	xmlConf, err := NewConfig("xml", "testXml.conf")
	if err != nil {
		t.Fatal(err)
	}

	if xmlConf.String("appname") != "beeapi" {
		t.Fatal("appname not equal beeapi")
	}
}
