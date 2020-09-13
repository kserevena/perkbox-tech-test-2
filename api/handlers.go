package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kserevena/perkbox-tech-test-2/database"
	"log"
	"net/http"
)

func PostCouponHandler(writer http.ResponseWriter, request *http.Request) {

	coupon := new(database.Coupon)
	err := json.NewDecoder(request.Body).Decode(&coupon)

	if err != nil {
		http.Error(writer, "badly formed body", http.StatusBadRequest)
		return
	}

	coupon.Id = uuid.New().String()

	err = database.NewDBClient().InsertCoupon(coupon)
	if err != nil {
		http.Error(writer, "Failed to store coupon", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(coupon)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func PutCouponHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
}

func GetCouponHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	couponId := vars["couponId"]

	coupon, err := database.NewDBClient().GetCoupon(couponId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(writer).Encode(coupon)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetCouponsHandler(writer http.ResponseWriter, request *http.Request) {

	coupons, err := database.NewDBClient().GetCoupons(nil)
	if err != nil {
		log.Printf("error retrieving coupons from database: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(coupons)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
