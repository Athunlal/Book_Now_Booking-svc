package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Train struct {
	TrainNumber uint               `json:"trainNumber" bson:"trainNumber,omitempty"`
	TrainName   string             `json:"trainName" bson:"trainName,omitempty" validate:"required,min=2,max=50"`
	Route       primitive.ObjectID `json:"route,omitempty" bson:"route,omitempty"`
	TrainType   string             `json:"traintype" bson:"traintype,omitempty"`
	Seatsid     primitive.ObjectID `json:"seatid,omitempty" bson:"_id,omitempty"`
	Compartment []Compartment      `json:"compartment,omitempty" bson:"compartment,omitempty"`
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
	SearcheResponse []Train `json:"searcheresponse,omitempty" bson:"searcheresponse,omitempty"`
}

type SeatDetails struct {
	SeatNumber     int    `json:"seatnumber,omitempty" bson:"seatnumber,omitempty"`
	SeatType       string `json:"seattype,omitempty" bson:"seattype,omitempty"`
	IsReserved     bool   `json:"isreserved,omitempty" bson:"isreserved,omitempty"`
	HasPowerOutlet bool   `json:"haspoweroutlet,omitempty" bson:"haspoweroutlet,omitempty"`
}

type Seats struct {
	Price        int           `json:"price,omitempty" bson:"price,omitempty"`
	Availability bool          `json:"availability,omitempty" bson:"availability,omitempty"`
	TypeOfSeat   string        `json:"typeofseate,omitempty" bson:"typeofseate,omitempty"`
	Compartment  string        `json:"compartment,omitempty" bson:"compartment,omitempty"`
	SeatDetails  []SeatDetails `json:"seatDetails,omitempty" bson:"seatDetails,omitempty"`
}

type SeatData struct {
	Price         float32
	NumbserOfSeat int
	Compartment   string `json:"compartment,omitempty" bson:"compartment,omitempty"`
	TypeOfSeat    string
}
