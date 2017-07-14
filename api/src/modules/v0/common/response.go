//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: --
package common

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

//ErrorBody - Error message structure
type ErrorBody struct {
	Src  string
	Code int
	Desc string
}

//SuccessBody - Response serializable body
type SuccessBody struct {
	Success bool
	Result  interface{}
}

//IDUriResult - Create and update type response
type IDUriResult struct {
	ID  string `json:"id"`
	Uri string `json:"uri"`
}

//Encode - Encode bytes into struct's (marshal)
func Encode(i interface{}) ([]byte, error) {
	return ffjson.Marshal(&i)
}

//JsonResponse - Create response
func JsonResponse(ctx *fasthttp.RequestCtx, status int, body interface{}) {
	ctx.SetContentType("application/json")
	if buff, err := Encode(&body); err == nil {
		ctx.SetBody(buff)
		ctx.SetStatusCode(status)
	} else {
		ctx.SetBody(NewErrorBuffer("Api/Common/Response/Serialize", "Error while serializing response", 500))
		CheckError("Response: Error while serializing response", err, false)
	}
}

//SuccessJsonResponse - add success structure json response
func SuccessJsonResponse(ctx *fasthttp.RequestCtx, status int, body interface{}) {
	JsonResponse(ctx, status, NewSuccessBody(true, body))
}

//ErrorJsonResponse - add error structure json response
func ErrorJsonResponse(ctx *fasthttp.RequestCtx, status int, src, desc string) {
	JsonResponse(ctx, status, NewError(src, desc, status))
}

//NewError - Get new error Body struct
func NewError(src, desc string, code int) *ErrorBody {
	return &ErrorBody{Src: src, Code: code, Desc: desc}
}

//NewErrorBuffer - Get new Error buffer
func NewErrorBuffer(src, desc string, code int) []byte {
	eBody := NewError(src, desc, code)
	bytes, _ := Encode(eBody)
	return bytes
}

//NewResponseBody - Get response body
func NewSuccessBody(success bool, result interface{}) *SuccessBody {
	return &SuccessBody{Success: success, Result: result}
}
