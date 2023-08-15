package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingRepo interface {
	FindbyTrainName(ctx context.Context, train domain.Train) (domain.Train, error)
	FindByTrainNumber(tx context.Context, train domain.Train) (domain.Train, error)
	FindTrainByRoutid(ctx context.Context, train domain.Train) (domain.SearchingTrainResponseData, error)
	FindTrianById(ctx context.Context, train domain.Train) (domain.Train, error)

	FindByStationName(ctx context.Context, station domain.Station) (domain.Station, error)
	FindroutebyName(ctx context.Context, route domain.Route) (domain.Route, error)

	FindRouteId(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)
	FindTheRoutMapById(ctx context.Context, routeData domain.Route) (domain.Route, error)

	GetSeatDetails(ctxctx context.Context, seatId primitive.ObjectID) (domain.Compartment2, error)

	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
	FindSeatbyCompartment(ctx context.Context, seat domain.Seats) (error, domain.Seats)
}
