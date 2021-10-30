package http

type sluggableErr interface {
	Slug() string
}

type Error struct {
	Message string
	Slug    string
}

type CreateSobRequest struct {
	Description string
	Id          string
	Name        string
}

type SobResponse struct {
	Description string
	Id          string
	Name        string
}
