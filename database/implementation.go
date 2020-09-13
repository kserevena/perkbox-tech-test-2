package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var mongodbConnectionString string

const databaseName = "pbtt"
const couponsCollection = "coupons"

type clientImpl struct{}

var mockClient DbClient = nil

func init() {
	mongodbConnectionString = os.Getenv("MONGODB_CONNECTION_STRING")
}

type GetCouponsQueryFilter struct {
	Brand string
	Name  string
}

// Set the package to use the given mock client. Useful for unit testing
// remember to UnsetMockClient() when test is complete.
func SetMockClient(client DbClient) {
	mockClient = client
}

// Remove the set mock client. Useful for unit testing
func UnsetMockClient() {
	mockClient = nil
}

func NewDBClient() DbClient {
	if mockClient != nil {
		return mockClient
	}

	return clientImpl{}
}

func (clientImpl) InsertCoupon(coupon *Coupon) error {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := getConnectedMongoClient(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	coupon.CreatedAt = time.Now()

	collection := client.Database(databaseName).Collection(couponsCollection)
	_, err = collection.InsertOne(ctx, coupon)

	if err != nil {
		return fmt.Errorf("error creating document in mongodb: %w", err)
	}

	return nil
}

func (clientImpl) GetCoupon(id string) (coupon *Coupon, err error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := getConnectedMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	result := new(Coupon)
	err = client.Database(databaseName).Collection(couponsCollection).FindOne(ctx, bson.D{{"_id", id}}).Decode(result)
	if err != nil {
		err := fmt.Errorf("error decoding db response for coupone with ID %s: %w", id, err)
		return nil, err
	}

	return result, nil
}

func (clientImpl) GetCoupons(filter *GetCouponsQueryFilter) (coupons []Coupon, err error) {

	queryFilters := bson.D{}

	if filter != nil {
		if filter.Brand != "" {
			queryFilters = append(queryFilters, bson.E{"brand", filter.Brand})
		}
		if filter.Name != "" {
			queryFilters = append(queryFilters, bson.E{"name", filter.Name})
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := getConnectedMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	cursor, err := client.Database(databaseName).Collection(couponsCollection).Find(ctx, queryFilters)
	if err != nil {
		return nil, err
	}

	var results = []Coupon{}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func getConnectedMongoClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbConnectionString))
	if err != nil {
		return nil, fmt.Errorf("error creating mongodb client: %w", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}
	return client, nil
}
