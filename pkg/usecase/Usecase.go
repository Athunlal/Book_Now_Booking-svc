package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	"github.com/athunlal/bookNowBooking-svc/pkg/usecase/DAO"
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

func (use *BookingUseCase) BookingHistory(ctx context.Context, userid int64) (*domain.BookingHistory, error) {

	cur, err := use.Repo.FindTicketByUserid(ctx, userid)
	if err != nil {
		return nil, err
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	bookingHistory, err := DAO.MapBookingResponse(ctx, cur)
	if err != nil {
		return nil, err
	}

	if len(bookingHistory.Ticket) == 0 {
		return nil, fmt.Errorf("no booking history found for user")
	}

	for i, ticket := range bookingHistory.Ticket {
		res, err := use.Repo.FindStationById(ctx, ticket.Sourcestationid)
		if err != nil {
			return nil, err
		}
		res2, err := use.Repo.FindStationById(ctx, ticket.DestinationStationid)
		if err != nil {
			return nil, err
		}
		bookingHistory.Ticket[i].DestinationStation = res2.StationName
		bookingHistory.Ticket[i].SourceStation = res.StationName
	}

	return &bookingHistory, nil

}

// CancelletionTicket implements interfaces.BookingUseCase.
func (use *BookingUseCase) CancelletionTicket(ctx context.Context, ticket domain.Ticket) error {
	res, err := use.Repo.GetTicketById(ctx, ticket)
	if err != nil {
		return err
	}
	if !res.IsValide {
		return errors.New("Already canceled")
	}
	err = use.Repo.UpdateAmount(ctx, domain.UserWallet{
		Userid:        ticket.Userid,
		WalletBalance: res.TotalAmount,
	})
	if err != nil {
		return err
	}

	for _, ch := range res.SeatNumbers {
		err := use.Repo.UpdateCompartment(ctx, ch, res.CompartmentId, true)
		if err != nil {
			return err
		}
	}

	err = use.Repo.UpdateTicketValidateStatus(ctx, domain.Ticket{
		TicketId: ticket.TicketId,
		IsValide: false,
	})
	if err != nil {
		return err
	}
	return nil
}

// ViewTicket implements interfaces.BookingUseCase.
func (use *BookingUseCase) ViewTicket(ctx context.Context, tickets domain.Ticket) (*domain.TicketResponse, error) {

	ticketDetails, err := use.Repo.GetTicketById(ctx, tickets)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)
	wg := &sync.WaitGroup{}
	dataCh1 := use.findStationById(ctx, wg, errCh, ticketDetails.Sourcestationid)
	dataCh2 := use.findStationById(ctx, wg, errCh, ticketDetails.DestinationStationid)

	res := DAO.BuildResponse(ticketDetails, dataCh1, dataCh2)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := utils.CheckError(errCh); err != nil {
		return nil, err
	}

	result := <-res
	if err := utils.IsValidTicket(result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (use *BookingUseCase) findStationById(ctx context.Context, wg *sync.WaitGroup, errCh chan error, stationId primitive.ObjectID) chan string {
	dataCh := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(dataCh)
		res, err := use.Repo.FindStationById(ctx, stationId)
		errCh <- err
		dataCh <- res.StationName
	}()
	return dataCh
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
	wallet, err := use.getWalletData(ctx, paymentData.Userid)
	if err != nil {
		return nil, err
	}

	ticket, err := use.getTicketData(ctx, paymentData.PNRnumber)
	if err != nil {
		return nil, err
	}

	if err := utils.PaymentCalculation(*wallet, &ticket); err != nil {
		return nil, err
	}

	if err := use.updatePaymentStatus(ctx, &ticket); err != nil {
		return nil, err
	}

	if err := use.updateSeatCompartments(ctx, &ticket); err != nil {
		return nil, err
	}

	if err := use.updateWalletBalance(ctx, wallet, &ticket); err != nil {
		return nil, err
	}

	return &domain.Payment{
		TicketId: ticket.TicketId,
	}, nil
}

func (use *BookingUseCase) deleteTicket(ctx context.Context, ticket domain.Ticket) error {
	return use.Repo.DeleteTicket(ctx, ticket)
}

func (use *BookingUseCase) getWalletData(ctx context.Context, userid int64) (*domain.UserWallet, error) {
	return use.Repo.FetchWalletDatabyUserid(ctx, domain.UserWallet{Userid: userid})
}

func (use *BookingUseCase) getTicketData(ctx context.Context, PNRnumber int64) (domain.Ticket, error) {
	return use.Repo.GetTicketByPNR(ctx, PNRnumber)
}

func (use *BookingUseCase) updatePaymentStatus(ctx context.Context, ticket *domain.Ticket) error {
	if ticket.PaymentStatus {
		return fmt.Errorf("Payment already done")
	}
	return use.Repo.UpdateTicket(ctx, domain.Ticket{
		TicketId:      ticket.TicketId,
		PaymentStatus: true,
		IsValide:      true,
	})
}

func (use *BookingUseCase) updateSeatCompartments(ctx context.Context, ticket *domain.Ticket) error {

	for _, ch := range ticket.SeatNumbers {
		err := use.Repo.UpdateCompartment(ctx, ch, ticket.CompartmentId, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (use *BookingUseCase) updateWalletBalance(ctx context.Context, wallet *domain.UserWallet, ticket *domain.Ticket) error {
	updateAmount := wallet.WalletBalance - ticket.TotalAmount
	return use.Repo.UpdateAmount(ctx, domain.UserWallet{
		Userid:        wallet.Userid,
		WalletBalance: updateAmount,
	})
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
	if err != nil {
		return domain.CheckoutDetails{}, err
	}

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
		Userid:               bookingData.Userid,
		PaymentStatus:        false,
		CompartmentId:        Compartmentid,
		Travelers:            travelers,
		IsValide:             false,
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

	if len(trainData.Compartment) < 1 {
		return domain.BookingResponse{}, errors.New("Compartment not found")
	}

	response, err := use.getSeatDetails(ctx, trainData)
	if err != nil {
		return domain.BookingResponse{}, nil
	}

	response = use.checkAvailablility(ctx, response)

	return response, nil
}

func (use *BookingUseCase) checkAvailablility(ctx context.Context, response domain.BookingResponse) domain.BookingResponse {
	for _, ch := range response.CompartmentDetails {
		if ok := utils.CheckAvailableStatus(ch.SeatDetails); !ok {
			use.Repo.UpdateAvailableStatus(ctx, ch.SeatIds, false)
		}
	}
	return response
}
func (use *BookingUseCase) getSeatDetails(ctx context.Context, trainData domain.Train) (domain.BookingResponse, error) {
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
	routeData, err := use.Repo.FindRouteByStationId(ctx, searcheData)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}
	trainData, err := use.Repo.FindTrainByRoutid(ctx, domain.Train{
		Route: routeData.RouteID,
	})
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}
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
