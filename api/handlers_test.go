package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kserevena/perkbox-tech-test-2/database"
	"github.com/kserevena/perkbox-tech-test-2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostCouponHandler(t *testing.T) {

	// Mock out database layer
	mockDbClient := new(mocks.DbClient)
	database.SetMockClient(mockDbClient)
	defer database.UnsetMockClient()

	// Define input data
	couponToStore := database.Coupon{
		Id:        "",
		Name:      "",
		Brand:     "",
		Value:     0,
		CreatedAt: time.Time{},
		Expiry:    time.Time{},
	}

	// Define expected calls
	mockDbClient.
		On("InsertCoupon", mock.Anything).
		Run(func(args mock.Arguments) {
			coupon := args.Get(0).(*database.Coupon)
			assert.NotEqual(t, couponToStore.Id, coupon.Id, "Coupon ID should be set before persisting")
		}).
		Return(nil)

	// Execute code
	recorder := httptest.NewRecorder()
	postBody := new(bytes.Buffer)
	request := httptest.NewRequest(http.MethodPost, "/coupons", postBody)

	assert.NoError(t, json.NewEncoder(postBody).Encode(couponToStore))

	PostCouponHandler(recorder, request)

	// Verify results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockDbClient.AssertExpectations(t)

	response := new(database.Coupon)
	assert.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	assert.NotEqual(t, couponToStore.Id, response.Id, "coupon ID should have been generated and set")
	response.Id = couponToStore.Id
	assert.Equal(t, couponToStore, *response)
}

func TestGetCouponHandler(t *testing.T) {

	// Mock out database layer
	mockDbClient := new(mocks.DbClient)
	database.SetMockClient(mockDbClient)
	defer database.UnsetMockClient()

	// Construct data to returned by mock database client
	coupon := database.Coupon{
		Id:        "testCouponId",
		Name:      "",
		Brand:     "",
		Value:     0,
		CreatedAt: time.Time{},
		Expiry:    time.Time{},
	}

	// Define expected calls
	mockDbClient.
		On("GetCoupon", coupon.Id).
		Return(&coupon, nil)

	// Execute code
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/coupons/"+coupon.Id, nil)
	request = mux.SetURLVars(request, map[string]string{"couponId": coupon.Id})

	GetCouponHandler(recorder, request)

	// Verify result
	assert.Equal(t, http.StatusOK, recorder.Code)
	responseBody := database.Coupon{}
	assert.NoError(t, json.NewDecoder(recorder.Body).Decode(&responseBody))
	assert.Equal(t, coupon, responseBody)
}

func TestGetCouponsHandler(t *testing.T) {

	// Mock out database layer
	mockDbClient := new(mocks.DbClient)
	database.SetMockClient(mockDbClient)
	defer database.UnsetMockClient()

	// Construct data to returned by mock database client
	coupon1 := database.Coupon{
		Id:        "testCoupon1Id",
		Name:      "",
		Brand:     "",
		Value:     0,
		CreatedAt: time.Time{},
		Expiry:    time.Time{},
	}

	coupon2 := database.Coupon{
		Id:        "testCoupon2Id",
		Name:      "",
		Brand:     "",
		Value:     0,
		CreatedAt: time.Time{},
		Expiry:    time.Time{},
	}

	expectedCoupons := []database.Coupon{coupon1, coupon2}

	// Define expected calls
	mockDbClient.
		On("GetCoupons").
		Return(expectedCoupons, nil)

	// Execute code
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/coupons", nil)

	GetCouponsHandler(recorder, request)

	// Verify result
	assert.Equal(t, http.StatusOK, recorder.Code)
	responseBody := []database.Coupon{}
	assert.NoError(t, json.NewDecoder(recorder.Body).Decode(&responseBody))
	assert.Equal(t, expectedCoupons, responseBody)
}
