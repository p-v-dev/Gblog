package tag

type CreateTagInput struct {
	Name string `json:"name"`
}

type TagOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
