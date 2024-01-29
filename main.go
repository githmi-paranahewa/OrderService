package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// "golang.org/x/tools/go/analysis/passes/appends"
)

// type Item struct {
// 	ID         string
// 	Name       string
// 	Price      float64
// 	Quantity   int
// 	OrderItems []OrderItem `gorm:"foreignKey:ItemID"`
// }

type OrderItem struct {
	ItemID   string
	Quantity int
}

type Order struct {
	ID     string
	Items  []OrderItem
	Total  float64
	Status string
}

var orders []Order

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	json.NewEncoder(w).Encode(orders)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for index, instance := range orders {
		if instance.ID == params["orderId"] {
			orders = append(orders[:index], orders[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(orders)
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for _, instance := range orders {
		if instance.ID == params["orderId"] {
			json.NewEncoder(w).Encode(instance)
			return
		}
	}
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)
	order.ID = strconv.Itoa(rand.Intn(100000000))
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for index, instance := range orders {
		if instance.ID == params["orderId"] {
			orders = append(orders[:index], orders[index+1:]...)
			var order Order
			_ = json.NewDecoder(r.Body).Decode(&order)
			order.ID = params["orderId"]
			orders = append(orders, order)
			json.NewEncoder(w).Encode(order)
		}
	}

}

func main() {
	r := mux.NewRouter()

	orders = append(orders, Order{ID: "1", Items: []OrderItem{{ItemID: "1", Quantity: 2}}, Total: 2000, Status: "Ongoing"})
	orders = append(orders, Order{ID: "2", Items: []OrderItem{{ItemID: "1", Quantity: 2}, {ItemID: "2", Quantity: 3}}, Total: 2000, Status: "Ongoing"})
	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/order", GetOrder).Methods("GET")
	r.HandleFunc("/order/{orderId}", GetOrderById).Methods("GET")
	r.HandleFunc("/order/{orderId}", UpdateOrder).Methods("PUT")
	r.HandleFunc("/order/{orderId}", DeleteOrder).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:9090", r))
}
