package client_service_read

import (
	"context"

	"github.com/Dsmit05/ot-example/service-main/internal/models"
	pb "github.com/Dsmit05/ot-example/service-read/pkg/api"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcCLi pb.MsgReaderClient
}

func NewClient(targetURL string) (*Client, error) {
	conns, err := grpc.Dial(targetURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial() error")
	}

	msgReaderClient := pb.NewMsgReaderClient(conns)

	return &Client{grpcCLi: msgReaderClient}, nil
}

func (c *Client) ReadMsg(ctx context.Context, id int64) (models.Message, error) {
	var resp models.Message

	grpcResponse, err := c.grpcCLi.ReadMsg(ctx, &pb.ReadMsgRequest{Id: id})
	if err != nil {
		return resp, errors.Wrap(err, "grpc ReadMsg() error")
	}

	resp.ID = grpcResponse.GetId()
	resp.Msg = grpcResponse.GetMsg()

	return resp, nil
}
