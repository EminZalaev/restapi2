package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"restapi2/internal/model"
	"sync"
)

type CacheBench struct {
	store map[int]model.BenchFullResult
	m     sync.Mutex
}

func NewCache() *CacheBench {
	ss := &CacheBench{}
	ss.m = sync.Mutex{}
	ss.store = make(map[int]model.BenchFullResult)
	return ss
}

func (ss *CacheBench) Load() error {
	log.Print("cache is loaded")
	data, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		return fmt.Errorf("error read json file: %w", err)
	}

	tmp := model.BenchFullResult{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("error unmarshall json file: %w", err)
	}

	ss.m.Lock()
	defer ss.m.Unlock()

	id := len(ss.store)
	for {
		if reflect.DeepEqual(ss.store[id], model.BenchFullResult{}) {
			ss.store[id] = tmp
			return nil
		}
		id++
	}

}

func (ss *CacheBench) Get(id int) *model.BenchFullResult {
	return &model.BenchFullResult{
		Metrics: ss.store[id].Metrics,
	}
}
