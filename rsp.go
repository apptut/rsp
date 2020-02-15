package rsp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const MsgSuccess = "Success"
const MsgFailed = "Failed"

func JsonpOK(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSONP(http.StatusOK, jsonBaseFormat(0, msg, data))
}

func JsonpErr(ctx *gin.Context, args ...interface{}) {
	data, msg, code := makeErrData(args...)
	ctx.JSONP(http.StatusOK, jsonBaseFormat(code, msg.Error(), data))
}

func makeOkData(args ...interface{}) (interface{}, string) {
	var data interface{}
	var msg = MsgSuccess

	if len(args) >= 1 {
		data = args[0]
	}

	if len(args) >= 2 {
		msg = args[1].(string)
	}

	return data, msg
}

func makeErrData(args ...interface{}) (interface{}, *Error, int) {
	var data interface{}
	code := 1
	var appErr *Error

	if len(args) >= 1 {
		first := args[0]
		if tmpRsp, ok := first.(*Error); ok {
			appErr = tmpRsp
		} else if tmpErr, ok := first.(error); ok {
			appErr = NewErr(tmpErr)
		} else {
			appErr = NewErrMsg(first.(string))
		}
	} else {
		appErr = NewErrMsg(MsgFailed)
	}

	if len(args) >= 2 {
		data = args[1]
	}

	if len(args) >= 3 {
		code = args[2].(int)
	}

	return data, appErr, code
}

// JsonOk(ctx, data, msg)
func JsonOk(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(0, msg, data))
}

// JsonErr(ctx, msg, data)
func JsonErr(ctx *gin.Context, args ...interface{}) {
	data, msg, code := makeErrData(args...)

	errMsgs := ctx.Errors
	ctx.Errors = append(errMsgs, &gin.Error{
		Err:  errors.New(msg.Error()),
		Type: gin.ErrorTypePrivate,
		Meta: map[string]interface{}{
			"file": msg.File(),
			"line": msg.Line(),
		},
	})

	ctx.JSON(http.StatusOK, jsonBaseFormat(code, msg.Error(), data))
}

func jsonBaseFormat(code int, msg string, data interface{}) interface{} {
	var baseData = make(map[string]interface{})
	baseData["errno"] = code
	baseData["msg"] = msg
	baseData["data"] = data
	return baseData
}
