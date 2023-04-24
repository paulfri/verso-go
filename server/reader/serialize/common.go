package serialize

type EmptyObject struct{}

type Self struct {
	Href string `json:"href"`
}

type Category struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
