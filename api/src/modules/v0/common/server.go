//Author: Peter Nagy <https://peternagy.ie>
//Since: 07, 2017
//Description: functions extending the server default behavior

package common

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

//HTTPNotFoundHandler - custom handler for not found errors
func HTTPNotFoundHandler(router *fasthttprouter.Router) {
	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte("{\"Src\":\"api/router\",\"Code\":404,\"Desc\":\"Resource not found\"}"))
	}
}

//HTTPNotAllowedHandler - custom handler for not allowed errors
func HTTPNotAllowedHandler(router *fasthttprouter.Router) {
	router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(405)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte("{\"Src\":\"api/router\",\"Code\":405,\"Desc\":\"Method Not Allowed\"}"))
	}
}

//DisableServerFeatures - disable server features
func DisableServerFeatures(router *fasthttprouter.Router) {
	router.HandleMethodNotAllowed = false
}
