package tools

import (
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
)

// log tx to csv
func WriteToCSV(obj []string, filepath string) {

	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Can't find bookKeeping file ", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(obj)
	if err != nil {
		log.Println("Failed to record transaction ", err)
	}
}

// checks for folder in project directory
func IsFolderEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func PrettyPrint(tx *types.Transaction) {
	msg, err := tx.AsMessage(types.NewEIP155Signer(big.NewInt(1)), tx.GasFeeCap())
	if err != nil {
		log.Println("AnalyzeTx(), ERR in AsMessage()", err)
		return
	}
	log.Println("-----------------------------------------------")
	log.Println("Analyzing the transaction that is to be sent..")
	log.Println("-----------------------------------------------")
	log.Println("Hash:", tx.Hash())
	log.Println("From:", msg.From().Hex())
	log.Println("To:", tx.To())
	log.Println("Nonce:", tx.Nonce())
	log.Println("ChainID:", tx.ChainId())
	log.Println()
	log.Println("Value:", WeiToEth(tx.Value()).String(), "ETH")
	log.Println("Size:", tx.Size())
	log.Println("Cost:", WeiToEth(tx.Cost()).String(), "ETH")
	log.Println()
	log.Println("Gas Limit:", tx.Gas(), "units")
	log.Println("Gas Fee Cap:", WeiToGwei(tx.GasFeeCap()).String(), "Gwei")
	log.Println("Gas Tip Cap:", WeiToGwei(tx.GasTipCap()).String(), "Gwei")
	log.Println("Gas Price:", WeiToGwei(tx.GasPrice()).String(), "Gwei")
	log.Println("--------------------------")
	log.Println("Data:", hex.EncodeToString(tx.Data()))
	// log.Println("--------------------------")
	// v, r, s := tx.RawSignatureValues()
	// log.Println("Raw Signature Values: v|", v, " r|", r, " s|", s)
	log.Println("-----------------------------------------------")
}
