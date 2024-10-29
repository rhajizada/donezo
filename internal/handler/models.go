package handler

type BoardRequest struct {
	Name string `json:"name"`
}

type CreateItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateItemRequest struct {
	CreateItemRequest
	Completed bool `json:"completed"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

var DefaultHealthResponse = StatusResponse{
	Status: "healthy",
}

var ValidTokenResponse = StatusResponse{
	Status: "token in valid",
}
