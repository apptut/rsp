package rsp

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
	"regexp"
)

const MsgSuccess = "Success"
const MsgFailed = "Failed"

func JsonpOK(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	safeJsonP(ctx, http.StatusOK, jsonBaseFormat(0, msg, data))
}

func JsonpErr(ctx *gin.Context, args ...interface{}) {
	data, msg, code := makeErrData(args...)
	safeJsonP(ctx, http.StatusOK, jsonBaseFormat(code, msg.Error(), data))
}

func safeJsonP(ctx *gin.Context, code int, obj interface{}) {
	callback := ctx.DefaultQuery("callback", "")
	if callback == "" {
		ctx.Render(code, render.JSON{Data: obj})
		return
	}

	if !regexp.MustCompile(`^[\w-.]{1,64}$`).MatchString(callback) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
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

// JsonOk (ctx, data, msg)
func JsonOk(ctx *gin.Context, args ...interface{}) {
	data, msg := makeOkData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(0, msg, data))
}

// JsonErr (ctx, msg, data)
func JsonErr(ctx *gin.Context, args ...interface{}) {
	data, msg, code := makeErrData(args...)
	ctx.JSON(http.StatusOK, jsonBaseFormat(code, msg.Error(), data))
}

func jsonBaseFormat(code int, msg string, data interface{}) interface{} {
	var baseData = make(map[string]interface{})
	baseData["errno"] = code
	baseData["msg"] = msg
	baseData["data"] = data
	return baseData
}
