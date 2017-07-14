//Author: Peter Nagy <https://peternagy.ie>
//Since: 07, 2017
//Description: controller for authorization requests
package auth

import (
	"../common"
	valid "github.com/asaskevich/govalidator"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"sync"
)

const (
	basePath = "/v0/auth"
	idSub    = "/:id"
)

var (
	once sync.Once
)

type HTTPController interface {
	EnableHTTPMethods(router *fasthttprouter.Router)
}

type Controller struct {
	m *Model
}

//NewController - get new instance of this
func NewController() *Controller {
	once.Do(func() {
		valid.SetFieldsRequiredByDefault(true)
	})
	return &Controller{}
}

//EnableHTTPMethods - enable http methods for this controller
func (c *Controller) EnableHTTPMethods(router *fasthttprouter.Router) {
	router.OPTIONS(basePath+"/*filepath", common.WithCORS(common.HandleOPTIONSRequest))
	router.GET(basePath+idSub, common.WithCORS(c.HandleGETRequest))
	router.POST(basePath+"/register", common.WithCORS(c.HandlePOSTRequest))
	router.POST(basePath+"/authorize", common.WithCORS(c.HandlePOSTRequest))
}

func (c *Controller) HandleGETRequest(ctx *fasthttp.RequestCtx) {
	common.SuccessJsonResponse(ctx, 200, true)
}

func (c *Controller) HandlePOSTRequest(ctx *fasthttp.RequestCtx) {
	var u *User

	if err := common.DecodeBody(ctx, u); err != nil {
		common.CheckError("auth_controller:HandlePOSTRequest", err, false)
		common.ErrorJsonResponse(ctx, 400, "api/auth/request/parse", "Error while parsing body")
		return
	}

	if ok, err := valid.ValidateStruct(u); !ok || err != nil {
		common.ErrorJsonResponse(ctx, 400, "api/auth/request/validate", "Invalid or missing data sent "+err.Error())
		return
	}

	ok, err := c.m.CreateUser(u)
	if err != nil {
		common.ErrorJsonResponse(ctx, 400, "api/auth/storage/db", "Error storing user, "+err.Error())
		return
	}

	common.SuccessJsonResponse(ctx, 200, ok)
}

func (c *Controller) HandlePOSTAuthorizeRequest(ctx *fasthttp.RequestCtx) {
	var uc *UserCredentials

	if err := common.DecodeBody(ctx, uc); err != nil {
		common.CheckError("auth_controller:HandlePOSTRequest", err, false)
		common.ErrorJsonResponse(ctx, 400, "api/auth/request/parse", "Error while parsing body")
		return
	}

	if ok, err := valid.ValidateStruct(uc); !ok || err != nil {
		common.ErrorJsonResponse(ctx, 400, "api/auth/request/validate", "Invalid or missing data sent "+err.Error())
		return
	}

}
