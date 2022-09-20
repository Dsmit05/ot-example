package controllers

import (
	"net/http"

	"github.com/Dsmit05/ot-example/service-main/internal/broker"
	"github.com/Dsmit05/ot-example/service-main/internal/models"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type WriteMsg struct {
	cli broker.Producer
}

func NewWriteMsg(cli broker.Producer) *WriteMsg {
	return &WriteMsg{cli: cli}
}

// @Summary WriteMsg
// @Tags msg
// @Description put msg in write service
// @ID write-msg
// @Accept json
// @Produce json
// @Param input body models.Message true "credentials"
// @Success 200 {string} string
// @Failure 400 {object} string "error"
// @Router /msg [POST]
func (r *WriteMsg) PutMsg(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	span.SetAttributes(attribute.String("controller", "PutMsg"))

	var inputData models.Message
	if err := c.ShouldBindJSON(&inputData); err != nil {
		span.SetStatus(1, "not convert input data")
		c.JSON(http.StatusInternalServerError, err.Error())

		return
	}

	if err := r.cli.SendMsg(c.Request.Context(), inputData); err != nil {
		span.SetStatus(1, "error from broker")
		c.JSON(http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, "put msg in broker")
}
