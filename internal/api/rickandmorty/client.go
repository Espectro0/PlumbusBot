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

type namedEntry[T any] struct {
	data      []T
	expiresAt time.Time
}

type namedCache[T any] struct {
	mu    sync.RWMutex
	cache map[string]namedEntry[T]
}

func newNamedCache[T any]() *namedCache[T] {
	return &namedCache[T]{
		cache: make(map[string]namedEntry[T]),
	}
}

func (c *namedCache[T]) get(key string) ([]T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

func (c *namedCache[T]) set(key string, data []T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = namedEntry[T]{
		data:      data,
		expiresAt: time.Now().Add(cacheTTL),
	}
}

type totalCache struct {
	mu        sync.RWMutex
	count     int
	expiresAt time.Time
}

func (c *totalCache) get() (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().Before(c.expiresAt) {
		return c.count, true
	}
	return 0, false
}

func (c *totalCache) set(count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = count
	c.expiresAt = time.Now().Add(cacheTTL)
}

type pageResponse[T any] struct {
	Info    Info `json:"info"`
	Results []T  `json:"results"`
}

var (
	charNameCache = newNamedCache[Character]()
	locNameCache  = newNamedCache[Location]()

	characterTotal = &totalCache{}
	locationTotal  = &totalCache{}
)

func doGetJSON(url string, dest any) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dest)
}

func fetchByID[T any](url string, kind string, id int) (*T, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s %d not found", kind, id)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fetchAll[T any](url string) ([]T, error) {
	var all []T
	for url != "" {
		var resp pageResponse[T]
		if err := doGetJSON(url, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Results...)
		url = resp.Info.Next
	}
	return all, nil
}

func getTotal(endpoint string, cache *totalCache) (int, error) {
	if count, ok := cache.get(); ok {
		return count, nil
	}

	var resp pageResponse[any]
	if err := doGetJSON(baseURL+endpoint, &resp); err != nil {
		return 0, err
	}

	cache.set(resp.Info.Count)
	return resp.Info.Count, nil
}

func GetTotalCharacter() (int, error) {
	return getTotal("/character", characterTotal)
}

func GetCharacterById(id int) (*Character, error) {
	url := fmt.Sprintf("%s/character/%d", baseURL, id)
	return fetchByID[Character](url, "Character", id)
}

func GetCharacterByName(name string) ([]Character, error) {
	if cached, ok := charNameCache.get(name); ok {
		return cached, nil
	}

	url := fmt.Sprintf("%s/character/?name=%s", baseURL, name)
	results, err := fetchAll[Character](url)
	if err != nil {
		return nil, fmt.Errorf("Character %q not found", name)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Character %q not found", name)
	}

	charNameCache.set(name, results)
	return results, nil
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

func GetTotalLocation() (int, error) {
	return getTotal("/location", locationTotal)
}

func GetLocationById(id int) (*Location, error) {
	url := fmt.Sprintf("%s/location/%d", baseURL, id)
	return fetchByID[Location](url, "Location", id)
}

func GetLocationByName(name string) ([]Location, error) {
	if cached, ok := locNameCache.get(name); ok {
		return cached, nil
	}

	url := fmt.Sprintf("%s/location/?name=%s", baseURL, name)
	results, err := fetchAll[Location](url)
	if err != nil {
		return nil, fmt.Errorf("Location %q not found", name)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Location %q not found", name)
	}

	locNameCache.set(name, results)
	return results, nil
}

func RandomLocation() (*Location, error) {
	count, err := GetTotalLocation()
	if err != nil {
		return nil, fmt.Errorf("Failed to get total: %w", err)
	}

	for attempt := 0; attempt < 5; attempt++ {
		id := rand.IntN(count) + 1
		loc, err := GetLocationById(id)
		if err == nil {
			return loc, nil
		}
	}

	return GetLocationById(1)
}
