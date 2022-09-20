package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Dsmit05/ot-example/service-read/internal/repository"
	pb "github.com/Dsmit05/ot-example/service-read/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	sr         *ServerReader
}

func NewServer(rep repository.ReposytoryI) *Server {
	sr := NewServerReader(rep)

	return &Server{sr: sr}
}

func (s *Server) StartGRPC(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "net.Listen() error")
	}

	newGrpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(time.Second*5),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))

	s.grpcServer = newGrpcServer
	pb.RegisterMsgReaderServer(s.grpcServer, s.sr)

	if err = s.grpcServer.Serve(listener); err != nil {
		return errors.Wrap(err, "grpcServer.Serve() error")
	}

	return nil
}

func (s *Server) StartHTTP(addrHTTP, addrGRPC string) error {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterMsgReaderHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
	if err != nil {
		return errors.Wrap(err, "RegisterMsgReaderHandlerFromEndpoint() error")
	}

	newHTTPServer := &http.Server{
		Addr:    addrHTTP,
		Handler: mux,
	}

	s.httpServer = newHTTPServer

	if err := s.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Print("HTTP server", "Server closed under request")
		} else {
			return errors.Wrap(err, "http.ListenAndServe() error")
		}
	}

	return nil
}

func (s *Server) Stop() {
	stop := make(chan bool)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func() {
		s.httpServer.SetKeepAlivesEnabled(false)

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown failed", err)
			return
		}

		s.grpcServer.GracefulStop()

		stop <- true
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()

		log.Fatal(ctx.Err())
	case <-stop:
		log.Print("Server", "Server closed under request")
	}
}
