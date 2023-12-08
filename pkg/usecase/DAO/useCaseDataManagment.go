package DAO

import (
	"context"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildResponse(req domain.Ticket, dataCh1, dataCh2 chan string) chan domain.TicketResponse {
	out := make(chan domain.TicketResponse)
	go func() {
		res := domain.TicketResponse{
			Sourcestation:      <-dataCh1,
			DestinationStation: <-dataCh2,
			Classname:          req.Classname,
			PNRnumber:          req.PNRnumber,
			SeatNumbers:        req.SeatNumbers,
			Username:           req.Username,
			TotalAmount:        req.TotalAmount,
			Trainname:          req.Trainname,
			Travelers:          req.Travelers,
			Trainnumber:        req.Trainnumber,
			IsValide:           req.IsValide,
		}
		out <- res
	}()
	return out
}

func MapBookingResponse(ctx context.Context, cur *mongo.Cursor) (domain.BookingHistory, error) {
	var bookingHistory domain.BookingHistory
	for cur.Next(ctx) {
		var ticket domain.Ticket
		if err := cur.Decode(&ticket); err != nil {
			return domain.BookingHistory{}, err
		}
		bookingHistory.Ticket = append(bookingHistory.Ticket, ticket)
	}

	return bookingHistory, nil
}

// func MapTrainSearchResponse(ctx context.Context, cur *mongo.Cursor) ([]domain.Train, error) {
// 	var TrainData []domain.Ticket
// 	for cur.Next(ctx) {
// 		var train domain.Train
// 		if err := cur.Decode(&train); err != nil {
// 			return []domain.Train{}, err
// 		}
// 		TrainData = append(TrainData, train.)
// 	}
// 	return TrainData, nil
// }
