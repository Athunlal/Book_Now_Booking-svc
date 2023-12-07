package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainDataBase struct {
	DB *mongo.Database
}

// FindByStationName implements interfaces.BookingRepo.
func (db *TrainDataBase) FindStationById(ctx context.Context, stationId primitive.ObjectID) (domain.Station, error) {
	filter := bson.M{"_id": stationId}
	var result domain.Station
	err := db.DB.Collection("station").FindOne(ctx, filter).Decode(&result)
	return result, err
}

// BookingHistory implements interfaces.BookingRepo.
func (db *TrainDataBase) BookingHistory(ctx context.Context, userid int64) (*domain.BookingHistory, error) {
	var bookingHistory domain.BookingHistory

	filter := bson.M{"userid": userid}
	cur, err := db.DB.Collection("tickets").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var ticket domain.Ticket
		if err := cur.Decode(&ticket); err != nil {
			return nil, err
		}
		bookingHistory.Ticket = append(bookingHistory.Ticket, ticket)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if len(bookingHistory.Ticket) == 0 {
		return nil, fmt.Errorf("no booking history found for user")
	}

	return &bookingHistory, nil
}

func (db *TrainDataBase) UpdateTicketValidateStatus(ctx context.Context, ticket domain.Ticket) error {
	collection := db.DB.Collection("tickets")
	filter := bson.M{"_id": ticket.TicketId}
	update := bson.M{"$set": bson.M{"isvalide": ticket.IsValide}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAvailableStatus implements interfaces.BookingRepo.
func (db *TrainDataBase) UpdateAvailableStatus(ctx context.Context, compartmentID primitive.ObjectID, status bool) error {
	collection := db.DB.Collection("seat")
	filter := bson.M{"_id": compartmentID}
	update := bson.M{"$set": bson.M{"availability": status}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTicket implements interfaces.BookingRepo.
func (db *TrainDataBase) DeleteTicket(ctx context.Context, ticket domain.Ticket) error {
	collection := db.DB.Collection("tickets")
	filter := bson.M{"_id": ticket.TicketId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// UpdateTicket implements interfaces.BookingRepo.
func (db *TrainDataBase) UpdateTicket(ctx context.Context, ticket domain.Ticket) error {
	collection := db.DB.Collection("tickets")
	filter := bson.M{"_id": ticket.TicketId}
	update := bson.M{"$set": bson.M{
		"paymentstatus": ticket.PaymentStatus,
		"isvalide":      ticket.IsValide,
	},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetTicketByPNRandUserId implements interfaces.BookingRepo.
func (db *TrainDataBase) GetTicketById(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error) {
	collection := db.DB.Collection("tickets")
	var tickets domain.Ticket
	filter := bson.M{"_id": ticket.TicketId}

	err := collection.FindOne(ctx, filter).Decode(&tickets)
	if err != nil {
		return domain.Ticket{}, err
	}
	return tickets, nil
}

// GetTicketByPNR implements interfaces.BookingRepo.
func (db *TrainDataBase) GetTicketByPNR(ctx context.Context, PNR int64) (domain.Ticket, error) {
	collectionRoute := db.DB.Collection("tickets")
	var ticket domain.Ticket

	filter := bson.M{"pnrnumber": PNR}

	err := collectionRoute.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		return domain.Ticket{}, err
	}
	return ticket, nil
}

// UpdateCompartment implements interfaces.BookingRepo.
func (db *TrainDataBase) UpdateCompartment(ctx context.Context, seatNumber int64, compartmentID primitive.ObjectID, status bool) error {
	collection := db.DB.Collection("seat")
	filter := bson.M{"_id": compartmentID, "seatDetails.seatnumber": seatNumber}
	update := bson.M{"$set": bson.M{"seatDetails.$.isreserved": status}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// CreatTicket implements interfaces.BookingRepo.
func (db *TrainDataBase) CreatTicket(ctx context.Context, ticketData domain.Ticket) error {
	collection := db.DB.Collection("tickets")
	_, err := collection.InsertOne(ctx, ticketData)
	if err != nil {
		return err
	}
	return nil
}

// FetchWalletDatabyUserid implements interfaces.BookingRepo.
func (db *TrainDataBase) FetchWalletDatabyUserid(ctx context.Context, wallet domain.UserWallet) (*domain.UserWallet, error) {
	collection := db.DB.Collection("wallet")
	filter := bson.M{"user_id": wallet.Userid}

	err := collection.FindOne(ctx, filter).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Update implements interfaces.BookingRepo.
func (db *TrainDataBase) UpdateAmount(ctx context.Context, wallet domain.UserWallet) error {
	collection := db.DB.Collection("wallet")
	filter := bson.M{"user_id": wallet.Userid}
	update := bson.M{"$set": bson.M{"walletBalance": wallet.WalletBalance}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// CreateWallet implements interfaces.BookingRepo.
func (db *TrainDataBase) CreateWallet(ctx context.Context, wallet domain.UserWallet) error {
	collection := db.DB.Collection("wallet")

	walletDocument := bson.M{
		"user_id":       wallet.Userid,
		"walletBalance": wallet.WalletBalance,
	}
	_, err := collection.InsertOne(ctx, walletDocument)
	if err != nil {
		return err
	}
	return nil
}

// GetSeatDetails retrieves seat details based on seat ID
func (db *TrainDataBase) GetSeatDetails(ctx context.Context, seatId primitive.ObjectID) (domain.Compartment2, error) {
	collectionSeat := db.DB.Collection("seat")
	var seatData domain.Compartment2
	filter := bson.M{"_id": seatId}
	err := collectionSeat.FindOne(ctx, filter).Decode(&seatData)
	return seatData, err
}

// FindTrianById implements interfaces.BookingRepo.
func (db *TrainDataBase) FindTrainById(ctx context.Context, train_id primitive.ObjectID) (domain.Train, error) {
	collectionRoute := db.DB.Collection("train")
	var trainData domain.Train

	filter := bson.M{"_id": train_id}

	err := collectionRoute.FindOne(ctx, filter).Decode(&trainData)
	if err != nil {
		return domain.Train{}, err
	}
	return trainData, nil
}

// FindTheRoutMapById implements interfaces.BookingRepo.
func (db *TrainDataBase) FindTheRoutMapById(ctx context.Context, routeData domain.Route) (domain.Route, error) {
	collectionRoute := db.DB.Collection("route")

	var route domain.Route

	filter := bson.M{"_id": routeData.RouteId}

	err := collectionRoute.FindOne(ctx, filter).Decode(&route)
	if err != nil {
		return domain.Route{}, err
	}
	return route, nil
}

// FindTrainByRoutid implements interfaces.BookingRepo.
func (db *TrainDataBase) queryTrainsByRoute(ctx context.Context, train domain.Train) (*mongo.Cursor, error) {
	filter := bson.M{"route": train.Route}
	cur, err := db.DB.Collection("train").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (db *TrainDataBase) mapTrainDataFromCursor(ctx context.Context, cur *mongo.Cursor) (domain.SearchingTrainResponseData, error) {
	var trainData domain.SearchingTrainResponseData

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var train domain.Train
		if err := cur.Decode(&train); err != nil {
			return trainData, err
		}
		trainData.TrainId = append(trainData.TrainId, train.TrainId.Hex())
		trainData.TrainNames = append(trainData.TrainNames, train.TrainName)
		trainData.TrainNumber = append(trainData.TrainNumber, train.TrainNumber)
		trainData.Traintype = append(trainData.Traintype, train.TrainType)
		trainData.StartingTime = append(trainData.StartingTime, train.StartingTime)
		trainData.EndingtingTime = append(trainData.EndingtingTime, train.EndingtingTime)
	}

	if len(trainData.TrainId) == 0 {
		return trainData, errors.New("No trains found for the specified route")
	}

	if err := cur.Err(); err != nil {
		return trainData, err
	}

	return trainData, nil
}

func (db *TrainDataBase) FindTrainByRoutid(ctx context.Context, train domain.Train) (domain.SearchingTrainResponseData, error) {
	cur, err := db.queryTrainsByRoute(ctx, train)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}

	return db.mapTrainDataFromCursor(ctx, cur)
}

// FindByStationName implements interfaces.BookingRepo.
func (db *TrainDataBase) FindByStationName(ctx context.Context, station domain.Station) (domain.Station, error) {
	filter := bson.M{"stationname": station.StationName}
	var result domain.Station
	err := db.DB.Collection("station").FindOne(ctx, filter).Decode(&result)
	return result, err
}

// FindByTrainNumber implements interfaces.BookingRepo.
func (db *TrainDataBase) FindByTrainNumber(tx context.Context, train domain.Train) (domain.Train, error) {
	filter := bson.M{"trainName": train.TrainName}
	var result domain.Train
	err := db.DB.Collection("train").FindOne(tx, filter).Decode(&result)

	return result, err
}

// FindSeatbyCompartment implements interfaces.BookingRepo.
func (db *TrainDataBase) FindSeatbyCompartment(ctx context.Context, seat domain.Seats) (error, domain.Seats) {
	filter := bson.M{"compartment": seat.Compartment}
	var result domain.Seats
	err := db.DB.Collection("seat").FindOne(ctx, filter).Decode(&result)
	return err, result
}

// FindbyTrainName implements interfaces.BookingRepo.
func (db *TrainDataBase) FindbyTrainName(ctx context.Context, train domain.Train) (domain.Train, error) {
	filter := bson.M{"trainName": train.TrainName}
	var result domain.Train
	err := db.DB.Collection("train").FindOne(ctx, filter).Decode(&result)
	return result, err
}

// FindroutebyName implements interfaces.BookingRepo.
func (db *TrainDataBase) FindroutebyName(ctx context.Context, route domain.Route) (domain.Route, error) {
	filter := bson.M{"routename": route.RouteName}
	var result domain.Route
	err := db.DB.Collection("route").FindOne(ctx, filter).Decode(&result)
	return result, err
}

// SearchTrain implements interfaces.BookingRepo.
func (db *TrainDataBase) FindRouteByStationId(ctx context.Context, searchData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {

	collection := db.DB.Collection("route")
	sourceStationID := searchData.SourceStationid
	destinationStationID := searchData.DestinationStationid

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"routemap.stationid": sourceStationID},
					{"routemap.stationid": destinationStationID},
				},
			},
		},
		{
			"$addFields": bson.M{
				"sstationIndex": bson.M{
					"$indexOfArray": []interface{}{"$routemap.stationid", sourceStationID},
				},
				"dstationIndex": bson.M{
					"$indexOfArray": []interface{}{"$routemap.stationid", destinationStationID},
				},
			},
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$ne": []interface{}{"$sstationIndex", -1}},
						{"$ne": []interface{}{"$dstationIndex", -1}},
						{"$lt": []interface{}{"$sstationIndex", "$dstationIndex"}},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       1,
				"routename": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}
	defer cursor.Close(context.Background())

	var results []domain.RouteResult

	if err := cursor.All(context.Background(), &results); err != nil {
		return domain.SearchingTrainResponseData{}, err
	}

	if len(results) == 0 {
		return domain.SearchingTrainResponseData{}, fmt.Errorf("no route found")
	}

	routeid, err := primitive.ObjectIDFromHex(results[0].ID)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}

	return domain.SearchingTrainResponseData{
		RouteID:   routeid,
		RouteName: results[0].RouteName,
	}, nil
}

// ViewTrain implements interfaces.BookingRepo.
func (db *TrainDataBase) ViewTrain(ctx context.Context) (*domain.SearchingTrainResponseData, error) {
	var Train []domain.Train
	cursor, err := db.DB.Collection("train").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var train domain.Train
		if err := cursor.Decode(&train); err != nil {
			return nil, err
		}
		Train = append(Train, train)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &domain.SearchingTrainResponseData{
		SearcheResponse: Train,
	}, nil
}

func NewTrainRepo(db *mongo.Database) interfaces.BookingRepo {
	return &TrainDataBase{
		DB: db,
	}
}
