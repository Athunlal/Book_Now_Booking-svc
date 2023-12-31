package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SeateAllocation(seateData domain.SeatData) domain.Seats {
	seateNumberStartFrom := 1
	seateDetail := []domain.SeatDetails{}

	//Allocating an empty array
	for i := 0; i < seateData.NumbserOfSeat; i++ {
		seateDetail = append(seateDetail, domain.SeatDetails{
			SeatType: "null",
		})
	}

	for i := 0; i < seateData.NumbserOfSeat-1; i++ {
		seateDetail[i].SeatNumber = seateNumberStartFrom
		seateNumberStartFrom++
		seateDetail[i].IsReserved = true
		if i == 0 || (i+1)%4 == 0 {
			seateDetail[i].SeatType = "side"
			seateDetail[i+1].SeatType = "side"
			seateDetail[i].HasPowerOutlet = true
			seateDetail[i+1].HasPowerOutlet = true

		} else {
			seateDetail[i].SeatType = "mid"
		}
	}

	allocated := domain.Seats{
		Price:        int(seateData.Price),
		Availability: true,
		TypeOfSeat:   seateData.TypeOfSeat,
		Compartment:  seateData.Compartment,
		SeatDetails:  seateDetail,
	}

	return allocated
}

func ConvertToPrimitiveTimestamp(pbTimestamp *timestamppb.Timestamp) primitive.Timestamp {
	seconds := pbTimestamp.GetSeconds()
	nanos := pbTimestamp.GetNanos()
	return primitive.Timestamp{T: uint32(seconds), I: uint32(nanos)}
}

func PriceCalculation(seatDetails domain.Compartment2, numberofTraverls int) float64 {
	return float64(seatDetails.Price) * float64(numberofTraverls)
}

func PaymentCalculation(wallet domain.UserWallet, ticket *domain.Ticket) error {
	if wallet.WalletBalance >= ticket.TotalAmount {
		return nil
	}
	return fmt.Errorf("Insufficient funds")
}

func GeneratePNR() int64 {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000000)
	randomNumber = randomNumber % 1000000
	return int64(randomNumber)
}

func IsValidTicket(result domain.TicketResponse) error {
	if !result.IsValide {
		return errors.New("Ticket canceld")
	}
	return nil
}

func CheckError(errCh chan error) error {
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}
