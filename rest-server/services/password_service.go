package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	pb "github.com/faizainur/hands-on-golang/rest-server/cryptos_pb"
	"github.com/faizainur/hands-on-golang/rest-server/models"
	"gorm.io/gorm"
)

type PasswordService struct {
	db           *gorm.DB
	cryptoClient pb.CryptosServiceClient
}

func NewPasswordService(db *gorm.DB, cl pb.CryptosServiceClient) *PasswordService {
	return &PasswordService{db, cl}
}

func (s *PasswordService) AddPassword(name, username, email, password string) (models.Password, error) {
	var pass models.Password
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := s.cryptoClient.EncryptData(ctx, &pb.CryptoRequest{
		Data: []byte(password),
		Type: 1,
	})
	if err != nil {
		return pass, err
	}

	pass.Name = name
	pass.Username = username
	pass.Email = email
	pass.Password = hex.EncodeToString(r.GetData())

	err = s.db.Create(&pass).Error
	if err != nil {
		return pass, err
	}

	return pass, nil
}

func (s *PasswordService) GetPassword(id uint) (models.Password, error) {
	var pass models.Password
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := s.db.Select("id", "name", "username", "email", "password").First(&pass, id).Error
	if err != nil {
		return pass, err
	}

	decodedPass, err := hex.DecodeString(pass.Password)
	if err != nil {
		return pass, err
	}
	r, err := s.cryptoClient.DecryptData(ctx, &pb.CryptoRequest{
		Data: decodedPass,
		Type: 1,
	})
	if err != nil {
		return pass, err
	}
	fmt.Println(r.GetData())
	pass.Password = string(r.GetData())
	return pass, nil
}

func (s *PasswordService) ListPasswords(limit, offset uint32) ([]models.Password, error) {
	type chanResponse struct {
		idx int
		old string
		new string
	}
	var pass []models.Password
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := s.db.Find(&pass).Error
	if err != nil {
		return pass, err
	}

	decryptedPasswordsChannel := make(chan chanResponse, len(pass))
	defer close(decryptedPasswordsChannel)

	for idx, val := range pass {
		go func(channel chan chanResponse, idx int, old string) {
			// fmt.Println("Decrypting...")
			decodedPass, err := hex.DecodeString(old)
			if err != nil {
				channel <- chanResponse{idx: idx, old: old, new: "error"}
			}
			r, err := s.cryptoClient.DecryptData(ctx, &pb.CryptoRequest{
				Data: decodedPass,
				Type: 1,
			})
			// fmt.Println("in: ", string(r.GetData()))
			if err != nil {
				channel <- chanResponse{idx: idx, old: old, new: "error"}
			}
			channel <- chanResponse{idx: idx, old: old, new: string(r.GetData())}
		}(decryptedPasswordsChannel, idx, val.Password)
	}

	for i := 0; i < len(pass); i++ {
		decryptedPasswords := <-decryptedPasswordsChannel
		pass[decryptedPasswords.idx].Password = decryptedPasswords.new
	}

	// fmt.Println("End of function")
	return pass, nil
}
