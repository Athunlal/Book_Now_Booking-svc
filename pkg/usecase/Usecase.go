package usecase

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	usecase "github.com/athunlal/bookNowBooking-svc/pkg/usecase/interface"
)

type BookingUseCase struct {
	Repo interfaces.BookingRepo
}

// SearchTrain implements interfaces.TrainUseCase.
func (use *BookingUseCase) SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	res, err := use.Repo.FindRouteId(ctx, searcheData)

	return res, err
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
