package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
	ID    string
	Items []OrderItem
	Total float64
}

var serviceURL = os.Getenv("SERVICE_URL")
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var tokenURL = os.Getenv("TOKEN_URL")

var clientCredsConfig = clientcredentials.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	TokenURL:     tokenURL,
}

var client2 = clientCredsConfig.Client(context.Background())

// HardCode

func makeClient() *http.Client {
	var ctx = context.Background()

	var _, err = clientCredsConfig.TokenSource(context.Background()).Token()

	// fmt.Printf("getting tokenLL: %v\n", token)
	if err != nil {
		fmt.Printf("Error getting tokenLL: %v\n", err)
		return nil
	}

	// fmt.Printf("Access Token: %s\n", token.AccessToken)

	client := oauth2.NewClient(ctx, clientCredsConfig.TokenSource(ctx))
	return client

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
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}

				item.Quantity = item.Quantity + orderitem.Quantity

				updatedItem = *item
				err = UpdateItem(client, serviceURL, itemID, updatedItem)
				if err != nil {
					fmt.Println("Error updating item:", err)
					return
				}
			}

			orders = append(orders[:index], orders[index+1:]...)
			json.NewEncoder(w).Encode(orders)
			return
		}

	}
	http.NotFound(w, r)
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
	http.NotFound(w, r)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	client := makeClient()

	var hasError = false
	var order Order
	var updatedItem Item

	_ = json.NewDecoder(r.Body).Decode(&order)
	order.ID = strconv.Itoa(rand.Intn(100000000))
	order.Total = 0
	for _, orderitem := range order.Items {
		var itemID = orderitem.ItemID
		item, err := GetItemByID(client, serviceURL, itemID)

		if err != nil || orderitem.Quantity > item.Quantity {
			hasError = true
			break
		}
		item.Quantity = item.Quantity - orderitem.Quantity
		order.Total = item.Price*float64(orderitem.Quantity) + order.Total
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

}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	client := makeClient()

	var preUpdatedItem, updatedItem Item

	for index, orderInstance := range orders {
		if orderInstance.ID == params["orderId"] {
			for _, orderitem := range orderInstance.Items {
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}

				item.Quantity = item.Quantity + orderitem.Quantity

				preUpdatedItem = *item
				err = UpdateItem(client, serviceURL, itemID, preUpdatedItem)
				if err != nil {
					fmt.Println("Error updating item:", err)
					return
				}
			}

			var updatedOrder Order
			err := json.NewDecoder(r.Body).Decode(&updatedOrder)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedOrder.ID = params["orderId"]
			updatedOrder.Total = 0
			for _, updatedOrderItems := range updatedOrder.Items {
				item, err := GetItemByID(client, serviceURL, updatedOrderItems.ItemID)
				if err != nil {
					fmt.Println("Error getting  item:", err)
					return
				}
				updatedOrder.Total = item.Price*float64(updatedOrderItems.Quantity) + updatedOrder.Total
			}

			orders[index] = updatedOrder

			json.NewEncoder(w).Encode(updatedOrder)

			for _, orderitem := range updatedOrder.Items {
				var itemID = orderitem.ItemID
				item, err := GetItemByID(client, serviceURL, itemID)
				if err != nil {
					fmt.Println("Error getting item:", err)
					return
				}
				item.Quantity = item.Quantity - orderitem.Quantity

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
	http.NotFound(w, r)

}

func main() {
	r := mux.NewRouter()
	// fmt.Println(client2)
	fmt.Println("SERVICE_URL", serviceURL)
	os.Setenv("SERVICE_URL", "https://4c49cc7f-a4f9-4bf4-937b-6e7dc9d97bae-dev.e1-eu-north-azure.choreoapis.dev/kdnv/itemservice/item-9e9/v1.0")
	os.Setenv("CLIENT_ID", "NmjtwiEs8yPft4wii9WGwc_TPIca")
	os.Setenv("CLIENT_SECRET", "m8HOuvjruaGnDXg6vMteXp9clcAaitemscoon")
	os.Setenv("TOKEN_URL", "https://sts.choreo.dev/oauth2/token")
	orders = append(orders, Order{ID: "1", Items: []OrderItem{{ItemID: "1", Quantity: 2}}, Total: 600})
	orders = append(orders, Order{ID: "2", Items: []OrderItem{{ItemID: "1", Quantity: 2}, {ItemID: "2", Quantity: 3}}, Total: 720})

	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/order", GetOrder).Methods("GET")
	r.HandleFunc("/order/{orderId}", GetOrderById).Methods("GET")
	r.HandleFunc("/order/{orderId}", UpdateOrder).Methods("PUT")
	r.HandleFunc("/order/{orderId}", DeleteOrder).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9090", r))

}

func GetItems(client *http.Client, serviceURL string) ([]Item, error) {
	url := fmt.Sprintf("%s/item/", serviceURL)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get items. Status code: %d", resp.StatusCode)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var items []Item
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func GetItemByID(client *http.Client, serviceURL, itemID string) (*Item, error) {
	url := fmt.Sprintf("%s/item/%s", serviceURL, itemID)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get item. Status code: %d", resp.StatusCode)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item Item
	err = json.Unmarshal(resBody, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateItem(client *http.Client, serviceURL string, itemID string, updatedItem Item) error {
	requestBody, err := json.Marshal(updatedItem)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/item/%s", serviceURL, itemID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to update item. Status code: %d", resp.StatusCode)
	}

	return nil
}
