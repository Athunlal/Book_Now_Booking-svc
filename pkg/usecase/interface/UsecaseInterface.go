package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
)

type BookingUseCase interface {
	SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)
	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
	SearchCompartment(ctx context.Context, trainid domain.Train) (domain.BookingResponse, error)
	SeatBooking(ctx context.Context, bookingData domain.BookingData) (domain.BookingResponse, error)
}
