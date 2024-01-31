package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// "os"

	// "fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/clientcredentials"
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

// var clientCredsConfig = clientcredentials.Config{
// 	ClientID:     "CLIENT_ID",
// 	ClientSecret: "CLIENT_SECRET",
// 	TokenURL:     "TOKEN_URL",
// }

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

	// client := clientCredsConfig.Client(context.Background())
	// // os.Setenv("ServiceURL", "SERVICE_URL")
	// serviceURL := os.Getenv("ServiceURL")
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",

		// Scopes: []string{"SCOPE1", "SCOPE2"},
		Endpoint: oauth2.Endpoint{
			// AuthURL:  "",
			TokenURL: "TOKEN_URL",
		},
	}
	// verifier := oauth2.GenerateVerifier()

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.Exchange(ctx, code)
	fmt.Printf("HHH")
	if err != nil {
		fmt.Printf("KKK")
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	_ = client

	// rootRouter := r.PathPrefix("/").Subrouter()
	// rootRouter.Use(authenticateMiddlewaretest)

	orders = append(orders, Order{ID: "1", Items: []OrderItem{{ItemID: "1", Quantity: 2}}, Total: 2000, Status: "Ongoing"})
	orders = append(orders, Order{ID: "2", Items: []OrderItem{{ItemID: "1", Quantity: 2}, {ItemID: "2", Quantity: 3}}, Total: 2000, Status: "Ongoing"})
	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/order", GetOrder).Methods("GET")
	r.HandleFunc("/order/{orderId}", GetOrderById).Methods("GET")
	r.HandleFunc("/order/{orderId}", UpdateOrder).Methods("PUT")
	r.HandleFunc("/order/{orderId}", DeleteOrder).Methods("DELETE")
	// http.Handle("/item", authenticateMiddleware(client, serviceURL)(r))
	log.Fatal(http.ListenAndServe(":9090", r))

}

// func authenticateMiddleware(client *http.Client, serviceURL string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			params := mux.Vars(r)
// 			a, err := client.Get(params[serviceURL])

// 			if err != nil {
// 				fmt.Println("url", a, "error", err)
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func authenticateMiddleware(next *http.Client, se string) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("middleware2")
// 		// next.ServeHTTP(w, r)
// 	})
// }

// func authenticateMiddlewaretest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("middleware2")
// 		_, err := GetOAuth2Token()
// 		if err != nil {
// 			fmt.Println("Error obtaining token:", err)
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		// client := clientCredsConfig.Client(context.Background())
// 		// os.Setenv("ServiceURL", "SERVICE_URL")
// 		// serviceURL := os.Getenv("ServiceURL")
// 		// h, err := os.LookupEnv(e)
// 		// a, err := client.Get(serviceURL)
// 		// if err != nil {
// 		// 	fmt.Println("url", a, "error", err)
// 		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		// 	return
// 		// }

// 		next.ServeHTTP(w, r)
// 	})
// }
// func GetOAuth2Token() (*http.Client, error) {
// 	// Get the OAuth2 token using the client credentials
// 	client := clientCredsConfig.Client(context.Background())
// 	return client, nil
// }
