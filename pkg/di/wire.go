//go:build wireinject
// +build wireinject

package di

import (
	"github.com/athunlal/bookNowBooking-svc/pkg/api"
	"github.com/athunlal/bookNowBooking-svc/pkg/api/handler"
	"github.com/athunlal/bookNowBooking-svc/pkg/config"
	"github.com/athunlal/bookNowBooking-svc/pkg/db"
	"github.com/athunlal/bookNowBooking-svc/pkg/repository"
	"github.com/athunlal/bookNowBooking-svc/pkg/usecase"
	"github.com/google/wire"
)

func InitApi(cfg config.Config) (*api.ServerHttp, error) {
	wire.Build(
		db.ConnectDataBase,
		repository.NewTrainRepo,
		usecase.NewBookingUseCase,
		handler.NewBookingHandler,
		api.NewServerHttp)
	return &api.ServerHttp{}, nil
}
