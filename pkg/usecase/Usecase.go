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

// Booking implements interfaces.BookingUseCase.
func (use *BookingUseCase) Booking(ctx context.Context, trainid domain.Train) (domain.BookingResponse, error) {
	trainData, err := use.Repo.FindTrianById(ctx, trainid)
	response := domain.BookingResponse{
		CompartmentDetails: make([]domain.CompartmentDetails, len(trainData.Compartment)),
	}
	for i, ch := range trainData.Compartment {
		response.CompartmentDetails[i].SeatIds = ch.Seatid
		response.CompartmentDetails[i].Price = ch.
		compartmentDetails := domain.CompartmentDetails{
			SeatDetails: make([]domain.SeatDetail, len(ch.Seatid)),
		}
		
	}

	return response, err
}

// SearchTrain implements interfaces.BookingUseCase.
func (use *BookingUseCase) SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	routeData, err := use.Repo.FindRouteId(ctx, searcheData)
	trainData, err := use.Repo.FindTrainByRoutid(ctx, domain.Train{
		Route: routeData.RouteID,
	})

	return trainData, err
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
