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
	"log"
	"os"
)

type QConsoleWriter struct {
	*log.Logger
}

func NewQConsoleWriter() QLogWriter {
	return &QConsoleWriter{Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

func (this *QConsoleWriter) Init() {

}

func (this *QConsoleWriter) Destroy() {

}

func (this *QConsoleWriter) WriteMsg(msg string, level QLogLevel) {
	this.Logger.Println(msg)
}
func (this *QConsoleWriter) Flush() {

}

func init() {
	Register("console", NewQConsoleWriter)
}
