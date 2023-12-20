package DAO

import (
	"net/http"
	"strconv"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"github.com/athunlal/bookNowBooking-svc/pkg/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertingArraytoString(arr []int64) string {
	var seatNumber string
	for i, num := range arr {
		if i == len(arr)-1 {
			seatNumber += strconv.FormatInt(num, 10)
		} else {
			seatNumber += strconv.FormatInt(num, 10) + ","
		}
	}
	return seatNumber
}

func convertingTavelers(Travelers []domain.Travelers) []*pb.Travelers {
	var travelors []*pb.Travelers
	for _, traveler := range Travelers {
		pbTraveler := &pb.Travelers{
			Travelername: traveler.Travelername,
		}
		travelors = append(travelors, pbTraveler)
	}
	return travelors
}

func convertTicketToViewBookingResponse(ticket domain.Ticket) *pb.ViewTicketResponse {
	return &pb.ViewTicketResponse{
		Sourestation:       ticket.SourceStation,
		Destinationstation: ticket.DestinationStation,
		Trainname:          ticket.Trainname,
		Trainnumber:        ticket.Trainnumber,
		Travelers:          convertingTavelers(ticket.Travelers),
		PnRnumber:          ticket.PNRnumber,
		Username:           ticket.Username,
		Classname:          ticket.Classname,
		Totalamount:        float32(ticket.TotalAmount),
		Seatnumbers:        convertingArraytoString(ticket.SeatNumbers),
		Isvalide:           ticket.IsValide,
	}
}

func CreateBookingHistroyResponse(res *domain.BookingHistory) *pb.BookingHistoryResponse {
	var viewTicketResponses []*pb.ViewTicketResponse
	for _, ticket := range res.Ticket {
		viewTicketResponse := convertTicketToViewBookingResponse(ticket)
		viewTicketResponses = append(viewTicketResponses, viewTicketResponse)
	}

	return &pb.BookingHistoryResponse{
		Response: viewTicketResponses,
	}
}

func CreateTravelersForCheckout(req *pb.CheckoutRequest) []domain.Travelers {
	travelers := []domain.Travelers{}
	for _, ch := range req.Travelers {
		traveler := domain.Travelers{
			Travelername: ch.Travelername,
		}
		travelers = append(travelers, traveler)
	}
	return travelers
}

func CreateCheckoutResponse(res domain.CheckoutDetails, req *pb.CheckoutRequest) *pb.CheckoutResponse {
	return &pb.CheckoutResponse{
		TrainName:   res.TrainName,
		Trainnumber: res.TrainNumber,
		Username:    res.Username,
		Travelers:   req.Travelers,
		Amount:      float32(res.Amount),
		PNRnumber:   res.PnrNumber,
	}
}

func MapBookingData(req *pb.CheckoutRequest, sourceStationid primitive.ObjectID, destinationStationid primitive.ObjectID) domain.BookingData {
	return domain.BookingData{
		CompartmentId:        req.Compartmentid,
		TrainId:              req.TrainId,
		Userid:               req.Userid,
		SourceStationid:      sourceStationid,
		DestinationStationid: destinationStationid,
		Travelers:            CreateTravelersForCheckout(req),
	}
}

func PrepareSearchData(req *pb.SearchTrainRequest) (domain.SearchingTrainRequstedData, error) {
	sourceid, err := primitive.ObjectIDFromHex(req.Sourcestationid)
	if err != nil {
		return domain.SearchingTrainRequstedData{}, err
	}

	destinationid, err := primitive.ObjectIDFromHex(req.Destinationstationid)
	if err != nil {
		return domain.SearchingTrainRequstedData{}, err
	}

	searchData := domain.SearchingTrainRequstedData{
		Date:                 req.Date,
		SourceStationid:      sourceid,
		DestinationStationid: destinationid,
	}

	return searchData, nil
}

func HandleSearchError(err error) (*pb.SearchTrainResponse, error) {
	return &pb.SearchTrainResponse{
		Status: http.StatusUnprocessableEntity,
		Error:  err.Error(),
	}, err
}

func ConvertToTrainDataList(res domain.SearchingTrainResponseData) []*pb.TrainData {
	var trainDataList []*pb.TrainData
	if len(res.TrainNames) < 1 {
		return nil
	}
	for i := range res.TrainNames {
		trainData := &pb.TrainData{
			Trainid:      res.TrainId[i],
			Trainname:    res.TrainNames[i],
			TrainNumber:  int64(res.TrainNumber[i]),
			StartingTime: res.StartingTime[i],
			Endingtime:   res.EndingtingTime[i],
		}
		trainDataList = append(trainDataList, trainData)
	}
	return trainDataList
}
func CreateSearchTrainResponse(trainDataList []*pb.TrainData) *pb.SearchTrainResponse {
	return &pb.SearchTrainResponse{
		Status:    http.StatusOK,
		Traindata: trainDataList,
	}
}

//Search compartment handler response
func SearchCompartmentBookingResponse(res domain.BookingResponse) *pb.SearchCompartmentResponse {
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
		}
		bookingResponse.Compartment[i] = pbCompartment
	}
	return bookingResponse
}

//Ticket Response
func CreateViewTicketResponse(res *domain.TicketResponse) *pb.ViewTicketResponse {
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

	response := &pb.ViewTicketResponse{
		Trainname:          res.Trainname,
		Trainnumber:        res.Trainnumber,
		Sourestation:       res.Sourcestation,
		Destinationstation: res.DestinationStation,
		PnRnumber:          res.PNRnumber,
		Username:           res.Username,
		Classname:          res.Classname,
		Compartment:        res.Classname,
		SeatNumber:         0,
		Totalamount:        float32(res.TotalAmount),
		Seatnumbers:        seatNumber,
		Isvalide:           res.IsValide,
		Travelers:          Travelers,
	}
	return response
}
