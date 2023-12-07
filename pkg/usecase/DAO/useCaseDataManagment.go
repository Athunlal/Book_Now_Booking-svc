package DAO

import "github.com/athunlal/bookNowBooking-svc/pkg/domain"

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
