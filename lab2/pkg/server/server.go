package server

import (
	"crypto/rand"
	"encoding/json"
	"lab2/pkg/crypt"
	"lab2/pkg/requests"
	"log"
	"net/http"
)

const keyLength = 32

func New() *App {
	app := &App{
		files: make(map[string][]byte),
	}
	http.HandleFunc("/gen", app.GenerateKey)
	http.HandleFunc("/create", app.CreateFile)
	return app
}

type App struct {
	files map[string][]byte
}

func (a *App) Serve(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func (a *App) GenerateKey(w http.ResponseWriter, r *http.Request) {
	var req requests.GenerateKeyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	randomString := make([]byte, keyLength)
	if _, err := rand.Read(randomString); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("Generated key: %s\n", randomString)

	key, err := crypt.Encrypt(req.PublicKey, randomString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)

	encoder.Encode(requests.GenerateKeyResponse{
		Key: key,
	})
}

func (a *App) CreateFile(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateFileRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	a.files[req.Filename] = req.Content
}

func (a *App) GetFile(w http.ResponseWriter, r *http.Request) {
	w.Write(nil)
}
