package pages

type Page struct {
	ID       string
	Title    string
	Author   string
	MenuName string
	Slug     string
	Kind     Kind
	Pics     []Pic
}

type Pic struct {
	ID   string
	Path string
	Alt  string
	Post Page
}

type Kind struct {
	ID    string
	Name  string
	Posts []Page
}
