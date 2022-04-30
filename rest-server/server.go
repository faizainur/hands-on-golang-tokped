package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	pb "github.com/faizainur/hands-on-golang/rest-server/cryptos_pb"
	"github.com/faizainur/hands-on-golang/rest-server/models"
	"github.com/faizainur/hands-on-golang/rest-server/routes"
	"github.com/faizainur/hands-on-golang/rest-server/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	DB_POSTGRES_HOST = "localhost"
	DB_POSTGRES_PORT = "5432"
	DB_POSTGRES_USER = os.Getenv("DB_POSTGRES_USER")
	DB_POSTGRES_PASS = os.Getenv("DB_POSTGRES_PASS")
	DB_POSTGRES_NAME = os.Getenv("DB_POSTGRES_NAME")
)

func main() {
	fmt.Println("Connecting to gRPCs server...")
	conn, err := grpc.Dial(":6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connected to gRPCs server")

	cryptosGrpcClient := NewCryptosGrpcClient(conn)

	fmt.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/UTC", DB_POSTGRES_HOST, DB_POSTGRES_USER, DB_POSTGRES_PASS, DB_POSTGRES_NAME, DB_POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Password{})

	passwordService := services.NewPasswordService(db, cryptosGrpcClient)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// Ping Endpoint : Used for check connection to the server
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(struct {
			Code   int16  `json:"code,omitempty"`
			Status string `json:"status,omitempty"`
		}{
			Code:   200,
			Status: "OK",
		})
	})

	passwordRouter := r.PathPrefix("/password").Subrouter().StrictSlash(true)
	{
		passwordRoutes := routes.NewPasswordRoutes(passwordService)

		passwordRouter.HandleFunc("", passwordRoutes.AddPassword).Methods("POST").Queries("username", "{username}", "password", "{password}", "email", "{email}", "name", "{name}")
		passwordRouter.HandleFunc("", passwordRoutes.ListPasswords).Methods("GET").Queries("limit", "{limit}", "offset", "{offset}")
		passwordRouter.HandleFunc("/{id}", passwordRoutes.GetPassword).Methods("GET")
	}

	http.Handle("/", r)
	http.ListenAndServe(":9000", nil)
}

func NewCryptosGrpcClient(conn grpc.ClientConnInterface) pb.CryptosServiceClient {
	return pb.NewCryptosServiceClient(conn)
}
