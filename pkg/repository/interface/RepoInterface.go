package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
)

type BookingRepo interface {
	FindbyTrainName(ctx context.Context, train domain.Train) (domain.Train, error)
	FindByTrainNumber(tx context.Context, train domain.Train) (domain.Train, error)

	FindByStationName(ctx context.Context, station domain.Station) (domain.Station, error)
	FindroutebyName(ctx context.Context, route domain.Route) (domain.Route, error)

	SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)

	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
	FindSeatbyCompartment(ctx context.Context, seat domain.Seats) (error, domain.Seats)
}
