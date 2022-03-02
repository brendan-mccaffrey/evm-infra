package wallets

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

var CachedKeys = make(map[int]*CachedKey)

type CachedKey struct {
	Wallet  *accounts.Wallet
	Account *accounts.Account
	Key     *keystore.Key
	Nonce   *big.Int
}

type ExecWallet struct {
	pw string
	Ks *keystore.KeyStore
}

func ImportPKey(c *cli.Context) error {
	//
	// Make temp ew
	ew, err := NewExecWallet(c.String("walletsPath"), c.String("pwPath"))
	if err != nil {
		log.Println("NewExecWallet", err)
		return err
	}
	return ew.ImportPrivateKeyFromString(c.String("pKey"))
}

func (ew *ExecWallet) ImportPrivateKeyFromString(pKey string) error {
	// convert pkey
	privatKeyECDSA, err := crypto.HexToECDSA(pKey)
	if err != nil {
		log.Println(err)
		return err
	}
	// get public
	publicKey := privatKeyECDSA.Public()
	publicKeyECDSA, success := publicKey.(*ecdsa.PublicKey)
	if !success {
		return errors.New("failed to cast: publicKey.(*ecdsa.PublicKey)")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	log.Println("Loaded", address)
	// load to wallet file
	_, err = ew.Ks.ImportECDSA(privatKeyECDSA, ew.pw)
	return err
}

func (ew *ExecWallet) GetAccountByIndex(i int, client *ethclient.Client) *CachedKey {
	// return from cache if we have
	if val, ok := CachedKeys[i]; ok {
		return val
	}
	wallet := ew.Ks.Wallets()[i]
	account := wallet.Accounts()[0]
	jsonBytes, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		log.Println("Reading file", err)
		return nil
	}
	pKey, err := keystore.DecryptKey(jsonBytes, ew.pw)
	if err != nil {
		log.Println("Decrypting key", err)
		return nil
	}
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		log.Println("Getting nonce", err)
		return nil
	}
	// success
	ret := &CachedKey{
		Wallet:  &wallet,
		Account: &account,
		Key:     pKey,
		Nonce:   big.NewInt(int64(nonce)),
	}
	CachedKeys[i] = ret
	return ret
}

func (ew *ExecWallet) GetAccountByAddr(client *ethclient.Client, addr common.Address) *CachedKey {
	var account accounts.Account

	// loop through accounts
	for _, wallet := range ew.Ks.Wallets()[:] {
		account = wallet.Accounts()[0]
		// if not the one we are looking for, go to next
		if account.Address != addr {
			continue
		}
		jsonBytes, err := ioutil.ReadFile(account.URL.Path)
		if err != nil {
			log.Println("Reading file", err)
			break
		}
		pKey, err := keystore.DecryptKey(jsonBytes, ew.pw)
		if err != nil {
			log.Println("Decrypting key", err)
			break
		}
		nonce, err := client.PendingNonceAt(context.Background(), account.Address)
		if err != nil {
			log.Println("Getting nonce", err)
			break
		}
		log.Println(nonce, account.Address)
		// success
		ret := &CachedKey{
			Wallet:  &wallet,
			Account: &account,
			Key:     pKey,
			Nonce:   big.NewInt(int64(nonce)),
		}
		return ret
	}
	return nil
}

func NewExecWallet(walletPath string, pwPath string) (*ExecWallet, error) {
	if walletPath == "" {
		walletPath = "../wallet-store"
	}
	//
	// Load key
	ks := keystore.NewKeyStore(walletPath, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		return nil, errors.New("keystore failed to load")
	}
	//
	// Read PW
	content, err := ioutil.ReadFile(pwPath)
	if err != nil {
		return nil, err
	}
	pw := string(content)
	return &ExecWallet{
		pw: pw,
		Ks: ks,
	}, nil
}

// I made this because I manually deleted a wallet, and
// I couldn't import it because the keystore had it cached
func (ew *ExecWallet) ClearCache() {
	CachedKeys = make(map[int]*CachedKey)
	log.Println(ew.Ks.Accounts())
	for _, account := range ew.Ks.Accounts() {
		ew.Ks.Delete(account, ew.pw)
	}
	log.Println(CachedKeys)
}
