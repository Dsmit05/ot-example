package api

import (
	_ "github.com/Dsmit05/ot-example/service-main/docs"
	"github.com/Dsmit05/ot-example/service-main/internal/api/controllers"
	"github.com/Dsmit05/ot-example/service-main/internal/broker"
	clir "github.com/Dsmit05/ot-example/service-main/internal/client-service-read"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	r *gin.Engine
}

func NewServer(grpcCli *clir.Client, ext *trace.TracerProvider, producer broker.Producer) *Server {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	rm := controllers.NewReadMsg(grpcCli)
	wm := controllers.NewWriteMsg(producer)

	mainG := r.Group("/msg")

	mainG.Use(otelgin.Middleware("service-main", otelgin.WithTracerProvider(ext)))
	{
		mainG.GET("/:id", rm.GetMsg)
		mainG.POST("", wm.PutMsg)
	}

	return &Server{r: r}
}

func (s *Server) Start(addr string) error {
	return s.r.Run(addr)
}
