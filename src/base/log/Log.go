package log

import (
	"fmt"
	"sync"
)

const (
	Console = "console" //控制台输出
	File    = "file"    // 文件输出
)

const (
	Trace = iota
	Debug
	Info
	Warn
	Error
	Critical
)

type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

type loggerType func() LoggerInterface

var adapters = make(map[string]loggerType)

func Register(name string, log loggerType) {
	if log == nil {
		panic("Log: Register provide is nil")
	}

	if _, dup := adapters[name]; dup {
		panic("Logs:Register called twice for provider " + name)
	}
	adapters[name] = log
}

type BeeLogger struct {
	lock    sync.Mutex
	level   int
	msg     chan *logMsg
	outputs map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

func NewLogger(channellen int64) *BeeLogger {
	bl := new(BeeLogger)
	bl.msg = make(chan *logMsg, channellen)
	bl.outputs = make(map[string]LoggerInterface)
	go bl.StartLogger()
	return bl
}

func (bl *BeeLogger) SetLogger(adapterName string, config string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if log, ok := adapters[adapterName]; ok {
		lg := log()
		lg.Init(config)
		bl.outputs[adapterName] = lg
		return nil
	} else {
		return fmt.Errorf("Logs:unknow adaptername %q (forgotten Register?", adapterName)
	}
}

func (bl *BeeLogger) DelLogger(adapterName string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if lg, ok := bl.outputs[adapterName]; ok {
		lg.Destroy()
		delete(bl.outputs, adapterName)
		return nil
	} else {
		return fmt.Errorf("Logs:unknow adapterName %q (forgotten Regester?)", adapterName)
	}
}

func (bl *BeeLogger) writeMsg(level int, msg string) error {
	if bl.level > level {
		return nil
	}

	lm := new(logMsg)
	lm.level = level
	lm.msg = msg
	bl.msg <- lm
	return nil
}

func (bl *BeeLogger) Setlevel(level int) {
	bl.level = level
}

func (bl *BeeLogger) StartLogger() {
	for {
		select {
		case bm := <-bl.msg:
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		}
	}
}

func (bl *BeeLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	bl.writeMsg(Trace, msg)
}

func (bl *BeeLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	bl.writeMsg(Debug, msg)
}

func (bl *BeeLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	bl.writeMsg(Info, msg)
}

func (bl *BeeLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	bl.writeMsg(Warn, msg)
}

func (bl *BeeLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	bl.writeMsg(Error, msg)
}

func (bl *BeeLogger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	bl.writeMsg(Critical, msg)
}

func (bl *BeeLogger) Flush() {
	for _, l := range bl.outputs {
		l.Flush()
	}
}

func (bl *BeeLogger) Close() {
	for {
		if len(bl.msg) > 0 {
			bm := <-bl.msg
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		} else {
			break
		}
	}
	for _, l := range bl.outputs {
		l.Flush()
		l.Destroy()
	}
}
