package repository

import (
	"context"
	"errors"
	"time"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrainDataBase struct {
	DB *mongo.Database
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
func (db *TrainDataBase) FindTrianById(ctx context.Context, train domain.Train) (domain.Train, error) {
	collectionRoute := db.DB.Collection("train")
	var trainData domain.Train

	filter := bson.M{"_id": train.TrainId}

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
func (db *TrainDataBase) FindTrainByRoutid(ctx context.Context, train domain.Train) (domain.SearchingTrainResponseData, error) {
	var trainData domain.SearchingTrainResponseData

	filter := bson.M{"route": train.Route}
	cur, err := db.DB.Collection("train").Find(ctx, filter)
	if err != nil {
		return trainData, err
	}
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

	if err := cur.Err(); err != nil {
		return trainData, err
	}
	return trainData, nil
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
func (db *TrainDataBase) FindRouteId(ctx context.Context, searchData domain.SearchingTrainRequstedData) (domain.SearchingTrainResponseData, error) {
	collectionRoute := db.DB.Collection("route")
	sourceStationID := searchData.SourceStationid
	destinationStationID := searchData.DestinationStationid

	var routeDoc struct {
		ID       primitive.ObjectID `bson:"_id"`
		Routemap []struct {
			StationID primitive.ObjectID `bson:"stationid"`
			Distance  float32            `bson:"distance"`
		} `bson:"routemap"`
	}

	// Use aggregation to find the index of the source and destination stations
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"routemap.stationid": bson.M{
					"$in": []primitive.ObjectID{sourceStationID, destinationStationID},
				},
			},
		},
		{
			"$addFields": bson.M{
				"sourceIndex": bson.M{
					"$indexOfArray": []interface{}{
						"$routemap.stationid", sourceStationID,
					},
				},
			},
		},
		{
			"$match": bson.M{
				"routemap.stationid": destinationStationID,
				"sourceIndex":        bson.M{"$gt": 0},
			},
		},
	}

	opts := options.Aggregate().SetMaxTime(2 * time.Second) // Set a reasonable max execution time

	cursor, err := collectionRoute.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return domain.SearchingTrainResponseData{}, errors.New("No train found for this route")
	}

	if err := cursor.Decode(&routeDoc); err != nil {
		return domain.SearchingTrainResponseData{}, err
	}

	response := domain.SearchingTrainResponseData{
		RouteID:   routeDoc.ID,
		Stationid: make([]primitive.ObjectID, len(routeDoc.Routemap)),
		Distance:  make([]float32, len(routeDoc.Routemap)),
	}

	for i, ch := range routeDoc.Routemap {
		response.Stationid[i] = ch.StationID
		response.Distance[i] = ch.Distance
	}

	return response, nil
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
