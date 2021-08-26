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

func (s *server) DeriveBech32AddressFromXpub(ctx context.Context, in *btckeys.DerivationRequest) (*btckeys.Address, error) {
	xpub := in.Xpub
	path := in.Path

	log.Printf("receive bech32 child wallet address request, xpub:%v, path:%v\n", xpub, path)

	if string(xpub[:4]) == "xprv" {
		log.Printf("refuse to serve because a xprv is uploaded and this can be dangerous\n")
		return &btckeys.Address{Address: ""}, ErrInvalidXpub
	}

	exkey, err := btckeys.NewKeyFromString(xpub)
	if err != nil {
		log.Printf("invalid xpub string, error: %v\n", err)
		return &btckeys.Address{Address: ""}, err
	}

	childExkey, err := exkey.Derive(path)
	if err != nil {
		log.Printf("failed to derive from parent extended key, error: %v\n", err)
		return &btckeys.Address{Address: ""}, err
	}

	bech32Address, err := childExkey.P2WPKHAddress()
	if err != nil {
		log.Printf("failed to generate bech32address, error: %v\n", err)
	}

	return &btckeys.Address{Address: bech32Address}, nil
}

func (s *server) GetMultiSigAddress(ctx context.Context, in *btckeys.MultiSigRequest) (*btckeys.MultiSigResponse, error) {
	publicKeyStrings := in.Pubkeys
	/*for _, publicKeyString := range publicKeyStrings {
		if btckeys.IsCompressedPublicKeyString(publicKeyString) {
			log.Printf("want compressed public key")
			return &btckeys.Address{Address: ""}, ErrWantCompressedPubKeys
		}
	}*/ // not necessary

	log.Printf("receive multisig address generation request,m: %v, n: %v\n", in.M, in.N)

	multiSigAddress, redeemScript, err := btckeys.GenerateMultiSigAddress(publicKeyStrings, int(in.M), int(in.N))
	if err != nil {
		log.Printf("failed to generate multisig address, error: %v\n", err)
		return &btckeys.MultiSigResponse{Address: "", Redeemscript: ""}, err
	}

	//return both multisig address and the redeem script
	return &btckeys.MultiSigResponse{Address: multiSigAddress, Redeemscript: redeemScript}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	btckeys.RegisterBtcKeysServer(s, &server{})
	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
