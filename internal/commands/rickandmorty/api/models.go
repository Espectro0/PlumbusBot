package rickandmorty

type Character struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	Species    string   `json:"species"`
	Subspecies string   `json:"type"`
	Gender     string   `json:"gender"`
	Episode    []string `json:"episode"`
	Image      string   `json:"image"`
	Origin     struct {
		Name string `json:"name"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
}

type APIResponse struct {
	Info struct {
		Count int `json:"count"`
		Pages int `json:"pages"`
	} `json:"info"`
	Results []Character `json:"results"`
}
