package repository

import (
	"context"
	"log"

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
func (db *TrainDataBase) SearchTrain(ctx context.Context, searcheData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	collectionTrain := db.DB.Collection("train")

	sourceStationID := searcheData.SourceStationid
	destinationStationID := searcheData.DestinationStationid

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "route",
			"localField":   "route",
			"foreignField": "_id",
			"as":           "train_route",
		}}},
		{{Key: "$unwind", Value: "$train_route"}},
		{{Key: "$match", Value: bson.M{
			"train_route.routemap.stationid": bson.M{"$all": []primitive.ObjectID{sourceStationID, destinationStationID}},
		}}},
	}

	cur, err := collectionTrain.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("Error executing aggregation pipeline: %v\n", err)
		return domain.SearchingTrainResponseData{}, err
	}
	defer cur.Close(context.Background())

	var trains []domain.Train
	for cur.Next(context.Background()) {
		var train domain.Train
		if err := cur.Decode(&train); err != nil {
			log.Printf("Error decoding document: %v\n", err)
			return domain.SearchingTrainResponseData{}, err
		}
		trains = append(trains, train)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Error reading cursor: %v\n", err)
		return domain.SearchingTrainResponseData{}, err
	}

	return domain.SearchingTrainResponseData{
		SearcheResponse: trains,
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
