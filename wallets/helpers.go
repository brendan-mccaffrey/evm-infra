package wallets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func LoadPkFromString(pkPath string) *ecdsa.PrivateKey {
	//
	// Load Flashbots Privatekey
	content, err := ioutil.ReadFile(pkPath)
	var fbPk *ecdsa.PrivateKey
	if err == nil {
		fbPk, err = crypto.HexToECDSA(string(content))
		if err != nil {
			log.Println("Unabl to load flashbots pkey")
			fbPk = nil
		}
	}

	return fbPk
}

func PkToAddress(pk *ecdsa.PrivateKey) string {
	return fmt.Sprintf("%x", elliptic.Marshal(secp256k1.S256(), pk.X, pk.Y))
}
