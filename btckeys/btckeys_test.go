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
	expectedAddress := "3CY2p4b8dKVdjoqqcscxTABYsNQViNybNp"
	expectedRedeemScript := "5221034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa2102ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd21032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed353ae"

	publicKeyStrings := make([]string, 3)
	publicKeyStrings[0] = "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
	publicKeyStrings[1] = "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd"
	publicKeyStrings[2] = "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3"
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
	} else {
		t.Logf("multisig address: %v\nreedeemScriptString: %v\n", multisigAddress, redeemScriptString)
		if multisigAddress == expectedAddress && redeemScriptString == expectedRedeemScript {
			t.Logf("allright as expected")
		} else {
			t.Logf("ERROR: not as expected")
		}
	}
}
