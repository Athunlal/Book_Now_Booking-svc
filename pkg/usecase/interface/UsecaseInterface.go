package interfaces

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
)

type BookingUseCase interface {
	SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error)
	ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error)
}
