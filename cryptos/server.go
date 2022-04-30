package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/faizainur/hands-on-golang/cryptos/cryptos"

	"google.golang.org/grpc"
)

const (
	SECRET_DIR_PATH = "/Users/faiz.rofiq/Documents"
)

func main() {
	fmt.Println("Hello World")
	key := loadSecret()
	fmt.Println("Key : ", string(key))

	tcpServer, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
	}
	s := grpc.NewServer()
	cryptos.RegisterGrpcServer(s, &cryptos.Server{Key: key})
	log.Printf("server listening at %v", tcpServer.Addr())
	if err := s.Serve(tcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadSecret() []byte {
	fmt.Println("Loading key...")

	key := make([]byte, 24)
	encodedKey, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", SECRET_DIR_PATH, "master.key"))
	if err != nil {
		fmt.Println("Generating new secret file...")
		key := generateKeyFile()
		return key
	}
	fmt.Println(string(encodedKey))
	hex.Decode(key, encodedKey)
	return key
}

func generateKey() []byte {
	key := make([]byte, 24)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
	}
	return key
}

func generateKeyFile() []byte {
	generatedKey := generateKey()
	encodedKey := hex.EncodeToString(generatedKey)
	fmt.Println(encodedKey)

	out, err := os.Create(fmt.Sprintf("%s/%s", SECRET_DIR_PATH, "master.key"))
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
	}
	defer out.Close()

	_, err = out.WriteString(encodedKey)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
	}
	return generatedKey

}
