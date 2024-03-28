package main

import (
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type entry struct {
	key   string
	value   interface{}
	expiration time.Time
}

type LRUCache struct {
	cache   map[string]*list.Element
	lruList *list.List
	maxSize   int
	expiration  time.Duration
	mutex      sync.Mutex

}
var cache = NewLRUCache(1024, 5*time.Second)
func GetHandler(w http.ResponseWriter, r *http.Request){
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key parameter missing", http.StatusBadRequest)
		return
	}
	value, found := cache.Get(key)
	if !found {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(value)
}
type Data struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
}
func SetHandler(w http.ResponseWriter, r *http.Request){
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data);
	err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	cache.Set(data.Key, data.Value)
	w.WriteHeader(http.StatusCreated)
}
func NewLRUCache(maxSize int, expiration time.Duration) *LRUCache {
    return &LRUCache{
        cache:      make(map[string]*list.Element),
        lruList:    list.New(),
        maxSize:    maxSize,
        expiration: expiration,
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if ele, found := c.cache[key]; found {
        if ele.Value.(*entry).expiration.Before(time.Now()) {
            c.removeElement(ele)
            return nil, false
        }
        c.lruList.MoveToFront(ele)
        return ele.Value.(*entry).value, true
    }
    return nil, false
}
func (c *LRUCache) Set(key string, value interface{}) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if ele, found := c.cache[key]; found {
        c.lruList.MoveToFront(ele)
        ele.Value.(*entry).value = value
        ele.Value.(*entry).expiration = time.Now().Add(c.expiration)
    } else {
        ele := c.lruList.PushFront(&entry{key, value, time.Now().Add(c.expiration)})
        c.cache[key] = ele

        if len(c.cache) > c.maxSize {
            c.removeOldest()
        }
    }
}

func (c *LRUCache) removeOldest() {
    ele := c.lruList.Back()
    if ele != nil {
        c.removeElement(ele)
    }
}

func (c *LRUCache) removeElement(e *list.Element) {
    c.lruList.Remove(e)
    delete(c.cache, e.Value.(*entry).key)
}

func main(){
	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/set", SetHandler)
	log.Fatal(http.ListenAndServe(":8080",nil))
}