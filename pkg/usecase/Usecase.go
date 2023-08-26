package usecase

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	usecase "github.com/athunlal/bookNowBooking-svc/pkg/usecase/interface"
	"github.com/athunlal/bookNowBooking-svc/pkg/usermodule"
	"github.com/athunlal/bookNowBooking-svc/pkg/usermodule/pb"
	"github.com/athunlal/bookNowBooking-svc/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingUseCase struct {
	Repo   interfaces.BookingRepo
	Client pb.ProfileManagementClient
}

// CreateWallet implements interfaces.BookingUseCase.
func (use *BookingUseCase) CreateWallet(ctx context.Context, wallet domain.UserWallet) error {
	err := use.Repo.CreateWallet(ctx, wallet)
	if err != nil {
		return err
	}
	return nil
}

// AddAmount implements interfaces.BookingUseCase.
func (use *BookingUseCase) AddAmount(ctx context.Context, wallet domain.UserWallet) error {
	err := use.Repo.AddAmount(ctx, wallet)
	if err != nil {
		return err
	}
	return nil
}

// Payment implements interfaces.BookingUseCase.
func (use *BookingUseCase) Payment(ctx context.Context, paymentData domain.Payment) (domain.Payment, error) {
	return domain.Payment{}, nil
}

// SeatBooking implements interfaces.BookingUseCase.
func (use *BookingUseCase) SeatBooking(ctx context.Context, bookingData domain.BookingData) (domain.CheckoutDetails, error) {

	TrainId, err := primitive.ObjectIDFromHex(bookingData.TrainId)
	trainData, err := use.Repo.FindTrainById(ctx, TrainId)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	Compartmentid, err := primitive.ObjectIDFromHex(bookingData.CompartmentId)
	seatDetail, err := use.Repo.GetSeatDetails(ctx, Compartmentid)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	_, err = utils.CheckSeatAvailable(seatDetail)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	userData, err := usermodule.GetUserData(use.Client, bookingData.Userid)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	return domain.CheckoutDetails{
		TrainName:   trainData.TrainName,
		TrainNumber: int64(trainData.TrainNumber),
		Username:    userData.Username,
		Traveler:    []domain.Traveler{},
	}, nil
}

// Booking implements interfaces.BookingUseCase.
func (use *BookingUseCase) SearchCompartment(ctx context.Context, trainid domain.Train) (domain.BookingResponse, error) {
	trainData, err := use.Repo.FindTrainById(ctx, trainid.TrainId)
	if err != nil {
		return domain.BookingResponse{}, err
	}

	response := domain.BookingResponse{
		CompartmentDetails: make([]domain.CompartmentDetails, len(trainData.Compartment)),
	}

	for i, compartment := range trainData.Compartment {
		res, err := use.Repo.GetSeatDetails(ctx, compartment.Seatid)
		if err != nil {
			return domain.BookingResponse{}, err
		}

		response.CompartmentDetails[i].SeatIds = compartment.Seatid
		response.CompartmentDetails[i].Price = res.Price
		response.CompartmentDetails[i].Availability = res.Availability
		response.CompartmentDetails[i].TypeOfSeat = res.TypeOfSeat
		response.CompartmentDetails[i].Compartment = res.Compartment

		seatDetails := make([]domain.SeatDetail, len(res.SeatDetails))
		for j, seatDetail := range res.SeatDetails {
			seatDetails[j] = domain.SeatDetail{
				SeatNumbers:    seatDetail.SeatNumber,
				SeatType:       seatDetail.SeatType,
				IsReserved:     seatDetail.IsReserved,
				HasPowerOutlet: seatDetail.HasPowerOutlet,
			}
		}
		response.CompartmentDetails[i].SeatDetails = seatDetails
	}

	return response, nil
}

// SearchTrain implements interfaces.BookingUseCase.
func (use *BookingUseCase) SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	routeData, err := use.Repo.FindRouteById(ctx, searcheData)
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

func NewBookingUseCase(repo interfaces.BookingRepo, client pb.ProfileManagementClient) usecase.BookingUseCase {
	return &BookingUseCase{
		Repo:   repo,
		Client: client,
	}
}
