package db

type ResponseTask struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
}

type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
