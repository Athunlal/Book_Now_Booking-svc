package usecase

import (
	"context"
	"fmt"

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

// ViewTicket implements interfaces.BookingUseCase.
func (use *BookingUseCase) ViewTicket(ctx context.Context, tickets domain.Ticket) (*domain.Ticket, error) {
	res, err := use.Repo.GetTicketById(ctx, tickets)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return &res, nil
}

// UpdateAmount implements interfaces.BookingUseCase.
func (use *BookingUseCase) UpdateAmount(ctx context.Context, wallet domain.UserWallet) error {
	res, err := use.Repo.FetchWalletDatabyUserid(ctx, wallet)
	if err != nil {
		return err
	}
	totalAmount := res.WalletBalance + wallet.WalletBalance
	updateWallet := domain.UserWallet{
		Userid:        wallet.Userid,
		WalletBalance: totalAmount,
	}
	err = use.Repo.UpdateAmount(ctx, updateWallet)
	return err
}

// CreateWallet implements interfaces.BookingUseCase.
func (use *BookingUseCase) CreateWallet(ctx context.Context, wallet domain.UserWallet) error {
	err := use.Repo.CreateWallet(ctx, wallet)
	if err != nil {
		return err
	}
	return nil
}

// Payment implements interfaces.BookingUseCase.
func (use *BookingUseCase) Payment(ctx context.Context, paymentData domain.Payment) (*domain.Payment, error) {
	wallet, err := use.Repo.FetchWalletDatabyUserid(ctx, domain.UserWallet{
		Userid: paymentData.Userid,
	})
	if err != nil {
		return nil, err
	}
	ticket, err := use.Repo.GetTicketByPNR(ctx, paymentData.PNRnumber)

	if err := utils.PaymentCalculation(*wallet, ticket); err != nil {
		return nil, err
	}

	for _, ch := range ticket.SeatNumbers {
		err := use.Repo.UpdateCompartment(ctx, ch, ticket.CompartmentId)
		if err != nil {
			return nil, err
		}
	}

	return &domain.Payment{
		TicketId: ticket.TicketId,
	}, nil
}

// SeatBooking implements interfaces.BookingUseCase.
func (use *BookingUseCase) SeatBooking(ctx context.Context, bookingData domain.BookingData) (domain.CheckoutDetails, error) {

	//fetching train data
	TrainId, err := primitive.ObjectIDFromHex(bookingData.TrainId)
	trainData, err := use.Repo.FindTrainById(ctx, TrainId)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	//fetching compartment data
	Compartmentid, err := primitive.ObjectIDFromHex(bookingData.CompartmentId)
	seatDetail, err := use.Repo.GetSeatDetails(ctx, Compartmentid)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	price := utils.PriceCalculation(seatDetail, len(bookingData.Travelers))

	//check seat availability
	seatNumber, err := utils.CheckSeatAvailable(len(bookingData.Travelers), seatDetail)
	//fetch user data
	userData, err := usermodule.GetUserData(use.Client, bookingData.Userid)
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	travelers := []domain.Travelers{}
	for _, ch := range bookingData.Travelers {
		travler := domain.Travelers{
			Travelername: ch.Travelername,
		}
		travelers = append(travelers, travler)
	}

	pnr := utils.GeneratePNR()

	ticket := domain.Ticket{
		Trainid:              trainData.TrainId,
		Trainname:            trainData.TrainName,
		Trainnumber:          int64(trainData.TrainNumber),
		Sourcestationid:      bookingData.SourceStationid,
		DestinationStationid: bookingData.DestinationStationid,
		PNRnumber:            pnr,
		Username:             userData.Username,
		Classname:            seatDetail.TypeOfSeat,
		SeatNumbers:          seatNumber,
		TotalAmount:          price,
		CompartmentId:        Compartmentid,
		Travelers:            travelers,
		IsValide:             true,
	}

	err = use.Repo.CreatTicket(ctx, ticket)

	if err != nil {
		return domain.CheckoutDetails{}, err
	}

	return domain.CheckoutDetails{
		TrainName:   trainData.TrainName,
		TrainNumber: int64(trainData.TrainNumber),
		Username:    userData.Username,
		Amount:      price,
		Traveler:    travelers,
		PnrNumber:   pnr,
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
