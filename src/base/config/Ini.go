package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	DEFAULT_SECTION = "default"
	bNumComment     = []byte{'#'}
	bSemComment     = []byte{';'}
	bEmpty          = []byte{}
	bEqual          = []byte{'='}
	bDQuote         = []byte{'"'}
	sectionStart    = []byte{'['}
	sectionEnd      = []byte{']'}
)

type IniConfig struct {
}

type IniConfigContainer struct {
	fileName       string
	data           map[string]map[string]string
	sectionComment map[string]string
	keyComment     map[string]string
	sync.RWMutex
}

func (this *IniConfig) Parse(name string) (ConfigContainer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	cfg := &IniConfigContainer{
		file.Name(),
		make(map[string]map[string]string),
		make(map[string]string),
		make(map[string]string),
		sync.RWMutex{},
	}
	cfg.Lock()

	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)
	section := DEFAULT_SECTION

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}

		if bytes.Equal(line, bEmpty) {
			continue
		}

		line = bytes.TrimSpace(line)

		var bComment []byte
		switch {
		case bytes.HasPrefix(line, bNumComment):
			bComment = bNumComment
		case bytes.HasPrefix(line, bSemComment):
			bComment = bSemComment
		}

		if bComment != nil {
			line = bytes.TrimLeft(line, string(bComment))
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}

		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			section = string(line[1 : len(line)-1])
			section = strings.ToLower(section)
			if comment.Len() > 0 {
				cfg.sectionComment[section] = comment.String()
				comment.Reset()
			}

			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
		} else {
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}

			keyval := bytes.SplitN(line, bEqual, 2)
			val := bytes.TrimSpace(keyval[1])
			if bytes.HasPrefix(val, bDQuote) {
				val = bytes.Trim(val, `"`)
			}

			key := string(bytes.TrimSpace(keyval[0]))
			key = strings.ToLower(key)
			cfg.data[section][key] = string(val)

			if comment.Len() > 0 {
				cfg.keyComment[section+"."+key] = comment.String()
				comment.Reset()
			}
		}
	}
	return cfg, nil
}

func (this *IniConfigContainer) getdata(key string) string {
	this.RLock()
	defer this.RUnlock()

	if len(key) == 0 {
		return ""
	}

	var section, k string
	key = strings.ToLower(key)
	sectionKey := strings.Split(key, "::")
	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionKey[0]
	}

	if v, ok := this.data[section]; ok {
		if vv, o := v[k]; o {
			return vv
		}
	}
	return ""
}

func (this *IniConfigContainer) DIY(key string) (v interface{}, err error) {
	key = strings.ToLower(key)
	if v, ok := this.data[key]; ok {
		return v, nil
	}

	return v, errors.New("Key not find")
}

func init() {
	Register(Ini, &IniConfig{})
}

func (this *IniConfigContainer) Bool(key string) (bool, error) {
	key = strings.ToLower(key)
	return strconv.ParseBool(this.getdata(key))
}

func (this *IniConfigContainer) Int(key string) (int, error) {
	key = strings.ToLower(key)
	return strconv.Atoi(this.getdata(key))
}

func (this *IniConfigContainer) Int64(key string) (int64, error) {
	key = strings.ToLower(key)
	return strconv.ParseInt(this.getdata(key), 10, 64)
}

func (this *IniConfigContainer) Float(key string) (float64, error) {
	key = strings.ToLower(key)
	return strconv.ParseFloat(this.getdata(key), 64)
}

func (this *IniConfigContainer) String(key string) string {
	key = strings.ToLower(key)
	return this.getdata(key)
}

func (this *IniConfigContainer) Set(key, val string) error {
	this.Lock()
	defer this.Unlock()
	if len(key) == 0 {
		return errors.New("Key is empty")
	}

	var section, k string
	key = strings.ToLower(key)
	sectionKey := strings.Split(key, "::")
	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionKey[0]
	}

	if _, ok := this.data[section]; !ok {
		this.data[section] = make(map[string]string)
	}
	this.data[section][k] = val
	return nil
}
