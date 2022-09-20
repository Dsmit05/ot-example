package api

import (
	"context"

	"github.com/Dsmit05/ot-example/service-read/internal/repository"
	pb "github.com/Dsmit05/ot-example/service-read/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerReader struct {
	pb.UnimplementedMsgReaderServer
	rep repository.ReposytoryI
}

func NewServerReader(rep repository.ReposytoryI) *ServerReader {
	return &ServerReader{rep: rep}
}

func (s *ServerReader) ReadMsg(ctx context.Context, msgReq *pb.ReadMsgRequest) (*pb.ReadMsgResponse, error) {
	msg, err := s.rep.GetMsg(ctx, msgReq.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.ReadMsgResponse{
		Id:  msgReq.GetId(),
		Msg: msg,
	}, nil
}
