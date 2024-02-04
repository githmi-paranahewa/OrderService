package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	// "os"
	"strconv"
	// "time"

	"os"

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
	ID       string
	Name     string
	Price    float64
	Quantity int
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

// var clientID = os.Getenv("CLIENT_ID")
// var clientSecret = os.Getenv("CLIENT_SECRET")

// var _ = os.Setenv("tokenURL", "TOKEN_URL") //os.Getenv("TOKEN_URL")
// var tokenURL = os.Getenv("TOKEN_URL")

// os.Setenv("ServiceURL", "SERVICE_URL")

var clientCredsConfig = clientcredentials.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	TokenURL:     tokenURL,
}

var serviceURL = os.Getenv("SERVICE_URL")
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var tokenURL = os.Getenv("TOKEN_URL")

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
	client := makeClient()
	params := mux.Vars(r)
	var updatedItem Item
	for index, orderInstance := range orders {
		if orderInstance.ID == params["orderId"] {

			for _, orderitem := range orderInstance.Items {
				fmt.Println("oreder", orderitem.ItemID)
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}

				fmt.Println("orderQ", orderitem.Quantity, "itemquan", item.Quantity)

				item.Quantity = item.Quantity + orderitem.Quantity

				// Print the retrieved item
				fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
				updatedItem = *item
				err = UpdateItem(client, serviceURL, itemID, updatedItem)
				if err != nil {
					fmt.Println("Error updating item:", err)
					return
				}

			}

			orders = append(orders[:index], orders[index+1:]...)
			break
		}

	}
	json.NewEncoder(w).Encode(orders)
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for _, orderInstance := range orders {
		if orderInstance.ID == params["orderId"] {
			json.NewEncoder(w).Encode(orderInstance)
			return
		}
	}
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	client := makeClient()

	// items, err := GetItems(client, serviceURL)
	// if err != nil {
	// 	fmt.Println("Error getting items:", err)
	// 	return
	// }

	// // Print the retrieved items
	// for _, item := range items {
	// 	fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
	// }

	var hasError = false
	var order Order
	var updatedItem Item
	// _ = json.NewDecoder(resp.Body).Decode(&respond)

	_ = json.NewDecoder(r.Body).Decode(&order)
	order.ID = strconv.Itoa(rand.Intn(100000000))

	for _, orderitem := range order.Items {
		fmt.Println("oreder", orderitem.ItemID)
		var itemID = orderitem.ItemID
		item, err := GetItemByID(client, serviceURL, itemID)

		fmt.Println("orderQ", orderitem.Quantity, "itemquan", item.Quantity)
		if err != nil || orderitem.Quantity > item.Quantity {
			hasError = true
			break
		}
		item.Quantity = item.Quantity - orderitem.Quantity
		order.Total = item.Price*float64(orderitem.Quantity) + order.Total

		// Print the retrieved item
		fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
		updatedItem = *item
		err = UpdateItem(client, serviceURL, itemID, updatedItem)
		if err != nil {
			fmt.Println("Error updating item:", err)
			return
		}

	}

	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
	fmt.Println("respondorderID", order.ID)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "pkgication/json")
	// params := mux.Vars(r)
	// for index, orderInstance := range orders {
	// 	if instance.ID == params["orderId"] {
	// 		orders = append(orders[:index], orders[index+1:]...)
	// 		var order Order
	// 		_ = json.NewDecoder(r.Body).Decode(&order)
	// 		order.ID = params["orderId"]
	// 		orders = append(orders, order)
	// 		json.NewEncoder(w).Encode(order)
	// 	}
	// }

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	client := makeClient()
	var preUpdatedItem, updatedItem Item
	fmt.Println("hi order update")
	for index, orderInstance := range orders {
		fmt.Println("hi ou index", index)
		fmt.Println("hi outside of ", orderInstance.ID, " para ", params["orderId"])
		if orderInstance.ID == params["orderId"] {

			for _, orderitem := range orderInstance.Items {
				fmt.Println("oreder", orderitem.ItemID)
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}

				fmt.Println("orderQ", orderitem.Quantity, "itemquan", item.Quantity)

				item.Quantity = item.Quantity + orderitem.Quantity

				// Print the retrieved item
				fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
				preUpdatedItem = *item
				err = UpdateItem(client, serviceURL, itemID, preUpdatedItem)
				if err != nil {
					fmt.Println("Error updating item:", err)
					return
				}

			}

			fmt.Println("hi inside of ", orderInstance.ID)
			var updatedOrder Order
			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA \n AAAAAAAAAAAA")
			err := json.NewDecoder(r.Body).Decode(&updatedOrder)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Update the existing order
			updatedOrder.ID = params["orderId"]
			orders[index] = updatedOrder

			// Respond with the updated order
			json.NewEncoder(w).Encode(updatedOrder)

			for _, orderitem := range updatedOrder.Items {
				fmt.Println("oreder", orderitem.ItemID)
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}

				fmt.Println("orderQUPDATE", orderitem.Quantity, "itemquan", item.Quantity)

				item.Quantity = item.Quantity - orderitem.Quantity

				// Print the retrieved item
				fmt.Printf("ID: %s, Name: %s, Price: %f, Quantity: %d\n", item.ID, item.Name, item.Price, item.Quantity)
				updatedItem = *item
				err = UpdateItem(client, serviceURL, itemID, updatedItem)
				if err != nil {
					fmt.Println("Error updating item:", err)
					return
				}

			}
			return
		}
	}

	// If the order with the specified ID is not found
	http.NotFound(w, r)

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

	// os.Setenv("ServiceURL", "SERVICE_URL")
	// os.Setenv("ServiceURL", "SERVICE_URL")
	// serviceURL := os.Getenv("ServiceURL")
	// h, err := os.LookupEnv(e)

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

func GetItems(client *http.Client, serviceURL string) ([]Item, error) {
	// Construct the URL for getting items
	url := fmt.Sprintf("%s/item/", serviceURL)

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get items. Status code: %d", resp.StatusCode)
	}

	// Read the response body
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON array into a slice of Item
	var items []Item
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func GetItemByID(client *http.Client, serviceURL, itemID string) (*Item, error) {
	// Construct the URL for getting a specific item by ID
	url := fmt.Sprintf("%s/item/%s", serviceURL, itemID)

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get item. Status code: %d", resp.StatusCode)
	}

	// Read the response body
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into an Item
	var item Item
	err = json.Unmarshal(resBody, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateItem(client *http.Client, serviceURL string, itemID string, updatedItem Item) error {
	// Convert the updatedItem struct to JSON

	requestBody, err := json.Marshal(updatedItem)
	if err != nil {
		return err
	}

	// Construct the URL with the itemID variable
	url := fmt.Sprintf("%s/item/%s", serviceURL, itemID)

	// Create a PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	fmt.Println("request updated", req, "afterma", string(requestBody))

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Perform the PUT request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to update item. Status code: %d", resp.StatusCode)
	}

	return nil
}
