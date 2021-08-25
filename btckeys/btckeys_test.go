package btckeys

import (
	//"github.com/btcsuite/btcutil/bech32"
	//"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	hd "github.com/btcsuite/btcutil/hdkeychain"
)

func Test_bech32(t *testing.T) {
	expectedAddress := "bc1qdq9hx6ss94s22dfphce9l42f3swv5mkc5rc5jw"
	xpub := "xpub6G2be6v3iwTjaQEWdgxzc3wjohmoApRjxm22VLVxaqoPyFev1tKSdscGeyrYXqMDG74MKFbXXk2h56ds99VvrbdmimeCWWHZnAxYDteTBcC"
	key, err := hd.NewKeyFromString(xpub)
	if err != nil {
		t.Logf("%v", err)
	}
	pubkey, _ := key.ECPubKey()
	t.Logf("public key: %v\n", pubkey.SerializeCompressed())
	witnessProg := btcutil.Hash160(pubkey.SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	address := addressWitnessPubKeyHash.EncodeAddress()
	t.Logf("bech32 address: %v\n", address)
	t.Logf("generated as expected?: %v\n", address == expectedAddress)
}

func Test_derive(t *testing.T) {
	xpub := "xpub6G2be6v3iwTjaQEWdgxzc3wjohmoApRjxm22VLVxaqoPyFev1tKSdscGeyrYXqMDG74MKFbXXk2h56ds99VvrbdmimeCWWHZnAxYDteTBcC"
	path := "0'/0'"

	exkey, err := NewKeyFromString(xpub)
	if err != nil {
		t.Logf("failed to generate extended key from xpub: %v\n", err)
	}

	childExkey, err := exkey.Derive(path)
	if err != nil {
		t.Logf("failed to derive from parent: %v\n", err)
	} else {
		t.Logf("child xpub: %v\n", childExkey.key.String())
	}
}

func Test_multisig(t *testing.T) {
	publicKeyStrings := make([]string, 3)
	publicKeyStrings[0] = "04a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd"
	publicKeyStrings[1] = "046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187"
	publicKeyStrings[2] = "0411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e83"
	for i, p := range publicKeyStrings {
		if IsCompressedPublicKeyString(p) {
			t.Logf("string %v is compressed\n", i)
		} else {
			t.Logf("string %v is uncompressed\n", i)
		}
	}
	multisigAddress, redeemScriptString, err := GenerateMultiSigAddress(publicKeyStrings, 2, 3)
	if err != nil {
		t.Logf("multisig error: %v\n", err)
	}
	t.Logf("multisig address: %v\nreedeemScriptString: %v\n", multisigAddress, redeemScriptString)
}
