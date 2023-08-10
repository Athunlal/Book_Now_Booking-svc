package api

import (
	"log"
	"net"

	"github.com/athunlal/bookNowBooking-svc/pkg/api/handler"
	"github.com/athunlal/bookNowBooking-svc/pkg/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServerHttp struct {
	Engine *gin.Engine
}

func NewGrpcServer(BookingHandler *handler.BookingHandler, grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalln("Failed to listen to the GRPC Port", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBookingManagementServer(grpcServer, BookingHandler)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Could not serve the GRPC Server: ", err)
	}
}

func NewServerHttp(BookingHandler *handler.BookingHandler) *ServerHttp {
	engine := gin.New()
	go NewGrpcServer(BookingHandler, "8893")

	engine.Use(gin.Logger())
	return &ServerHttp{
		Engine: engine,
	}
}

func (ser *ServerHttp) Start() {
	ser.Engine.Run(":9001")
}
