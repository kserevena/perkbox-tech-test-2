package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kserevena/perkbox-tech-test-2/api"
	"log"
	"net/http"
	"os"
)

func main() {

	log.Printf("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	router.HandleFunc("/coupons", api.PostCouponHandler).Methods(http.MethodPost)
	router.HandleFunc("/coupons", api.GetCouponsHandler).Methods(http.MethodGet)
	router.HandleFunc("/coupons/{couponId}", api.GetCouponHandler).Methods(http.MethodGet)
	router.HandleFunc("/coupons/{couponId}", api.PutCouponHandler).Methods(http.MethodPut)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
