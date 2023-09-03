package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"github.com/athunlal/bookNowBooking-svc/pkg/pb"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/usecase/interface"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	useCasse interfaces.BookingUseCase
	pb.BookingManagementServer
}

func NewBookingHandler(usecase interfaces.BookingUseCase) *BookingHandler {
	return &BookingHandler{
		useCasse: usecase,
	}
}

func (h *BookingHandler) CancelTicket(ctx context.Context, req *pb.CancelTicketRequest) (*pb.CancelTicketResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Ticketid)
	if err != nil {
		return nil, err
	}
	err = h.useCasse.CancelletionTicket(ctx, domain.Ticket{
		TicketId: id,
		Userid:   req.Userid,
	})

	if err != nil {
		return nil, err
	}

	return &pb.CancelTicketResponse{
		Status: http.StatusOK,
		Error:  "",
	}, nil

}

func (h *BookingHandler) ViewTicket(ctx context.Context, req *pb.ViewTicketRequest) (*pb.ViewTicketResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Ticketid)
	if err != nil {
		return nil, err
	}
	res, err := h.useCasse.ViewTicket(ctx, domain.Ticket{
		TicketId: id,
	})
	if err != nil {
		return nil, err
	}

	Travelers := []*pb.Travelers{}

	for _, ch := range res.Travelers {
		traveler := &pb.Travelers{
			Travelername: ch.Travelername,
		}

		Travelers = append(Travelers, traveler)
	}

	var seatNumber string
	for i, num := range res.SeatNumbers {
		if i == len(res.SeatNumbers)-1 {
			seatNumber += strconv.FormatInt(num, 10)
		} else {
			seatNumber += strconv.FormatInt(num, 10) + ","
		}
	}

	return &pb.ViewTicketResponse{
		Trainname:            res.Trainname,
		Travelers:            Travelers,
		Trainnumber:          res.Trainnumber,
		Sourgestationid:      res.Sourcestationid.Hex(),
		Destinationstationid: res.DestinationStationid.Hex(),
		PnRnumber:            res.PNRnumber,
		Userid:               res.Userid,
		Username:             res.Username,
		Classname:            res.Classname,
		Compartmentid:        res.CompartmentId.Hex(),
		Totalamount:          float32(res.TotalAmount),
		Seatnumbers:          seatNumber,
		Isvalide:             false,
	}, nil
}

func (h *BookingHandler) UpdateAmount(ctx context.Context, req *pb.UpdateAmountRequest) (*pb.UpdateAmountResponse, error) {
	wallet := domain.UserWallet{
		Userid:        req.Userid,
		WalletBalance: float64(req.WalletBalance),
	}

	err := h.useCasse.UpdateAmount(ctx, wallet)
	if err != nil {
		return &pb.UpdateAmountResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, err
	}
	return &pb.UpdateAmountResponse{
		Status: http.StatusOK,
	}, nil
}

func (h *BookingHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	wallet := domain.UserWallet{
		Userid:        req.Userid,
		WalletBalance: float64(req.WalletBalance),
	}
	err := h.useCasse.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	// Construct and return a response
	response := &pb.CreateWalletResponse{
		Status: http.StatusOK,
	}
	return response, nil
}

func (h *BookingHandler) Payment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {

	res, err := h.useCasse.Payment(ctx, domain.Payment{
		Userid:    req.Userid,
		PNRnumber: req.PNRnumber,
	})

	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{
		Ticketid: res.TicketId.Hex(),
	}, nil
}

func (h *BookingHandler) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {

	travelers := []domain.Travelers{}
	for _, ch := range req.Travelers {
		traveler := domain.Travelers{
			Travelername: ch.Travelername,
		}
		travelers = append(travelers, traveler)
	}

	sourceStationid, err := primitive.ObjectIDFromHex(req.Sourcestationid)
	if err != nil {
		return nil, err
	}
	destinationStationid, err := primitive.ObjectIDFromHex(req.Destinationstationid)
	if err != nil {
		return nil, err
	}

	res, err := h.useCasse.SeatBooking(ctx, domain.BookingData{
		CompartmentId:        req.Compartmentid,
		TrainId:              req.TrainId,
		Userid:               req.Userid,
		SourceStationid:      sourceStationid,
		DestinationStationid: destinationStationid,
		Travelers:            travelers,
	})

	return &pb.CheckoutResponse{
		TrainName:   res.TrainName,
		Trainnumber: res.TrainNumber,
		Username:    res.Username,
		Travelers:   req.Travelers,
		Amount:      float32(res.Amount),
		PNRnumber:   res.PnrNumber,
	}, err
}

func (h *BookingHandler) SearchCompartment(ctx context.Context, req *pb.SearchCompartmentRequest) (*pb.SearchCompartmentResponse, error) {
	trainID, err := primitive.ObjectIDFromHex(req.Trainid)
	if err != nil {
		return nil, err
	}

	bookingData := domain.Train{TrainId: trainID}
	res, err := h.useCasse.SearchCompartment(ctx, bookingData)
	if err != nil {
		return nil, err
	}

	// Construct BookingResponse
	bookingResponse := &pb.SearchCompartmentResponse{
		Compartment: make([]*pb.Compartment, len(res.CompartmentDetails)),
	}

	for i, compartment := range res.CompartmentDetails {
		var status string
		if compartment.Availability {
			status = "Available"
		}
		pbCompartment := &pb.Compartment{
			Compartmentid:     compartment.SeatIds.Hex(),
			Price:             strconv.Itoa(compartment.Price),
			Typeofseat:        compartment.TypeOfSeat,
			CompartmentName:   compartment.Compartment,
			Availablitystatus: status,
			// SeatDetails:     make([]*pb.SeatDetails, len(compartment.SeatDetails)),
		}

		// for j, seatDetail := range compartment.SeatDetails {
		// 	pbSeatDetail := &pb.SeatDetails{
		// 		Isreserved: strconv.FormatBool(seatDetail.IsReserved),
		// 		Seattype:   seatDetail.SeatType,
		// 		Seatnumber: int64(seatDetail.SeatNumbers),
		// 	}
		// 	pbCompartment.SeatDetails[j] = pbSeatDetail
		// }

		bookingResponse.Compartment[i] = pbCompartment
	}

	return bookingResponse, nil
}

func (h *BookingHandler) SearchTrain(ctx context.Context, req *pb.SearchTrainRequest) (*pb.SearchTrainResponse, error) {
	sourceid, err := primitive.ObjectIDFromHex(req.Sourcestationid)
	if err != nil {
		log.Fatal("Converting the string to primitive.ObjectId err", err)
	}
	destinationid, err := primitive.ObjectIDFromHex(req.Destinationstationid)
	if err != nil {
		log.Fatal("Converting the string to primitive.ObjectId err", err)
	}
	searchData := domain.SearchingTrainRequstedData{
		Date:                 req.Date,
		SourceStationid:      sourceid,
		DestinationStationid: destinationid,
	}

	res, err := h.useCasse.SearchTrain(ctx, searchData)
	if err != nil {
		return &pb.SearchTrainResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, err
	}

	// Convert the domain search result to protobuf TrainData
	var trainDataList []*pb.TrainData
	for i, _ := range res.TrainNames {
		trainData := &pb.TrainData{
			Trainid:      res.TrainId[i],
			Trainname:    res.TrainNames[i],
			TrainNumber:  int64(res.TrainNumber[i]),
			StartingTime: res.StartingTime[i],
			Endingtime:   res.EndingtingTime[i],
		}
		trainDataList = append(trainDataList, trainData)
	}

	response := &pb.SearchTrainResponse{
		Status:    http.StatusOK,
		Traindata: trainDataList,
	}

	return response, nil
}
