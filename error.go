package rsp

import (
	"encoding/json"
	"errors"
	"runtime"
	"strconv"
)
type Error struct {
	file string
	line int
	msg  string
}

func (a Error) File() string {
	return a.file
}

func (a Error) Line() int {
	return a.line
}

func (a Error) Error() string {
	return a.msg
}

// 同时打印错误信息和栈信息
func (a Error) StackError() string {
	return a.file + "|" + strconv.Itoa(a.line) + "|" + a.msg
}

// 以json方式打印栈信息
func (a Error) StackJsonError() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"file": a.file,
		"line": a.line,
		"msg":  a.msg,
	})
}

func NewErrMsg(msg string) *Error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	return &Error{
		msg:  msg,
		file: file,
		line: line,
	}
}

func NewErr(err error) *Error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	return &Error{
		msg:  err.Error(),
		file: file,
		line: line,
	}
}

func Err(err error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	return errors.New(file + "|" + strconv.Itoa(line) + "|" + err.Error())
}
