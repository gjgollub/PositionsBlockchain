package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Blockchain []Block

type Block struct {
	Index         int
	TimeStamp     string
	Security_Id   string
	Security_Type string
	Quantity      int
	PreviousHash  string
	Hash          string
}

func calculateHash(block Block) string {
	record := (string(block.Index) + block.TimeStamp + block.Security_Id + block.Security_Type + string(block.Quantity) + block.PreviousHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, secId string, secType string, quantity int) (Block, error) {

	var result Block

	result.Index = (oldBlock.Index + 1)
	result.TimeStamp = time.Now().String()
	result.Security_Id = secId
	result.Security_Type = secType
	result.Quantity = quantity
	result.Hash = calculateHash(oldBlock)
	result.PreviousHash = oldBlock.Hash

	return result, nil
}

func validateBlock(oldBlock Block, newBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}