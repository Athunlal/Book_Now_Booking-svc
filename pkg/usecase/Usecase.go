package usecase

import (
	"context"
	"fmt"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	usecase "github.com/athunlal/bookNowBooking-svc/pkg/usecase/interface"
)

type BookingUseCase struct {
	Repo interfaces.BookingRepo
}

// SearchTrain implements interfaces.TrainUseCase.
func (use *BookingUseCase) SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	routeData, err := use.Repo.FindRouteId(ctx, searcheData)
	route := domain.Train{
		Route: routeData.RouteID,
	}
	routid := domain.Route{
		RouteId: routeData.RouteID,
	}
	res1, err := use.Repo.FindTheRoutMapById(ctx, routid)
	res2, err := use.Repo.FindTrainByRoutid(ctx, route)
	response := domain.SearchingTrainResponseData{
		TrainNames:      res2.TrainNames,
		SearcheResponse: make([]domain.Train, len(res1.RouteMap)),
	}

	for i, data := range res1.RouteMap {
		response.SearcheResponse[i].Distance = data.Distance
	}

	fmt.Println("Train name :", response.TrainNames)

	return res2, err
}

// ViewTrain implements interfaces.TrainUseCase.
func (use *BookingUseCase) ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error) {
	res, err := use.Repo.ViewTrain(ctx)
	return res, err
}

func NewBookingUseCase(repo interfaces.BookingRepo) usecase.BookingUseCase {
	return &BookingUseCase{
		Repo: repo,
	}
}
