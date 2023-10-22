package requests

type None struct{}

type RequestWithPublicRSA struct {
	PublicKey []byte `json:"public_rsa"`
}

type GenerateRSAPairRequest = None
type GenerateRSAPairResponse struct {
	PublicKey  []byte `json:"public_rsa"`
	PrivateKey []byte `json:"private_rsa"`
}

type GenerateKeyRequest = RequestWithPublicRSA

type GenerateKeyResponse struct {
	Key []byte `json:"key"`
}

type CreateFileRequest struct {
	RequestWithPublicRSA
	GetFileResponse
}

type GetFileRequest struct {
	RequestWithPublicRSA
	Filename string `json:"filename"`
}

type GetFileResponse struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}
