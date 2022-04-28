package main

import (
	"fmt"
	"log"
	"net/http"

	pb "github.com/faizainur/hands-on-golang/rest-server/cryptos_pb"

	"github.com/faizainur/hands-on-golang/rest-server/routes"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Connecting to gRPCs server...")
	conn, err := grpc.Dial(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connected to gRPCs server")

	cryptosGrpcClient := NewCryptosGrpcClient(conn)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	passwordRouter := r.PathPrefix("/password").Subrouter()

	passwordRoutes := routes.NewPasswordRoutes(cryptosGrpcClient)

	passwordRouter.HandleFunc("", passwordRoutes.AddPassword).Methods("POST").Queries("username_email", "{username_email}", "password", "{password}")
	passwordRouter.HandleFunc("/{id}", passwordRoutes.GetPassword).Methods("GET").Queries("password", "{password}")

	http.Handle("/", r)
	http.ListenAndServe(":4000", nil)
}

func NewCryptosGrpcClient(conn grpc.ClientConnInterface) pb.CryptosServiceClient {
	return pb.NewCryptosServiceClient(conn)
}
