package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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

func main() {
	// Load .env file for development environment
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		DB_POSTGRES_HOST = os.Getenv("DB_POSTGRES_HOST")
		DB_POSTGRES_PORT = os.Getenv("DB_POSTGRES_PORT")
		DB_POSTGRES_USER = os.Getenv("DB_POSTGRES_USER")
		DB_POSTGRES_PASS = os.Getenv("DB_POSTGRES_PASS")
		DB_POSTGRES_NAME = os.Getenv("DB_POSTGRES_NAME")
		// SERVER_IP_ADDRESS           = os.Getenv("SERVER_IP_ADDRESS")
		SERVER_PORT                 = os.Getenv("SERVER_PORT")
		CRYPTOS_GRPC_SERVER_ADDRESS = os.Getenv("CRYPTOS_GRPC_SERVER_ADDR")
		CRYPTOS_GRPC_SERVER_PORT    = os.Getenv("CRYPTOS_GRPC_SERVER_PORT")
	)

	log.Println("Connecting to gRPCs server...")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", CRYPTOS_GRPC_SERVER_ADDRESS, CRYPTOS_GRPC_SERVER_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to gRPCs server")

	cryptosGrpcClient := NewCryptosGrpcClient(conn)

	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/UTC", DB_POSTGRES_HOST, DB_POSTGRES_USER, DB_POSTGRES_PASS, DB_POSTGRES_NAME, DB_POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Password{})

	passwordService := services.NewPasswordService(db, cryptosGrpcClient)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// Ping Endpoint : check connection to the server
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
		passwordRouter.HandleFunc("/{id}", passwordRoutes.UpdatePassword).Methods("POST").Queries("password", "{password}")
		passwordRouter.HandleFunc("/{id}", passwordRoutes.DeletePassword).Methods("DELETE")
	}

	http.Handle("/", r)
	log.Printf("Listening at :%s\n", SERVER_PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil)
}

func NewCryptosGrpcClient(conn grpc.ClientConnInterface) pb.CryptosServiceClient {
	return pb.NewCryptosServiceClient(conn)
}
