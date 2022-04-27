package cryptos

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	pb "github.com/faizainur/cryptos/cryptos_pb"
	"google.golang.org/grpc"
)

type CryptoRequest = pb.CryptoRequest
type CryptoResponse = pb.CryptoResponse

type Server struct {
	pb.UnimplementedCryptosServiceServer
	Key []byte
}

func NewClient(conn grpc.ClientConnInterface) pb.CryptosServiceClient {
	return pb.NewCryptosServiceClient(conn)
}

func RegisterGrpcServer(grpcServer *grpc.Server, s *Server) {
	pb.RegisterCryptosServiceServer(grpcServer, s)
}

func (s *Server) EncryptData(ctx context.Context, in *CryptoRequest) (*CryptoResponse, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		// panic(err)
		return nil, err
	}
	chipertext := gcm.Seal(nonce, nonce, in.GetData(), nil)

	return &CryptoResponse{
		Data: chipertext,
		Type: 1,
	}, nil
}

func (s *Server) DecryptData(ctx context.Context, in *CryptoRequest) (*CryptoResponse, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
		return nil, err
	}

	nonce := in.GetData()[:gcm.NonceSize()]
	chipertext := in.GetData()[gcm.NonceSize():]
	plaintext, err := gcm.Open(nonce, nonce, chipertext, nil)
	if err != nil {
		fmt.Errorf("Error : ", err.Error())
		return nil, err
	}
	return &pb.CryptoResponse{
		Data: plaintext,
		Type: in.GetType(),
	}, nil
}
