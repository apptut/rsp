package rsp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const MsgSuccess = "Success"
const MsgFailed = "Failed"

func JsonpOk(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSONP(http.StatusOK, jsonBaseFormat(0, msg, data))
}

func JsonpErr(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(1, msg, data))
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

func makeErrData(args ...interface{}) (interface{}, string) {
	var data interface{}
	var msg = MsgFailed

	if len(args) >= 1 {
		msg = args[0].(string)
	}

	if len(args) >= 2 {
		data = args[1]
	}
	return data, msg
}

// JsonOk(ctx, data, msg)
func JsonOk(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(0, msg, data))
}

// JsonErr(ctx, msg, data)
func JsonErr(ctx *gin.Context, args ...interface{}) {
	data, msg := makeErrData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(1, msg, data))
}

func jsonBaseFormat(code int, msg string, data interface{}) interface{} {
	var baseData = make(map[string]interface{})
	baseData["errno"] = code
	baseData["msg"] = msg
	baseData["data"] = data
	return baseData
}

// JsonError(ctx, code, msg)
func JsonError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusOK, jsonErrorFormat(code, msg))
}

func jsonErrorFormat(code int, msg string) interface{} {
	var errData = make(map[string]interface{})
	errData["errno"] = code
	errData["msg"] = msg
	return errData
}
