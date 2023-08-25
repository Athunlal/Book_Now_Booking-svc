package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
)

type BookingUseCase interface {
	SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)
	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
	SearchCompartment(ctx context.Context, trainid domain.Train) (domain.BookingResponse, error)
	SeatBooking(ctx context.Context, bookingData domain.BookingData) (domain.CheckoutDetails, error)
	Payment(ctx context.Context, paymentData domain.Payment) (domain.Payment, error)
	AddAmount(ctx context.Context, wallet domain.UserWallet) error

	CreateWallet(ctx context.Context, wallet int64) error
}
