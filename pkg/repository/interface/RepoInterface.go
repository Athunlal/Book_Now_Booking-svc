package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo interface {
	FindbyTrainName(ctx context.Context, train domain.Train) (domain.Train, error)
	FindByTrainNumber(tx context.Context, train domain.Train) (domain.Train, error)
	FindTrainByRoutid(ctx context.Context, train domain.Train) ([]domain.Train, error)
	FindTrainById(ctx context.Context, train_id primitive.ObjectID) (domain.Train, error)

	FindByStationName(ctx context.Context, station domain.Station) (domain.Station, error)
	FindroutebyName(ctx context.Context, route domain.Route) (domain.Route, error)

	FindRouteByStationId(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)
	FindTheRoutMapById(ctx context.Context, routeData domain.Route) (domain.Route, error)

	GetSeatDetails(ctxctx context.Context, seatId primitive.ObjectID) (domain.Compartment2, error)

	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
	FindSeatbyCompartment(ctx context.Context, seat domain.Seats) (error, domain.Seats)

	CreateWallet(ctx context.Context, wallet domain.UserWallet) error
	UpdateAmount(ctx context.Context, wallet domain.UserWallet) error
	FetchWalletDatabyUserid(ctx context.Context, wallet domain.UserWallet) (*domain.UserWallet, error)

	CreatTicket(ctx context.Context, ticketData domain.Ticket) error

	UpdateCompartment(ctx context.Context, seatNumber int64, compartmentID primitive.ObjectID, status bool) error

	GetTicketByPNR(ctx context.Context, PNR int64) (domain.Ticket, error)

	GetTicketById(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error)

	UpdateTicket(ctx context.Context, ticket domain.Ticket) error
	DeleteTicket(ctx context.Context, ticket domain.Ticket) error

	UpdateAvailableStatus(ctx context.Context, compartmentID primitive.ObjectID, status bool) error

	UpdateTicketValidateStatus(ctx context.Context, ticket domain.Ticket) error

	FindStationById(ctx context.Context, stationId primitive.ObjectID) (domain.Station, error)

	FindTicketByUserid(ctx context.Context, userId int64) (*mongo.Cursor, error)

	FindTrainByDate(ctx context.Context, date string) ([]domain.Train, error)
}
