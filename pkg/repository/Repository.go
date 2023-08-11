package repository

import (
	"context"
	"errors"

	"github.com/athunlal/bookNowBooking-svc/pkg/domain"
	interfaces "github.com/athunlal/bookNowBooking-svc/pkg/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainDataBase struct {
	DB *mongo.Database
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
	var trainNames domain.SearchingTrainResponseData

	filter := bson.M{"route": train.Route}
	cur, err := db.DB.Collection("train").Find(ctx, filter)
	if err != nil {
		return trainNames, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var train domain.Train
		if err := cur.Decode(&train); err != nil {
			return trainNames, err
		}
		trainNames.TrainNames = append(trainNames.TrainNames, train.TrainName)
	}

	if err := cur.Err(); err != nil {
		return trainNames, err
	}
	return trainNames, nil
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
			// Add other fields you need, e.g., "time", "distance"
		} `bson:"routemap"`
	}

	filter := bson.M{
		"routemap.stationid": bson.M{
			"$in": []primitive.ObjectID{sourceStationID, destinationStationID},
		},
	}

	err := collectionRoute.FindOne(ctx, filter).Decode(&routeDoc)
	if err != nil {
		return domain.SearchingTrainResponseData{}, err
	}

	routeMap := routeDoc.Routemap

	isTrue := false
	for j, ch := range routeMap {
		if ch.StationID == sourceStationID {
			for i := j + 1; i < len(routeMap); i++ {
				if routeMap[i].StationID == destinationStationID {
					isTrue = true
				}
			}
		}
	}

	if isTrue {
		return domain.SearchingTrainResponseData{
			RouteID: routeDoc.ID,
		}, nil
	}

	return domain.SearchingTrainResponseData{}, errors.New("No train find this route")
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
