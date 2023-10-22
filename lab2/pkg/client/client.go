package client

import (
	"bytes"
	"encoding/json"
	"lab2/pkg/crypt"
	"lab2/pkg/requests"
	"net/http"
)

func New() *App {
	app := new(App)
	http.HandleFunc("/key", app.GetSessionKey)
	http.HandleFunc("/gen", app.GenerateRSAPair)
	http.HandleFunc("/create", app.CreateFile)
	return app
}

type App struct {
	SessionKey []byte
	PrivateKey []byte
	PublicKey  []byte
}

func (a *App) Serve(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func (a *App) GenerateRSAPair(w http.ResponseWriter, r *http.Request) {
	privateKey, publicKey, err := crypt.GenerateRSAPair()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	resp := requests.GenerateRSAPairResponse{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	a.PrivateKey = privateKey
	a.PublicKey = publicKey

	encoder := json.NewEncoder(w)

	encoder.Encode(resp)
}

func (a *App) GetSessionKey(w http.ResponseWriter, r *http.Request) {
	request := requests.GenerateKeyRequest{
		PublicKey: a.PublicKey,
	}
	marshaled, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := http.Post("http://localhost:8081/gen", http.DetectContentType(marshaled), bytes.NewBuffer(marshaled))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	var response requests.GenerateKeyResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	a.SessionKey = response.Key
	encoder := json.NewEncoder(w)

	encoder.Encode(requests.GenerateKeyResponse{
		Key: response.Key,
	})
}

func (a *App) CreateFile(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateFileRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	decryptedKey, err := crypt.Decrypt(a.PrivateKey, a.SessionKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedData, err := crypt.EncryptAES(decryptedKey, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Content = encryptedData
	req.PublicKey = a.SessionKey
	marshaled, err := json.Marshal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = http.Post("http://localhost:8081/create", http.DetectContentType(marshaled), bytes.NewBuffer(marshaled))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
