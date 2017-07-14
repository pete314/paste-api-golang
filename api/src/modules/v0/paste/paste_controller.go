//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: --
package paste

import (
	"../common"
	valid "github.com/asaskevich/govalidator"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	basePath = "/v0/paste"
	idSub    = "/:id"
)

type Controller struct {
	m *Model
}

type HTTPController interface {
	EnableHTTPMethods(router *fasthttprouter.Router)
}

func NewController() *Controller {
	return &Controller{m: NewModel()}
}

//EnableHTTPMethods - enable http methods for this controller
func (c *Controller) EnableHTTPMethods(router *fasthttprouter.Router) {
	router.OPTIONS(basePath+"/*filepath", common.WithCORS(common.HandleOPTIONSRequest))
	router.GET(basePath+idSub, common.WithCORS(c.HandleGETRequest))
	router.PUT(basePath, common.WithCORS(c.HandlePUTRequest))
}

func (c *Controller) HandleGETRequest(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.UserValue("id").(string)

	if !valid.IsUUIDv4(id) {
		common.ErrorJsonResponse(ctx, 400, "api/paste/request/parse", "Invalid ID in request")
		return
	}

	n, err := c.m.GetEntry(id)
	if err != nil {
		common.ErrorJsonResponse(ctx, 500, "api/paste/storage/db", "Error while retrieving entry")
		return
	}

	common.SuccessJsonResponse(ctx, 200, n)
}

func (c *Controller) HandlePUTRequest(ctx *fasthttp.RequestCtx) {
	var body *Note
	if err := common.DecodeBody(ctx, &body); err != nil {
		common.CheckError("paste::handlePUTRequest", err, false)
		common.ErrorJsonResponse(ctx, 400, "api/paste/request/parse", "Failed to parse request body")
		return
	}

	if ok, err := valid.ValidateStruct(body); !ok || err != nil {
		common.ErrorJsonResponse(ctx, 400, "api/paste/request/validate", "Value field required, non empty "+err.Error())
		return
	}

	InitNoteDefaults(body)
	n, err := c.m.StoreEntry(body)
	if err != nil {
		common.ErrorJsonResponse(ctx, 500, "api/paste/storage/db", "Internal error, while saving entry")
		return
	}

	common.SuccessJsonResponse(ctx, 200, n)
}
