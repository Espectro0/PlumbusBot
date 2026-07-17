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

type CharacterAPIResponse struct {
	Info    Info        `json:"info"`
	Results []Character `json:"results"`
}

type Location struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Dimension string   `json:"dimension"`
	Residents []string `json:"residents"`
}

type LocationAPIResponse struct {
	Info    Info       `json:"info"`
	Results []Location `json:"results"`
}

type Info struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}
