package dto

type BookQueryDTO struct {
	ID    int64  `json:"id"`
	StrID string `json:"str_id"`
	Query string `json:"query"`
	Page  int64  `json:"page"`
	Sort  string `json:"sort"`
}

type BookTitle struct {
	Eng    string `json:"eng"`
	JP     string `json:"jp"`
	Pretty string `json:"pretty"`
}

type BookTag struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int64  `json:"count"`
}

type BookPage struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"heigth"`
}

type BookResponse struct {
	ID         int64       `json:"id"`
	MediaID    int64       `json:"media_id"`
	Title      BookTitle   `json:"title"`
	Favorites  int64       `json:"favorites"`
	Thumbnail  string      `json:"thumbnail"`
	Cover      string      `json:"cover"`
	Scanlator  string      `json:"scanlator"`
	Uploaded   string      `json:"uploaded"`
	EpochTime  int64       `json:"epoch_time"`
	Characters []*BookTag  `json:"characters"`
	Pages      []*BookPage `json:"pages"`
	Tags       []string    `json:"tags"`
	NumPages   int64       `json:"num_pages"`
	RawTags    []*BookTag  `json:"-"`
}

type BooksResponse struct {
	Books []*BookResponse `json:"books"`
}
