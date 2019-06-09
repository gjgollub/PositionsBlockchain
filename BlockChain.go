package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var Blockchain []Block
var templates = template.Must(template.ParseFiles("home.html"))

type Page struct {
	Title string
	Body  []byte
}

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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, "home")
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "home", p)
}

func main() {
	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/home", makeHandler(homeHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
