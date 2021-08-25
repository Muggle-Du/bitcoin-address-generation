package btckeys

import (
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	hd "github.com/btcsuite/btcutil/hdkeychain"
)

type ExtendedKey struct {
	key *hd.ExtendedKey
}

const (
	HardenedKeyStart = hd.HardenedKeyStart
)

var HardenedSymbol = map[string]bool{
	"'": true,
	"h": true,
	"H": true,
}

var (
	ErrInvalidPath = errors.New("invalid derivation path")
)

func NewKeyFromString(xpubOrxprv string) (*ExtendedKey, error) {
	key, err := hd.NewKeyFromString(xpubOrxprv)
	if err != nil {
		return nil, err
	}

	exkey := &ExtendedKey{key: key}
	return exkey, nil
}

func (exkey *ExtendedKey) P2WPKHAddress() (address string, err error) {
	pubkey, _ := exkey.key.ECPubKey()
	witnessProg := btcutil.Hash160(pubkey.SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	address = addressWitnessPubKeyHash.EncodeAddress()
	return address, nil
}

func (exkey *ExtendedKey) Derive(path string) (*ExtendedKey, error) {
	childExkey := &ExtendedKey{}

	indexes, err := indexesFromPath(path)
	if err != nil {
		return nil, err
	}
	if len(indexes) == 0 {
		return exkey, nil
	}

	for i, index := range indexes {
		if i == 0 {
			childExkey.key, err = exkey.key.Derive(index)
			if err != nil {
				return nil, err
			}
		} else {
			childExkey.key, err = childExkey.key.Derive(index)
			if err != nil {
				return nil, err
			}
		}
	}

	return childExkey, nil
}

func indexesFromPath(path string) ([]uint32, error) {
	indexes := []uint32{}
	rawIndexes := strings.Split(path, "/")

	if strings.ToLower(strings.TrimSpace(rawIndexes[0])) == "m" {
		if len(rawIndexes) < 2 {
			return indexes, nil
		} else {
			rawIndexes = rawIndexes[1:]
		}
	}

	for _, rawIndex := range rawIndexes {
		rawIndex = strings.Trim(rawIndex, " ")
		isHardened := false

		if len(rawIndex) > 1 && HardenedSymbol[string(rawIndex[len(rawIndex)-1])] {
			isHardened = true
			rawIndex = string(rawIndex[:len(rawIndex)-1])
		}

		index, err := strconv.Atoi(rawIndex)
		finalIndex := uint32(0)
		if err != nil {
			return nil, ErrInvalidPath
		}

		if isHardened {
			finalIndex = uint32(index + HardenedKeyStart)
		} else {
			finalIndex = uint32(index)
		}

		indexes = append(indexes, finalIndex)
	}
	return indexes, nil
}
