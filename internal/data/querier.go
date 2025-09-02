package data

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Querier struct {
	URL       string
	data      ParsedData
	dataMutex sync.RWMutex
}

func (q *Querier) Query() (*Data, error) {

	resp, err := http.Get(q.URL)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil, err
	}

	return &data, nil
}

// Tick fetches the latest data and updates the internal state
func (q *Querier) Tick() {
	data, err := q.Query()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return
	}
	parsed := ParseData(data)

	q.SetData(parsed)
}

func (q *Querier) GetData() ParsedData {
	q.dataMutex.RLock()
	defer q.dataMutex.RUnlock()
	return q.data
}

func (q *Querier) SetData(data ParsedData) {
	q.dataMutex.Lock()
	defer q.dataMutex.Unlock()
	q.data = data
}
