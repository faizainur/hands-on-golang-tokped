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

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	if os.Getenv("ENV") == "development" {
		godotenv.Load()
	}
	var (
		// SERVER_IP_ADDRESS = os.Getenv("SERVER_IP_ADDRESS")
		SERVER_PORT     = os.Getenv("SERVER_PORT")
		SECRET_DIR      = os.Getenv("SECRET_DIR")
		SECRET_FILENAME = os.Getenv("SECRET_FILENAME")
	)

	// Check if the program able to read the env variables
	// but the variables is not configured properly
	if SECRET_DIR == "" || SECRET_FILENAME == "" {
		log.Fatalf("%s\n", "Error : Please set the env variables properly")
		return
	}

	// If path to the secret file is not provided (both env varibles is not configured)
	// then use the default path
	// default path = $(PWD)/.tmp-handson-tokped
	if SECRET_DIR == "" && SECRET_FILENAME == "" {
		// SECRET_PATH = fmt.Sprintf("%s/%s", GetKeyDirPath(), "master.key")
		log.Println("Secret file is not provided, using the default key")
		SECRET_DIR = GetKeyDirPath()
		SECRET_FILENAME = "master.key"
	}

	key := loadSecret(SECRET_FILENAME, SECRET_DIR)

	tcpServer, err := net.Listen("tcp", fmt.Sprintf(":%s", SERVER_PORT))
	if err != nil {
		log.Fatalf("Error : %s ", err.Error())
	}
	s := grpc.NewServer()
	cryptos.RegisterGrpcServer(s, &cryptos.Server{Key: key})
	reflection.Register(s)
	log.Printf("Server listening at %v\n", tcpServer.Addr())
	if err := s.Serve(tcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadSecret(filename, path string) []byte {
	log.Printf("Loading key from '%s'\n", fmt.Sprintf("%s/%s", path, filename))

	key := make([]byte, 24)
	encodedKey, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		log.Println("Cannot find existing master key...")
		log.Println("Generating new master key file...")
		key := generateKeyFile(path, filename)
		return key
	}
	log.Println("Success, key is loaded!")
	hex.Decode(key, encodedKey)
	return key
}

func generateKey() []byte {
	key := make([]byte, 24)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Error : ", err.Error())
	}
	return key
}

func generateKeyFile(path, filename string) []byte {
	generatedKey := generateKey()
	encodedKey := hex.EncodeToString(generatedKey)

	filePath := fmt.Sprintf("%s/%s", path, "master.key")
	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error : ", err.Error())
	}
	defer out.Close()

	_, err = out.WriteString(encodedKey)
	if err != nil {
		log.Fatal("Error : ", err.Error())
	}

	log.Printf("File is generated at %s", filePath)
	return generatedKey
}

func GetKeyDirPath() string {
	homeDIr, errHomeDir := os.UserHomeDir()
	if errHomeDir != nil {
		log.Fatal(errHomeDir.Error())
	}
	dirPath := fmt.Sprintf("%s/%s", homeDIr, ".tmp-handson-tokped")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			log.Fatal("Error : ", err)
		}
	}
	return dirPath
}
