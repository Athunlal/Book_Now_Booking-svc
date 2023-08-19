package usermodule

import (
	"context"

	"log"

	"github.com/athunlal/bookNowBooking-svc/pkg/config"
	"github.com/athunlal/bookNowBooking-svc/pkg/usermodule/pb"
	"google.golang.org/grpc"
)

func UserModule(cfg *config.Config) pb.ProfileManagementClient {
	grpcConn, err := grpc.Dial(cfg.Usersvcurl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect to the GRPC Server")
	}
	return pb.NewProfileManagementClient(grpcConn)
}

func GetUserData(c pb.ProfileManagementClient, userid int64) (*pb.UserDataResponse, error) {
	res, err := c.GetUserData(context.Background(), &pb.UserDataRequest{
		Userid: userid,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UserDataResponse{
		Username: res.Username,
	}, err
}
