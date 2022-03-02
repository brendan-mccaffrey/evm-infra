package wallets

// ----------------
// TODO: Fix errors
// ----------------

import (
	"context"

	"fmt"
	"log"
	"math/big"

	"pangea/bindings/token"
	"pangea/tools"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Address struct {
	Addr string `json:"address"`
}

type Addresses struct {
	Addrs []*Address `json:"bots"`
}

// func FundWalletsFromFile(c *cli.Context) error {
// 	// Read json
// 	jsonFile, err := os.Open("./app/assets/bots.json")
// 	if err != nil {
// 		return err
// 	}
// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	var addresses Addresses
// 	json.Unmarshal(byteValue, &addresses)

// 	// Make temp ew
// 	ew, err := NewExecWallet("", c.String("pwPath"))
// 	if err != nil {
// 		log.Println("NewExecWallet")
// 		return err
// 	}

// 	var toAddrs []common.Address
// 	// load wallets
// 	for _, addr := range addresses.Addrs {
// 		ew.ImportPrivateKeyFromString(addr.Addr)
// 		toAddrs = append(toAddrs, common.HexToAddress(addr.Addr))
// 	}

// 	// fromAddr := common.HexToAddress(c.String("fromAddr"))
// 	// client, err := ethclient.Dial("wss://eth-rinkeby.alchemyapi.io/v2/4oY8Mnw8Ya9ZIWsB-zX3-_LFaRlBPuqK")
// 	// if err != nil {
// 	// 	log.Println("making client")
// 	// 	return err
// 	// }
// 	// chainId, err := client.ChainID(context.Background())
// 	// if err != nil {
// 	// 	log.Println("getting chainID")
// 	// 	return err
// 	// }

// 	//
// 	// Start GasPrices service
// 	gasProvider := tools.NewGasProvider()
// 	go gasProvider.Start()
// 	for {
// 		v := gasProvider.Gas.Result.Instant
// 		if v == nil {
// 			log.Println("Waiting for first gas update")
// 		} else {
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// 	totalAmt := new(big.Int).Mul(big.NewInt(int64(c.Uint64("amount"))), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(c.Uint64("decimals"))), nil))
// 	FundWallets(client, gasProvider, fromAddr, c.String("pwPath"), totalAmt, toAddrs, c.String("token"))
// 	return nil
// }

// func FundWalletsOld(client *ethclient.Client, gasProvider tools.GasProvider, fromAddr common.Address, pwPath string, totalAmt *big.Int, toWallets []common.Address, tokenAddr string) {
// 	isNative := (tokenAddr == "")

// 	// get balances, chainId, tokenSymbol
// 	balances := GetBalances(client, fromAddr, totalAmt, toWallets, tokenAddr, isNative)
// 	chainId, err := client.ChainID(context.Background())
// 	if err != nil {
// 		log.Println("getting chainID ", err)
// 	}
// 	tokenSymbol, err := tools.TokenSymbol(client, chainId, tokenAddr)
// 	if err != nil {
// 		log.Println("getting token symbol ", err)
// 	}

// 	// log
// 	log.Println()
// 	log.Println("-----------------------------------------")
// 	log.Println("    Funding wallets with", tokenSymbol)
// 	log.Println("-----------------------------------------")

// 	// Create ew
// 	ew, err := NewExecWallet("./app/wallet_info", pwPath)
// 	if err != nil {
// 		log.Println("NewExecWallet", err)
// 	}

// 	// If we don't have wallet for fromWallet, and user wants to load it, load it
// 	if !(ew.Ks.HasAddress(fromAddr) || loadFromAddr(ew, fromAddr)) {
// 		return
// 	}

// 	// confirm action and calculate amount to each wallet
// 	amountsTo := getAmountsIfCorrect(chainId, totalAmt, tokenSymbol, fromAddr, toWallets, balances)
// 	if amountsTo == nil {
// 		return
// 	}

// 	// tx data
// 	sender := ew.GetAccountByAddr(client, fromAddr)
// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
// 	if err != nil {
// 		log.Println("Getting nonce of fromAcct ", err)
// 	}

// 	gasWei := tools.GweiToWei(decimal.NewFromFloat(gasProvider.Gas.Result.Instant.FeeCap))
// 	gasPrice := gasWei

// 	// If cant afford
// 	if !canAfford(client, fromAddr, totalAmt, toWallets, gasPrice, tokenAddr) {
// 		log.Println("Sender cannot afford txs")
// 		return
// 	}

// 	// Send native token
// 	if isNative {
// 		FundNative(client, sender, nonce, amountsTo, toWallets, gasPrice, chainId)
// 		return
// 	}

// 	// Send ERC20
// 	fundERC20(client, sender, amountsTo, toWallets, gasPrice, common.HexToAddress(tokenAddr), chainId)
// }

func fundERC20(client *ethclient.Client, sender *CachedKey, amountsTo []*big.Int, toWallets []common.Address, gasPrice *big.Int, tokenAddr common.Address, chainId *big.Int) {
	// create token instance
	instance, err := token.NewToken(tokenAddr, client)
	if err != nil {
		log.Println("creating token instance ", err)
	}

	for i, wallet := range toWallets {

		// If wallet doesn't need funds, skip
		if amountsTo[i].Cmp(big.NewInt(0)) == 0 {
			continue
		}

		// create tx
		auth, err := bind.NewKeyedTransactorWithChainID(sender.Key.PrivateKey, chainId)
		if err != nil {
			log.Println("creating bind ops", err)
		}

		// send tx
		tx, err := instance.TokenTransactor.Transfer(auth, wallet, amountsTo[i])
		if err != nil {
			log.Println("creating tx ", err)
		}

		// success
		log.Println("Sent to ", wallet.Hex()[:6], ": ", tx.Hash())
		log.Println("GAS: ", tx.Gas())
	}
}

func GetBalances(client *ethclient.Client, fromAddr common.Address, totalAmt *big.Int, accounts []common.Address, tokenAddr string, isNative bool) []*big.Int {
	var balances []*big.Int

	if isNative {

		// append each wallet balance
		for i := 0; i < len(accounts); i++ {
			bal, err := client.BalanceAt(context.Background(), accounts[i], nil)
			if err != nil {
				log.Println("balanceAt ", err)
				return make([]*big.Int, 0)
			}
			balances = append(balances, bal)
		}
	} else {

		// get contract instance
		instance, err := token.NewToken(common.HexToAddress(tokenAddr), client)
		if err != nil {
			log.Println("creating token instance ", err)
			return make([]*big.Int, 0)
		}

		// append each wallet balance
		for i := 0; i < len(accounts); i++ {
			bal, err := instance.BalanceOf(&bind.CallOpts{}, accounts[i])
			if err != nil {
				log.Println("querying balance ", err)
				return make([]*big.Int, 0)
			}
			balances = append(balances, bal)
		}
	}
	return balances
}

func canAfford(client *ethclient.Client, fromAddr common.Address, totalAmt *big.Int, toWallets []common.Address, gasPrice *big.Int, tokenAddr string) bool {

	// get native balance of sender
	bal, err := client.BalanceAt(context.Background(), fromAddr, nil)
	if err != nil {
		log.Println("balanceAt ", err)
	}

	// if native
	if tokenAddr == "" {

		// sender needs totalAmt + totalGas
		totalGas := new(big.Int).Mul(new(big.Int).Mul(tools.NativeTransferGas(), gasPrice), big.NewInt(int64(len(toWallets))))
		requiredAmt := new(big.Int).Add(totalAmt, totalGas)

		return bal.Cmp(requiredAmt) >= 0
	}
	// else erc20
	totalGas := new(big.Int).Mul(new(big.Int).Mul(tools.ERC20TransferGas(), gasPrice), big.NewInt(int64(len(toWallets))))

	// get erc balance
	instance, err := token.NewToken(common.HexToAddress(tokenAddr), client)
	if err != nil {
		log.Println("creating token instance ", err)
	}
	balERC, err := instance.BalanceOf(&bind.CallOpts{}, fromAddr)
	if err != nil {
		log.Println("querying balance ", err)
	}

	// sender needs totalAmt (erc20) and totalGas (native)
	return bal.Cmp(totalGas) >= 0 && balERC.Cmp(totalAmt) >= 0
}

// func CollectFromWallets(client *ethclient.Client, gasProvider gas_price.GasPricesProvider, toAddr common.Address, pwPath string, fromWallets []common.Address, tokenAddr string) {
// 	isNative := (tokenAddr == "")

// 	// get chainID
// 	chainId, err := client.ChainID(context.Background())
// 	if err != nil {
// 		log.Println("getting chainId: ", err)
// 	}

// 	// Create ew
// 	ew, err := NewExecWallet("../wallet-store", pwPath)
// 	if err != nil {
// 		log.Println("NewExecWallet", err)
// 	}

// 	// trim fromWallets to accounts for which we have pKey
// 	fromAddrs := make([]common.Address, 0)
// 	for _, wallet := range fromWallets {
// 		// If we don't have wallet for fromWallet, and user wants to load it, load it
// 		if ew.Ks.HasAddress(wallet) || loadFromAddr(ew, wallet) {
// 			fromAddrs = append(fromAddrs, wallet)
// 		}
// 	}

// 	// wait for gas price
// 	var gasPrice *big.Int
// 	for {
// 		gasPrice = gasProvider.Get().Fast
// 		if gasPrice == nil {
// 			log.Println("Waiting for first gas update in func")
// 		} else {
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}

// 	// send txs
// 	if isNative {
// 		for _, addr := range fromAddrs {
// 			sender := ew.GetAccountByAddr(client, addr)
// 			sendAllNative(client, sender, toAddr, gasPrice, chainId)
// 		}
// 		return
// 	}
// 	for _, addr := range fromAddrs {
// 		sender := ew.GetAccountByAddr(client, addr)
// 		sendAllERC20(client, sender, toAddr, tokenAddr, gasPrice, chainId)
// 	}
// }

func sendAllERC20(client *ethclient.Client, sender *CachedKey, toAddr common.Address, tokenAddr string, gasPrice *big.Int, chainId *big.Int) {
	log.Println()
	// create instance
	instance, err := token.NewToken(common.HexToAddress(tokenAddr), client)
	if err != nil {
		log.Println("creating token instance ", err)
	}

	// get native balance
	bal, err := client.BalanceAt(context.Background(), sender.Account.Address, nil)
	if err != nil {
		log.Println("balanceAt ", err)
		return
	}

	// confirm wallet can pay gas
	if bal.Cmp(tools.ERC20TransferGas()) <= 0 {
		log.Println("Wallet", sender.Account.Address, "does not have enough gas to return funds")
		return
	}

	// get erc20 balance
	bal, err = instance.BalanceOf(&bind.CallOpts{}, sender.Account.Address)
	if err != nil {
		log.Println("querying balance ", err)
	}
	if bal.Cmp(big.NewInt(0)) <= 0 {
		log.Println("Address:", sender.Account.Address.Hex()[:6], "has balance 0 of", tokenAddr)
		return
	}

	// create tx
	auth, err := bind.NewKeyedTransactorWithChainID(sender.Key.PrivateKey, chainId)
	if err != nil {
		log.Println("creating bind ops", err)
	}

	// send tx
	tx, err := instance.TokenTransactor.Transfer(auth, toAddr, bal)
	if err != nil {
		log.Println("creating tx ", err)
	}

	// success
	log.Println()
	log.Println("Returning", tools.ToDecimal(bal, 6), "tokens from", sender.Account.Address.Hex()[:6], "to", tx.To().Hex()[:6])
	log.Println("| Tx Hash:", tx.Hash())
}

func loadFromAddr(ew *ExecWallet, fromAddr common.Address) bool {
	log.Println("fromAddress [", fromAddr.String(), "] is not recognized..")

	// ask user to import fromAddr
	fmt.Println("Would you like to import the wallet? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	if answer != "y" {
		log.Println("Aborting..")
		return false
	}

	// import fromAddr
	fmt.Println("Enter the private key for the fromAddress: ")
	var key string
	fmt.Scanln(&key)
	err := ew.ImportPrivateKeyFromString(key)
	if err != nil {
		log.Println("importPrivateKeyFromString ", err)
		return false
	} else {
		log.Println("Imported wallet")
		return true
	}
}
