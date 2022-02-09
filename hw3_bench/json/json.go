package json

type user struct {
	Browsers []string `json:"browsers,intern"`
	Company  string   `json:"company,intern"`
	Country  string   `json:"country,intern"`
	Email    string   `json:"email,intern"`
	Job      string   `json:"job,intern"`
	Name     string   `json:"name,intern"`
	Phone    string   `json:"phone,intern"`
}
