package main

import (
	"context"
	"errors"
	"log"
	"net"

	btckeys "github.com/Muggle-Du/btckeys/btckeys"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var (
	ErrInvalidXpub           = errors.New("not a valid xpub string")
	ErrWantCompressedPubKeys = errors.New("want compressed public key strings")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	btckeys.UnimplementedBtcKeysServer
}

func (s *server) DeriveBech32AddressFromXpub(ctx context.Context, in *btckeys.DerivationRequest) (btckeys.Address, error) {
	xpub := in.Xpub
	path := in.Path

	if string(xpub[:4]) == "xprv" {
		log.Printf("refuse to serve because a xprv is uploaded and this can be dangerous")
		return btckeys.Address{Address: ""}, ErrInvalidXpub
	}

	exkey, err := btckeys.NewKeyFromString(xpub)
	if err != nil {
		log.Printf("invalid xpub string, error: %v\n", err)
		return btckeys.Address{Address: ""}, err
	}

	childExkey, err := exkey.Derive(path)
	if err != nil {
		log.Printf("failed to derive from parent extended key, error: %v\n", err)
		return btckeys.Address{Address: ""}, err
	}

	bech32Address, err := childExkey.P2WPKHAddress()
	if err != nil {
		log.Printf("failed to generate bech32address, error: %v\n", err)
	}

	return btckeys.Address{Address: bech32Address}, nil
}

/*
func (s *server) GenerateMultiSigAddress(ctx context.Context, in *btckeys.MultiSigRequest) (btckeys.Address, error) {
	publicKeyStrings := in.Pubkeys
	for _, publicKeyString := range publicKeyStrings {
		if btckeys.IsCompressedPublicKeyString(publicKeyString) {
			log.Printf("want compressed public key")
			return btckeys.Address{Address: ""}, ErrWantCompressedPubKeys
		}
	}

	multiSigAddress, _, err := btckeys.GenerateMultiSigAddress(publicKeyStrings, int(in.M), int(in.N))
	if err != nil {
		log.Printf("failed to generate multisig address, error: %v\n", err)
		return btckeys.Address{Address: ""}, err
	}
	return btckeys.Address{Address: multiSigAddress}, nil
}
*/
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	btckeys.RegisterBtcKeysServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
