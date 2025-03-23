package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"
	"os"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v);
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}


func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

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

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func (s *APIServer) Run() {
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	router := http.NewServeMux()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)


	http.ListenAndServe(s.listenAddr, router);
}

// functions return an error to actively promote error handling
// this ensures we actively chec, and look for errors rather than
// relying on error handlers

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error{
	if r.Method == "GET" {
		return s.handleGetAccount(w, r);
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r);
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r);
	}

	return fmt.Errorf("method not allowed %s", r.Method);
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("john", "lee");
	return WriteJSON(w, http.StatusOK, account);
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}