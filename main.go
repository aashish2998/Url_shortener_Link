package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {

	hasher := md5.New()

	hasher.Write([]byte(OriginalURL)) // it converts the original url to byte slice

	data := hasher.Sum(nil)

	fmt.Println("print the the hasher data value", data)

	hash := hex.EncodeToString(data)

	fmt.Println("Encode to string ", hash)
	fmt.Println("final string ", hash[:8])

	return hash[:8]
}

func createURL(originalURL string) string {
	shortUrl := generateShortURL(originalURL)
	id := shortUrl // Use the short url as id for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortUrl:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("Error not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "wassup People")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body ", http.StatusBadRequest)
		return
	}

	shortUrl := createURL(data.URL)
	//fmt.Fprintf(w, shortUrl)

	response := struct {
		ShortUrl string `json:"short_url"`
	}{ShortUrl: shortUrl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
		return // Stop further execution
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	// fmt.Println(" Creating url shortner ")
	// OriginalURL := "https://github.com/aashish2998"
	// generateShortURL(OriginalURL)

	//Register to handler function to handle all requestes to the root URL
	http.HandleFunc("/", handler)
	http.HandleFunc("/shortner", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	//start HTTP server on PORT 3000
	fmt.Println("Starting the server on port 3000")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("There is error on the server", err)
	}
}
