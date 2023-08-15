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
func (use *BookingUseCase) Booking(ctx context.Context, trainid domain.Train) (domain.SeatsForBookingResponse, error) {
	trainData, err := use.Repo.FindTrianById(ctx, trainid)
	var compartmentDetails domain.SeatsForBookingResponse

	for i, comp := range trainData.Compartment {
		seatData, err := use.Repo.GetSeatDetails(ctx, comp.Seatid)
		if err != nil {
			return compartmentDetails, err
		}

		compartmentDetails.SeatId = append(compartmentDetails.SeatId, seatData.SeatId)
		compartmentDetails.Price = append(compartmentDetails.Price, seatData.Price)
		compartmentDetails.Availability = append(compartmentDetails.Availability, seatData.Availability)
		compartmentDetails.TypeOfSeat = append(compartmentDetails.TypeOfSeat, seatData.TypeOfSeat)

		seatDetails := domain.SeatDetailsForBookingRespose{
			SeatNumber:     []int{},
			SeatType:       []string{},
			IsReserved:     []bool{},
			HasPowerOutlet: []bool{},
		}

		for _, seat := range seatData.SeatDetails {
			seatDetails.SeatNumber = append(seatDetails.SeatNumber, seat.SeatNumber)
			seatDetails.SeatType = append(seatDetails.SeatType, seat.SeatType)
			seatDetails.IsReserved = append(seatDetails.IsReserved, seat.IsReserved)
			seatDetails.HasPowerOutlet = append(seatDetails.HasPowerOutlet, seat.HasPowerOutlet)
		}

		compartmentDetails.SeatDetails = append(compartmentDetails.SeatDetails, seatDetails)
		compartmentDetails.SeatId = append(compartmentDetails.SeatId, trainData.Compartment[i].Seatid)
	}

	return compartmentDetails, err
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
