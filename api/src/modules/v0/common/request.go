//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: collection of methods for handling requests
package common

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"strconv"
)

var (
	corsAllowHeaders     = "Authorization,cache-control,content-type,x-xsrf-token"
	corsAllowMethods     = "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

//Decode bytes into struct's (unmarshal)
func Decode(data []byte, v interface{}) error {
	return ffjson.Unmarshal(data, &v)
}

//DecodeBody decodes request json body
func DecodeBody(ctx *fasthttp.RequestCtx, v interface{}) error {
	return ffjson.Unmarshal(ctx.PostBody(), &v)
}

//WithCORS - enable CORS support
func WithCORS(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		h(ctx)
	})
}

//HandleOPTIONSRequest returns an empty response with headers
func HandleOPTIONSRequest(ctx *fasthttp.RequestCtx) {
}

//GetUserValueIntOrDefault - get request user (parameter) value or default
func GetUserValueIntOrDefault(ctx *fasthttp.RequestCtx, key string, d int) int {
	ps, ok := ctx.UserValue(key).(string)
	if !ok {
		return d
	}

	n, err := strconv.Atoi(ps)
	if err != nil {
		return d
	}

	return n
}
