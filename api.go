package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v);
}

type apiFunc func(http.ResponseWriter, *http.Request) error

var (
	mongoClient *mongo.Client
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) connectToDatabase() {
	var uri string
	if uri = os.Getenv("MONGO_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)

	mongoClient = client

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("nj-online-sports-db").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

}

func (s *APIServer) disconnectFromDatabase() {
    if err := mongoClient.Disconnect(context.TODO()); err != nil {
        panic(err)
    }
}

func (s *APIServer) Run() {

	router := http.NewServeMux()

	router.HandleFunc("/account", s.handleAccount)
	router.HandleFunc("/account/{id}", s.handleGetAccount)
	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router);
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {

	// error handling in another PR
	// if err != nil {
	// 	return
	// }

	if r.Method == "GET" {
		s.handleGetAccount(w, r);
	}

	if r.Method == "POST" {
		s.handleCreateAccount(w, r);
	}

	if r.Method == "DELETE" {
		s.handleDeleteAccount(w, r);
	}

	fmt.Errorf("method not allowed %s", r.Method);
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {

	// The get Account responsibility should be retrieving and looking up a account,
	// given the project is in the early stages, this will eventuall be replaced by a different
	// functionality involving the login process

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your first name!")
	scanner.Scan();
	firstName := scanner.Text()

	fmt.Println("Enter your last name!")
	scanner.Scan()
	lastName := scanner.Text()

	account := NewAccount(firstName, lastName)
	fmt.Println("this is account", account);

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error during scan:", err)
		return
	}

	collection := mongoClient.Database("nj-online-sports-db").Collection("accounts")

	_, err := collection.InsertOne(context.Background(), account)

	if err != nil {
		log.Fatal("MongoDB insertion failed:", err)
		return
	}

	WriteJSON(w, http.StatusOK, account);
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	// Will be created
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Will be created
}