package handler

import (
	"context"
	"net/http"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/pkg/api/response"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type resp = response.Response

type respErr struct {
	Success bool   `json:"success" example:"false"`
	Msg     string `json:"msg" example:"error message"`
	// Obj     any    `json:"obj,omitempty"`
}

type respSuc struct {
	Success bool   `json:"success" example:"true"`
	Msg     string `json:"msg" example:""`
	Obj     any    `json:"obj,omitempty"`
}

type respSucNoObj struct {
	Success bool   `json:"success" example:"true"`
	Msg     string `json:"msg" example:""`
}

func writeResponse(c *gin.Context, status int, success bool, msg string, obj interface{}) {
	c.JSON(status, resp{Success: success, Msg: msg, Obj: obj})
}

/* ----  Predefined response status ---- */

func writeSuccess(c *gin.Context, status int, msg string, obj interface{}) {
	c.JSON(status, resp{Success: true, Msg: msg, Obj: obj})
}

func writeFailure(c *gin.Context, status int, msg string, obj interface{}) {
	c.JSON(status, resp{Success: false, Msg: msg, Obj: obj})
}

/* ---- Predefined both HTTP and response statuses ---- */
const msgSuccess = "Operation succeeded!"

func writeOK(c *gin.Context) {
	c.JSON(http.StatusOK, resp{Success: true, Obj: nil, Msg: msgSuccess})
}

func writeObj(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, resp{Success: true, Obj: obj, Msg: msgSuccess})
}

func writeBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, resp{Success: false, Msg: msg})
}

func writeNotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, resp{Success: false, Msg: msg})
}

func writeServerInternal(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, resp{Success: false, Msg: msg})
}

func prepareTools(c *gin.Context, op string) (logger zerolog.Logger, ctx context.Context) {
	return log.With().Str("op", op).Str("request_id", requestid.Get(c)).Logger(),
		c.Request.Context()
}
