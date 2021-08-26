package btckeys

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrInvalidMOfNValue        = errors.New("invalid m and n value")
	ErrInvalidPublicKeyString  = errors.New("invalid public key string")
	ErrInvalidPublicKeysNumber = errors.New("wrong number of public keys")
)

//This will generate bech32 address from compressed pubkeys even if uncompressed pubkeys were provided
//to shorten the redeemscript length
func GenerateMultiSigAddress(publicKeyStrings []string, flagM int, flagN int) (multiSigAddress string, redeemScriptString string, err error) {
	publicKeys := make([][]byte, len(publicKeyStrings))

	for i, publicKeyString := range publicKeyStrings {
		publicKeyString = strings.TrimSpace(publicKeyString)
		p, err := hex.DecodeString(publicKeyString)
		if err != nil {
			return "", "", err
		}
		key, err := btcec.ParsePubKey(p, btcec.S256()) //ParsePubkey will do some validation work and result is uncompressed
		if err != nil {
			return "", "", err
		}
		compressedPublicKey := key.SerializeCompressed()
		publicKeys[i] = compressedPublicKey
	}

	redeemScript, err := NewMOfNRedeemScript(publicKeys, flagM, flagN)
	if err != nil {
		return "", "", err
	}
	redeemHash := btcutil.Hash160(redeemScript)
	addr, err := btcutil.NewAddressScriptHashFromHash(redeemHash, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}

	multiSigAddress = addr.EncodeAddress()
	redeemScriptString = hex.EncodeToString(redeemScript)

	return multiSigAddress, redeemScriptString, nil
}

func NewMOfNRedeemScript(publicKeys [][]byte, m int, n int) ([]byte, error) {
	if n < 1 || n > 15 {
		return nil, ErrInvalidMOfNValue
	}
	if m < 1 || m > n {
		return nil, ErrInvalidMOfNValue
	}

	if len(publicKeys) != n {
		return nil, ErrInvalidPublicKeysNumber
	}

	mOPCode := txscript.OP_1 + (m - 1)
	nOPCode := txscript.OP_1 + (n - 1)

	builder := txscript.NewScriptBuilder()
	builder.AddOp(byte(mOPCode))
	for _, publicKey := range publicKeys {
		builder.AddData(publicKey)
	}
	builder.AddOp(byte(nOPCode))
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	redeemScript, err := builder.Script()
	if err != nil {
		return nil, err
	}
	return redeemScript, nil
}

func IsCompressedPublicKeyString(publicKeyString string) bool {
	p, err := hex.DecodeString(publicKeyString)
	if err != nil {
		return false
	}
	return btcec.IsCompressedPubKey(p)
}
