package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Train struct {
	TrainId        primitive.ObjectID  `json:"trainid,omitempty" bson:"_id,omitempty"`
	TrainNumber    uint                `json:"trainNumber" bson:"trainNumber,omitempty"`
	TrainName      string              `json:"trainName" bson:"trainName,omitempty" validate:"required,min=2,max=50"`
	Route          primitive.ObjectID  `json:"route,omitempty" bson:"route,omitempty"`
	TrainType      string              `json:"traintype" bson:"traintype,omitempty"`
	StartingTime   string              `json:"startingtime,omitempty" bson:"startingtime,omitempty"`
	EndingtingTime string              `json:"endingtingtime,omitempty" bson:"endingtingtime,omitempty"`
	Distance       float32             `json:"distance" bson:"distance,omitempty"`
	Time           primitive.Timestamp `json:"time" bson:"time,omitempty"`
	Compartment    []Compartment       `json:"compartment,omitempty" bson:"compartment,omitempty"`
}
type Compartment struct {
	Seatid primitive.ObjectID `json:"seatid,omitempty" bson:"_id,omitempty"`
}
type Station struct {
	StationId   primitive.ObjectID `json:"stationid" bson:"_id,omitempty"`
	StationName string             `json:"stationname" bson:"stationname,omitempty"`
	City        string             `json:"city" bson:"city,omitempty"`
	StationType string             `json:"stationtype" bson:"stationtype,omitempty"`
}
type Route struct {
	RouteId   primitive.ObjectID `json:"routeid" bson:"_id,omitempty"`
	RouteName string             `json:"routename" bson:"routename,omitempty"`
	RouteMap  []RouteStation     `json:"routemap,omitempty" bson:"routemap,omitempty"`
}
type RouteStation struct {
	StationId primitive.ObjectID     `json:"stationid" bson:"stationid,omitempty"`
	Distance  float32                `json:"distance,omitempty" bson:"distance,omitempty"`
	Time      *timestamppb.Timestamp `json:"time,omitempty" bson:"time,omitempty"`
}
type SearchingTrainRequstedData struct {
	Date                 string             `json:"data" bson:"data,omitempty"`
	SourceStationid      primitive.ObjectID `json:"sourcestationid,omitempty" bson:"sourcestationid,omitempty"`
	DestinationStationid primitive.ObjectID `json:"destinationstationid,omitempty" bson:"destinationstationid,omitempty"`
}
type SearchingTrainResponseData struct {
	TrainId         []string
	TrainNames      []string           `json:"trainname" bson:"trainname,omitempty"`
	RouteID         primitive.ObjectID `json:"routeID,omitempty" bson:"routeID,omitempty"`
	RouteName       string
	TrainNumber     []uint   `json:"trainNumber" bson:"trainNumber,omitempty"`
	Traintype       []string `json:"traintype" bson:"traintype,omitempty"`
	StartingTime    []string `json:"startingtime,omitempty" bson:"startingtime,omitempty"`
	EndingtingTime  []string `json:"endingtingtime,omitempty" bson:"endingtingtime,omitempty"`
	Stationid       []primitive.ObjectID
	Distance        []float32 `json:"distance,omitempty" bson:"distance,omitempty"`
	SearcheResponse []Train   `json:"searcheresponse,omitempty" bson:"searcheresponse,omitempty"`
}

type SeatDetails struct {
	Seatid         primitive.ObjectID
	SeatNumber     int    `json:"seatnumber,omitempty" bson:"seatnumber,omitempty"`
	SeatType       string `json:"seattype,omitempty" bson:"seattype,omitempty"`
	IsReserved     bool   `json:"isreserved,omitempty" bson:"isreserved,omitempty"`
	HasPowerOutlet bool   `json:"haspoweroutlet,omitempty" bson:"haspoweroutlet,omitempty"`
}

type Seats struct {
	SeatId       primitive.ObjectID `json:"seatid,omitempty" bson:"seatid,omitempty"`
	Price        int                `json:"price,omitempty" bson:"price,omitempty"`
	Availability bool               `json:"availability,omitempty" bson:"availability,omitempty"`
	TypeOfSeat   string             `json:"typeofseate,omitempty" bson:"typeofseate,omitempty"`
	Compartment  string             `json:"compartment,omitempty" bson:"compartment,omitempty"`
	SeatDetails  []SeatDetails      `json:"seatDetails,omitempty" bson:"seatDetails,omitempty"`
}

//Booking
type CompartmentDetails struct {
	SeatIds      primitive.ObjectID
	Price        int
	Availability bool
	TypeOfSeat   string
	Compartment  string
	SeatDetails  []SeatDetail
}
type SeatDetail struct {
	SeatNumbers    int
	SeatType       string
	SeatPosition   string
	IsReserved     bool
	HasPowerOutlet bool
}
type BookingResponse struct {
	CompartmentDetails []CompartmentDetails
}

//Adding seat
type SeatData struct {
	Price         float32
	NumbserOfSeat int
	Compartment   string `json:"compartment,omitempty" bson:"compartment,omitempty"`
	TypeOfSeat    string
}

type TrainAndRouteData struct {
	Routeid   primitive.ObjectID
	TrainName []string
	Stationid []primitive.ObjectID
	Distace   []float32
}

type SeatDetail2 struct {
	SeatNumber     int    `json:"seatnumber,omitempty" bson:"seatnumber,omitempty"`
	SeatType       string `json:"seattype,omitempty" bson:"seattype,omitempty"`
	IsReserved     bool   `json:"isreserved,omitempty" bson:"isreserved,omitempty"`
	HasPowerOutlet bool   `json:"haspoweroutlet,omitempty" bson:"haspoweroutlet,omitempty"`
}

type Compartment2 struct {
	SeatId       primitive.ObjectID `json:"_id" bson:"_id"`
	Price        int                `json:"price,omitempty" bson:"price,omitempty"`
	Availability bool               `json:"availability,omitempty" bson:"availability,omitempty"`
	TypeOfSeat   string             `json:"typeofseate,omitempty" bson:"typeofseate,omitempty"`
	Compartment  string             `json:"compartment,omitempty" bson:"compartment,omitempty"`
	SeatDetails  []SeatDetail2      `json:"seatDetails,omitempty" bson:"seatDetails,omitempty"`
}
type RouteResult struct {
	ID        string `bson:"_id"`
	RouteName string `bson:"routename"`
}
type RouteResponse struct {
	RouteResult []RouteResult
}
type BookingData struct {
	TrainId              string      `json:"trainid,omitempty" bson:"trainid,omitempty"`
	CompartmentId        string      `json:"compartmentid,omitempty" bson:"compartmentid,omitempty"`
	Userid               int64       `json:"userid,omitempty" bson:"userid,omitempty"`
	Travelers            []Travelers `json:"travelers"`
	SourceStationid      primitive.ObjectID
	DestinationStationid primitive.ObjectID
}
type CheckoutDetails struct {
	TrainName          string
	TrainNumber        int64
	SourceStation      string
	DestinationStation string
	Username           string
	Amount             float64
	PnrNumber          int64
	Traveler           []Travelers
}

type Travelers struct {
	Travelername string `json:"travelername" bson:"travelername"`
}

type Payment struct {
	TicketId             primitive.ObjectID `json:"ticketid,omitempty" bson:"ticketid,omitempty"`
	Paymentid            primitive.ObjectID `json:"paymentid,omitempty" bson:"paymentid,omitempty"`
	Trainname            string             `json:"trainname,omitempty" bson:"price,omitempty"`
	Sourcestationid      primitive.ObjectID `bson:"sourcestationid"`
	DestinationStationid primitive.ObjectID `bson:"destinationstationid"`
	TrainNumber          int64              `json:"trainnumber"`
	UserName             string             `json:"username"`
	Userid               int64
	PNRnumber            int64
	Travelers            []Travelers `json:"travelers"`
}

type Ticket struct {
	TicketId             primitive.ObjectID `bson:"train_id,omitempty"`
	Trainname            string             `bson:"trainname"`
	Trainnumber          int64              `bson:"trainnumber"`
	Sourcestationid      primitive.ObjectID `bson:"sourcestationid"`
	DestinationStationid primitive.ObjectID `bson:"destinationstationid"`
	PNRnumber            int64              `bson:"pnrnumber"`
	Userid               int64              `bson:"userid"`
	Username             string             `bson:"username"`
	Classname            string             `bson:"classname"`
	CompartmentId        primitive.ObjectID `bson:"compartmentid"`
	TotalAmount          float64            `bson:"amount"`
	SeatNumbers          []int64            `bson:"seatnumbers"`
	IsValide             bool
	Travelers            []Travelers `json:"travelers" bson:"travelers"`
}

type UserWallet struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Userid        int64              `bson:"userid,omitempty"`
	Username      string             `bson:"username"`
	Email         string             `bson:"email"`
	Amount        float64            `bson:"amount"`
	WalletBalance float64            `bson:"walletBalance"`
}
