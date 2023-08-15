package handler

import (
	"context"
	"log"
	"net/http"

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

func (h *BookingHandler) Booking(ctx context.Context, req *pb.BookingRequest) (*pb.BookingResponse, error) {
	trainid, err := primitive.ObjectIDFromHex(req.Trainid)
	if err != nil {
		log.Fatal("Converting the string to primitive.ObjectId err", err)
	}
	bookingData := domain.Train{TrainId: trainid}

	res, err := h.useCasse.Booking(ctx, bookingData)
	if err != nil {
		// Handle error if needed
		return nil, err
	}

	bookingResponse := &pb.BookingResponse{
		Compartment: []*pb.Compartment{},
	}

	for _, comp := range res.Compartment {
		compartment := &pb.Compartment{
			CompartmentName: res.Compartment[], // Assuming 'CompartmentName' is the correct field name
			SeatDetails:     []*pb.SeatDetails{},
		}

		for _, seatDetail := range comp.SeatDetails {
			seatDetailMsg := &pb.SeatDetails{
				Price:          seatDetail.Price,
				Isreserved:     seatDetail.IsReserved,
				Seattype:       seatDetail.SeatType,
				Seatnumber:     int64(seatDetail.SeatNumber),
				Haspoweroutlet: seatDetail.HasPowerOutlet,
			}

			compartment.SeatDetails = append(compartment.SeatDetails, seatDetailMsg)
		}

		bookingResponse.Compartment = append(bookingResponse.Compartment, compartment)
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
