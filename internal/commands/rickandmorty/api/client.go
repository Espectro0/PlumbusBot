package rickandmorty

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

const baseURL = "https://rickandmortyapi.com/api"

const cacheTTL = 5 * time.Minute

var httpClient = &http.Client{Timeout: 5 * time.Second}

var (
	cacheMu      sync.RWMutex
	nameCache    = make(map[string]nameEntry)
	totalCount   int
	totalExpires time.Time
	totalMu      sync.RWMutex
)

type nameEntry struct {
	chars     []Character
	expiresAt time.Time
}

func doGetJSON(url string, dest any) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dest)
}

func GetTotalCharacter() (int, error) {
	totalMu.RLock()
	if time.Now().Before(totalExpires) {
		totalMu.RUnlock()
		return totalCount, nil
	}
	totalMu.RUnlock()

	var apiResp APIResponse
	if err := doGetJSON(baseURL+"/character", &apiResp); err != nil {
		return 0, err
	}

	totalMu.Lock()
	totalCount = apiResp.Info.Count
	totalExpires = time.Now().Add(cacheTTL)
	totalMu.Unlock()

	return apiResp.Info.Count, nil
}

func GetCharacterById(id int) (*Character, error) {
	url := fmt.Sprintf("%s/character/%d", baseURL, id)

	resp, err := httpClient.Get(url)
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

func GetCharacterByName(name string) ([]Character, error) {
	key := name

	cacheMu.RLock()
	entry, ok := nameCache[key]
	cacheMu.RUnlock()
	if ok && time.Now().Before(entry.expiresAt) {
		return entry.chars, nil
	}

	url := fmt.Sprintf("%s/character/?name=%s", baseURL, name)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Character %q not found", name)
	}

	var apiResp struct {
		Results []Character `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Results) == 0 {
		return nil, fmt.Errorf("Character %q not found", name)
	}

	cacheMu.Lock()
	nameCache[key] = nameEntry{
		chars:     apiResp.Results,
		expiresAt: time.Now().Add(cacheTTL),
	}
	cacheMu.Unlock()

	return apiResp.Results, nil
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

	return GetCharacterById(1)
}
