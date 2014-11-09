// Copyright 2014 layeka Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package qlog

import (
	"fmt"
	"path"
	"runtime"
)

type QLogLevel int

const (
	TRACE QLogLevel = 1 << iota
	DEBUG
	INFO
	ERROR
	WARN
	CRITICAL
)

type QLogWriter interface {
	Init()
	Destroy()
	WriteMsg(msg string, level QLogLevel)
	Flush()
}
type QLogWriterFactory func() QLogWriter

var (
	adapters = make(map[string]QLogWriterFactory)
)

func Register(name string, factory QLogWriterFactory) {
	if factory == nil {
		panic("qlog: Register factory is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("qlog: Register called twice for factory " + name)
	}
	adapters[name] = factory
}

type qlogbuffer struct {
	msg   string
	level QLogLevel
}
type QLoger struct {
	writers        map[string]QLogWriter
	buffers        chan *qlogbuffer
	level          QLogLevel
	enableFuncCall bool
	funcCallDepth  int
}

func NewQLoger(bufferlen int) *QLoger {
	qloger := new(QLoger)
	qloger.enableFuncCall = true
	qloger.funcCallDepth = 2
	qloger.buffers = make(chan *qlogbuffer, bufferlen)
	qloger.writers = make(map[string]QLogWriter)
	go qloger.startLoger()
	return qloger
}
func (this *QLoger) startLoger() {
	for {
		select {
		case buffer := <-this.buffers:
			for _, writer := range this.writers {
				writer.WriteMsg(buffer.msg, buffer.level)
			}
		}
	}
}
func (this *QLoger) writeMsg(msg string, level QLogLevel) {
	buffer := new(qlogbuffer)
	buffer.level = level
	if this.enableFuncCall {
		_, file, line, ok := runtime.Caller(this.funcCallDepth)
		if ok {
			_, filename := path.Split(file)
			buffer.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
		} else {
			buffer.msg = msg
		}
	} else {
		buffer.msg = msg
	}
	this.buffers <- buffer
}
func (this *QLoger) Trace(format string, v ...interface{}) {
	if TRACE&this.level == TRACE {
		this.writeMsg(fmt.Sprintf("[T] "+format, v...), TRACE)
	}
}
func (this *QLoger) Debug(format string, v ...interface{}) {
	if DEBUG&this.level == DEBUG {
		this.writeMsg(fmt.Sprintf("[D] "+format, v...), DEBUG)
	}
}
func (this *QLoger) Info(format string, v ...interface{}) {
	if INFO&this.level == INFO {
		this.writeMsg(fmt.Sprintf("[I] "+format, v...), INFO)
	}
}
func (this *QLoger) Error(format string, v ...interface{}) {
	if ERROR&this.level == ERROR {
		this.writeMsg(fmt.Sprintf("[E] "+format, v...), ERROR)
	}
}
func (this *QLoger) Warn(format string, v ...interface{}) {
	if WARN&this.level == WARN {
		this.writeMsg(fmt.Sprintf("[W] "+format, v...), WARN)
	}
}
func (this *QLoger) Critical(format string, v ...interface{}) {
	if CRITICAL&this.level == CRITICAL {
		this.writeMsg(fmt.Sprintf("[C] "+format, v...), CRITICAL)
	}
}
