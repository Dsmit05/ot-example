package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Dsmit05/ot-example/service-main/internal/models"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
)

type grpcClient interface {
	ReadMsg(ctx context.Context, id int64) (models.Message, error)
}

type ReadMsg struct {
	cli grpcClient
}

func NewReadMsg(cli grpcClient) *ReadMsg {
	return &ReadMsg{cli: cli}
}

// @Summary ReadMsg
// @Tags msg
// @Description get msg from read service
// @ID get-msg
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Message
// @Failure 400 {object} string "error"
// @Router /msg/{id} [GET]
func (r *ReadMsg) GetMsg(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	// add info in Tags span
	span.SetAttributes(attribute.String("controller", "GetMsg"))

	msgID := c.Param("id")

	id, err := strconv.Atoi(msgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	// example need validate and check time result:

	tr := otel.Tracer("validate")
	ctx, spanValid := tr.Start(ctx, "validate")

	time.Sleep(time.Millisecond * 15)

	spanValid.End()

	// ---------------------------------------------

	msg, err := r.cli.ReadMsg(ctx, int64(id))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// need use mapping status code:
		// https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md
		c.JSON(int(st.Code()), st.Message())

		return
	}

	// add info in Events span
	span.AddEvent("msg", trace.WithAttributes(attribute.Int64("id", msg.ID), attribute.String("msg", msg.Msg)))

	c.JSON(http.StatusOK, msg)
}
