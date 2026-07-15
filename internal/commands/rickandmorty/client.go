package rickandmorty

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
)

const baseURL = "https://rickandmortyapi.com/api"

func GetTotalCharacter() (int, error) {
	resp, err := http.Get(baseURL + "/character")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, err
	}

	return apiResp.Info.Count, nil
}

func GetCharacterById(id int) (*Character, error) {
	url := fmt.Sprintf("%s/character/%d", baseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Character %d not found", id)
	}

	var char Character
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		return nil, err
	}

	return &char, nil
}

func RandomCharacter() (*Character, error) {
	count, err := GetTotalCharacter()
	if err != nil {
		return nil, fmt.Errorf("Failed to get total: %w", err)
	}

	for attempt := 0; attempt < 5; attempt++ {
		id := rand.IntN(count) + 1
		char, err := GetCharacterById(id)
		if err == nil {
			return char, nil
		}
	}

	_ = count
	return GetCharacterById(1)
}
