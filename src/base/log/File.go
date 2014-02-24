package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileLogWriter struct {
	*log.Logger
	mw               *MuxWriter
	FileName         string `json:"fileName"`
	MaxLines         int    `json:"maxLines"`
	maxLinesCurLines int
	MaxSize          int `json:"maxSize"`
	maxSizeCurSize   int
	Daily            bool  `json:"daily"`
	MaxDays          int64 `json:"maxDays"`
	dailyOpendate    int
	Rotate           bool `json:"rotate"`
	startLock        sync.Mutex
	Level            int `json:"level`
}

type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

func NewFileWriter() LoggerInterface {
	w := &FileLogWriter{
		FileName: "",
		MaxLines: 1000000,
		MaxSize:  1 << 28, //256 MB
		Daily:    true,
		MaxDays:  7,
		Rotate:   true,
		Level:    Trace,
	}
	w.mw = new(MuxWriter)
	w.Logger = log.New(w.mw, "", log.Ldate|log.Ltime)
	return w
}

func (w *FileLogWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), w)
	if err != nil {
		return err
	}

	if len(w.FileName) == 0 {
		return errors.New("jsonConfig mush have fileName")
	}

	err = w.StartLogger()
	return err
}

func (w *FileLogWriter) StartLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}

	w.mw.SetFd(fd)
	err = w.initFd()
	if err != nil {
		return err
	}
	return nil
}

func (w *FileLogWriter) createLogFile() (*os.File, error) {
	fd, err := os.OpenFile(w.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (w *FileLogWriter) initFd() error {
	fd := w.mw.fd
	info, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("Get stat err:%s\n", err)
	}

	w.maxSizeCurSize = int(info.Size())
	w.dailyOpendate = time.Now().Day()
	if info.Size() > 0 {
		content, err := ioutil.ReadFile(w.FileName)
		if err != nil {
			return err
		}

		w.maxLinesCurLines = len(strings.Split(string(content), "\n"))
	} else {
		w.maxLinesCurLines = 0
	}
	return nil
}

func (w *FileLogWriter) doCheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()

	if (w.MaxLines > 0 && w.maxLinesCurLines >= w.MaxLines) ||
		(w.MaxSize > 0 && w.maxSizeCurSize >= w.MaxSize) ||
		(w.Daily && time.Now().Day() != w.dailyOpendate) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.FileName, err)
			return
		}
	}
	w.maxLinesCurLines++
	w.maxSizeCurSize += size
}

func (w *FileLogWriter) WriteMsg(msg string, level int) error {
	if level < w.Level {
		return nil
	}

	n := 24 + len(msg)
	w.doCheck(n)
	w.Logger.Println(msg)
	return nil
}

func (w *FileLogWriter) DoRotate() error {
	_, err := os.Lstat(w.FileName)
	if err == nil {
		num := 1
		fname := ""
		for ; err == nil && num <= 999; num++ {
			fname = w.FileName + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), num)
			_, err = os.Lstat(fname)
		}

		if err == nil {
			return fmt.Errorf("Rotate: Connot find free log number to rename %s\n", w.FileName)
		}

		w.mw.Lock()
		defer w.mw.Unlock()
		fd := w.mw.fd
		fd.Close()

		err = os.Rename(w.FileName, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s\n", err)
		}

		err = w.StartLogger()
		if err != nil {
			return fmt.Errorf("Rotate:%s\n", err)
		}

		go w.deleteOldLog()
	}
	return nil
}

func (w *FileLogWriter) deleteOldLog() {
	dir := filepath.Dir(w.FileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.MaxDays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.FileName)) {
				os.Remove(path)
			}
		}
		return nil
	})
}

func (w *FileLogWriter) Destroy() {
	w.mw.fd.Close()
}

func (w *FileLogWriter) Flush() {
	w.mw.fd.Sync()
}

func init() {
	Register(File, NewFileWriter)
}
