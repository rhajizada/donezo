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
