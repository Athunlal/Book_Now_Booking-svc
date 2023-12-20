package handler

import (
	"context"
	"net/http"

	"github.com/athunlal/bookNowBooking-svc/pkg/api/handler/DAO"
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

//Booking histroy
func (h *BookingHandler) BookingHistory(ctx context.Context, req *pb.BookingHistroyRequest) (*pb.BookingHistoryResponse, error) {
	res, err := h.useCasse.BookingHistory(ctx, req.Userid)
	if err != nil {
		return nil, err
	}
	return DAO.CreateBookingHistroyResponse(res), nil
}

//Ticket Cancelation
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
	res, err := h.useCasse.ViewTicket(ctx, domain.Ticket{TicketId: id})
	if err != nil {
		return nil, err
	}

	return DAO.CreateViewTicketResponse(res), nil
}

func (h *BookingHandler) UpdateAmount(ctx context.Context, req *pb.UpdateAmountRequest) (*pb.UpdateAmountResponse, error) {
	err := h.useCasse.UpdateAmount(ctx, domain.UserWallet{
		Userid:        req.Userid,
		WalletBalance: float64(req.WalletBalance),
	})
	if err != nil {
		return &pb.UpdateAmountResponse{}, err
	}
	return &pb.UpdateAmountResponse{
		Status: http.StatusOK,
	}, nil
}

func (h *BookingHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {

	err := h.useCasse.CreateWallet(ctx, domain.UserWallet{
		Userid:        req.Userid,
		WalletBalance: float64(req.WalletBalance),
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateWalletResponse{Status: http.StatusOK}
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
	sourceStationid, err := primitive.ObjectIDFromHex(req.Sourcestationid)
	if err != nil {
		return nil, err
	}
	destinationStationid, err := primitive.ObjectIDFromHex(req.Destinationstationid)
	if err != nil {
		return nil, err
	}

	res, err := h.useCasse.SeatBooking(ctx, DAO.MapBookingData(req, sourceStationid, destinationStationid))
	if err != nil {
		return nil, err
	}

	return DAO.CreateCheckoutResponse(res, req), nil
}

func (h *BookingHandler) SearchCompartment(ctx context.Context, req *pb.SearchCompartmentRequest) (*pb.SearchCompartmentResponse, error) {
	trainID, err := primitive.ObjectIDFromHex(req.Trainid)
	if err != nil {
		return nil, err
	}
	res, err := h.useCasse.SearchCompartment(ctx, domain.Train{TrainId: trainID})
	if err != nil {
		return nil, err
	}
	return DAO.SearchCompartmentBookingResponse(res), nil
}

func (h *BookingHandler) SearchTrain(ctx context.Context, req *pb.SearchTrainRequest) (*pb.SearchTrainResponse, error) {
	searchData, err := DAO.PrepareSearchData(req)
	if err != nil {
		return DAO.HandleSearchError(err)
	}
	res, err := h.useCasse.SearchTrain(ctx, searchData)
	if err != nil {
		return DAO.HandleSearchError(err)
	}
	return DAO.CreateSearchTrainResponse(DAO.ConvertToTrainDataList(res)), nil
}
