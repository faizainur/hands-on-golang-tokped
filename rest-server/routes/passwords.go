package routes

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "github.com/faizainur/hands-on-golang/rest-server/cryptos_pb"
	"github.com/gorilla/mux"
)

type Password struct {
	Username_email string `json:"username_email,omitempty"`
	Password       string `json:"password,omitempty"`
}

type PasswordRoutes struct {
	cryptoClient pb.CryptosServiceClient
}

func NewPasswordRoutes(c pb.CryptosServiceClient) *PasswordRoutes {
	return &PasswordRoutes{cryptoClient: c}
}

func (p *PasswordRoutes) AddPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Fprintf(w, "%s", "Add password")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("Password ", vars["password"])
	t, err := p.cryptoClient.EncryptData(ctx, &pb.CryptoRequest{
		Data: []byte(vars["password"]),
		Type: 1,
	})
	if err != nil {
		// log.Fatal(err)
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Password{
		Username_email: vars["username_email"],
		Password:       hex.EncodeToString(t.GetData()),
	})
}

func (p *PasswordRoutes) GetPassword(w http.ResponseWriter, r *http.Request) {
	// id := mux.Vars(r)["id"]
	password := mux.Vars(r)["password"]

	if len(password) < 1 {
		w.WriteHeader(400)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	decodedPassword, err := hex.DecodeString(password)
	fmt.Println(decodedPassword)
	if err != nil {
		// log.Fatal(err)
		w.WriteHeader(400)
		return
	}
	t, err := p.cryptoClient.DecryptData(ctx, &pb.CryptoRequest{
		Data: decodedPassword,
		Type: 1})

	if err != nil {
		// log.Fatal(err)
		w.WriteHeader(400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Password{
		Password: string(t.Data),
	})
}
