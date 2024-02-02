package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	// "os"
	"strconv"
	// "time"

	// "os"

	// "fmt"
	"log"
	// "math/rand"
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	// "golang.org/x/oauth2/clientcredentials"
	// "golang.org/x/tools/go/analysis/passes/appends"
)

type Item struct {
	ID         string
	Name       string
	Price      float64
	Quantity   int
	OrderItems []OrderItem `gorm:"foreignKey:ItemID"`
}

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
// 	ClientID:     os.Getenv("CLIENT_ID"),
// 	ClientSecret: os.Getenv("CLIENT_SECRET"),
// 	TokenURL:     os.Getenv("TOKEN_URL"),
// }

var clientCredsConfig = clientcredentials.Config{
	ClientID:     "NmjtwiEs8yPft4wii9WGwc_TPIca",
	ClientSecret: "m8HOuvjruaGnDXg6vMteXp9clcAa",
	TokenURL:     "https://sts.choreo.dev/oauth2/token",
}

func makeClient() *http.Client {
	var ctx = context.Background()

	var token, err = clientCredsConfig.TokenSource(context.Background()).Token()

	// fmt.Printf("getting tokenLL: %v\n", token)
	if err != nil {
		fmt.Printf("Error getting tokenLL: %v\n", err)
		return nil
	}

	fmt.Printf("Access Token: %s\n", token.AccessToken)

	client := oauth2.NewClient(ctx, clientCredsConfig.TokenSource(ctx))
	return client
	// resp, err := client.Get(serviceURL)
	// if err != nil {
	// 	fmt.Printf("Error making API request: %v\n", err)
	// 	return
	// }
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
	client := makeClient()
	resp, err := client.Get(serviceURL + "/item")
	if err != nil {
		fmt.Printf("Error making API request: %v\n", err)
		return
	}
	var respond Item
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("show respBody ", string(resBody))
	if err != nil {
		fmt.Println("error")
	}
	var items []Item

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Now, 'item' contains the structured data
	for _, item := range items {
		fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
	}
	// for _, item := range resBody {
	// err = json.Unmarshal(resBody, &respond)
	// fmt.Println()
	// if err != nil {
	// 	http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
	// 	return
	// }
	// }
	// err = json.Unmarshal([]byte(resBody), &respond)
	// err, _ = w.Write(resBody)
	// _ = json.NewDecoder().Decode(&respond)
	fmt.Println("After Unmarshal-respondbody")
	fmt.Println("RespondID start", respond.Name, "end ")

	// defer resp.Body.Close()
	// fmt.Printf("request: %v\n", resp)
	var order Order
	// var respond Item
	// _ = json.NewDecoder(resp.Body).Decode(&respond)

	_ = json.NewDecoder(r.Body).Decode(&order)
	order.ID = strconv.Itoa(rand.Intn(100000000))
	// fmt.Printf("respondID", respond.ItemID)
	// fmt.Printf("respondorderID", order.ID)
	// params:=mux.Vars(resp)
	// hasError := false
	// for _, item := range order.Items {
	// 	fmt.Println("HII2", item.ItemID)

	// 	fmt.Println("ZZZYYPO", item.ItemID)

	// 	ItemInst, _, err := resBody.GetItemById(item.ItemID)

	// 		fmt.Println("ZZZYYPO", order.Total)
	// 	if err != nil || ItemInst.Quantity < item.Quantity {
	// 		// Handle the error, for example, return an HTTP error response
	// 		hasError = true
	// 		break

	// 	}
	// 	ItemInst.Quantity = ItemInst.Quantity - item.Quantity
	// 	order.Total = ItemInst.Price*float64(item.Quantity) + order.Total
	// 	resp.UpdateItemQuantity(ItemInst, ItemInst.Quantity)

	// 	fmt.Println("HIIIPO", ItemInst.Quantity)
	// }
	// if hasError {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
	fmt.Println("respondorderID", order.ID)
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

// var config = &oauth2.Config{
// 	ClientID:     "Ux_BfWJHPssqfqyA9mkrcILfSAoa",
// 	ClientSecret: "k2f0zEfA4f20VjFO_MJYYgwfwNMa",

// 	// Scopes: []string{"SCOPE1", "SCOPE2"},
// 	Endpoint: oauth2.Endpoint{
// 		// AuthURL:  "",
// 		TokenURL: "https://sts.choreo.dev/oauth2/token",
// 	},
// }

var serviceURL = "https://4c49cc7f-a4f9-4bf4-937b-6e7dc9d97bae-dev.e1-eu-north-azure.choreoapis.dev/kdnv/itemservice/item-9e9/v1.0"

// var serviceURL = "http://localhost:9010/item"

// ClientSecret := os.Getenv("CLIENT_SECRET"),
// TokenURL :=     os.Getenv("TOKEN_URL")

// config := &oauth2.Config{
// 	ClientID:     "CLIENT_ID",
// 	ClientSecret: "CLIENT_SECRET",
// 	Endpoint: oauth2.Endpoint{
// 		// AuthURL:  "",
// 		TokenURL: "TOKEN_URL",
// 	},
// }

func main() {
	r := mux.NewRouter()

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

//			next.ServeHTTP(w, r)
//		})
//	}
//
//	func GetOAuth2Token() (*http.Client, error) {
//		// Get the OAuth2 token using the client credentials
//		client := conf.Client(context.Background())
//		return client
//	}
// func GetOAuth2Token() (*oauth2.TokenResponse, error) {
// 	// Get the OAuth2 token using the client credentials
// 	token, err := conf.Token(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }

// func GetOAuth2Token() (*oauth2.Token, error) {
// 	// Get the OAuth2 token using the client credentials
// 	token, err := clientCredsConfig.Token(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }

//	var clientCredsConfig = &clientcredentials.Config{
//		ClientID:     os.Getenv("CLIENT_ID"),
//		ClientSecret: os.Getenv("CLIENT_SECRET"),
//		TokenURL:     os.Getenv("TOKEN_URL"),
//	}
func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
