package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func init() {
	mongodbConnectionString = "mongodb://root:example@localhost:27017/?authSource=admin&readPreference=primary&ssl=false"
}

func Test_clientImpl_StoreAndGetCoupon(t *testing.T) {

	defer clearCollection(t)

	var testCoupon1 = &Coupon{
		Id:    "testCoupon1",
		Name:  "TestCoupon",
		Brand: "TestBrand",
		Value: 20,
		//CreatedAt: time.Date(2020, 02, 01, 12, 0, 0, 0, time.UTC),
		Expiry: time.Date(2020, 02, 01, 12, 0, 0, 0, time.UTC),
	}

	dbClient := clientImpl{}

	err := dbClient.InsertCoupon(testCoupon1)
	assert.NoError(t, err)

	retrievedCoupon, err := dbClient.GetCoupon(testCoupon1.Id)

	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), retrievedCoupon.CreatedAt, time.Millisecond*500)
	retrievedCoupon.CreatedAt = testCoupon1.CreatedAt
	assert.Equal(t, testCoupon1, retrievedCoupon)
}

func Test_clientImpl_GetCoupons(t *testing.T) {

	defer clearCollection(t)

	var testCoupon1 = &Coupon{
		Id:     "testCoupon1",
		Name:   "TestCoupon",
		Brand:  "TestBrand",
		Value:  20,
		Expiry: time.Date(2020, 02, 01, 12, 0, 0, 0, time.UTC),
	}

	var testCoupon2 = &Coupon{
		Id:     "testCoupon2",
		Name:   "TestCoupon2",
		Brand:  "TestBrand2",
		Value:  50,
		Expiry: time.Date(2022, 11, 19, 18, 0, 0, 0, time.UTC),
	}

	dbClient := clientImpl{}

	err := dbClient.InsertCoupon(testCoupon1)
	assert.NoError(t, err)

	err = dbClient.InsertCoupon(testCoupon2)
	assert.NoError(t, err)

	coupons, err := dbClient.GetCoupons(nil)
	assert.NoError(t, err)
	for i, coupon := range coupons {
		assert.WithinDuration(t, time.Now(), coupon.CreatedAt, time.Millisecond*500)
		coupon.CreatedAt = time.Time{}
		coupons[i] = coupon
	}
	testCoupon1.CreatedAt = time.Time{}
	testCoupon2.CreatedAt = time.Time{}
	assert.ElementsMatch(t, []Coupon{*testCoupon1, *testCoupon2}, coupons)
}

func TestGetCouponsFiltering(t *testing.T) {

	// Populate test data
	defer clearCollection(t)

	testData := make([]*Coupon, 0)

	testCoupon1 := Coupon{
		Id:     "testCoupon1",
		Name:   "TestCoupon1",
		Brand:  "TestBrand1",
		Value:  20,
		Expiry: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testData = append(testData, &testCoupon1)

	testCoupon2 := Coupon{
		Id:     "testCoupon2",
		Name:   "TestCoupon2",
		Brand:  "TestBrand2",
		Value:  50,
		Expiry: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testData = append(testData, &testCoupon2)

	testCoupon3 := Coupon{
		Id:     "testCoupon3",
		Name:   "TestCoupon3",
		Brand:  "TestBrand3",
		Value:  100,
		Expiry: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testData = append(testData, &testCoupon3)

	testCoupon4 := Coupon{
		Id:     "testCoupon4",
		Name:   "TestCoupon4",
		Brand:  "TestBrand3",
		Value:  200,
		Expiry: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testData = append(testData, &testCoupon4)

	dbClient := clientImpl{}

	for _, coupon := range testData {
		err := dbClient.InsertCoupon(coupon)
		assert.NoError(t, err)
	}

	// Define tests
	tests := []struct {
		name          string
		filter        GetCouponsQueryFilter
		expectedNames []string
	}{
		{
			name:          "filter on brand",
			filter:        GetCouponsQueryFilter{Brand: testCoupon2.Brand},
			expectedNames: []string{testCoupon2.Name},
		},
		{
			name: "filter on name",
			filter: GetCouponsQueryFilter{
				Brand: "",
				Name:  testCoupon3.Name,
			},
			expectedNames: []string{testCoupon3.Name},
		},
		{
			name: "filter on name and brand",
			filter: GetCouponsQueryFilter{
				Brand: testCoupon3.Brand,
				Name:  testCoupon4.Name,
			},
			expectedNames: []string{testCoupon4.Name},
		},
	}

	// Execute tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			coupons, err := dbClient.GetCoupons(&test.filter)
			assert.NoError(t, err)
			foundCouponNames := make([]string, len(coupons))
			for i, coupon := range coupons {
				foundCouponNames[i] = coupon.Name
			}
			assert.ElementsMatch(t, test.expectedNames, foundCouponNames)
		})
	}
}

func clearCollection(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := getConnectedMongoClient(ctx)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	collection := client.Database(databaseName).Collection(couponsCollection)
	collection.DeleteMany(ctx, bson.D{})
}
