package main

import "sync"


type rdb struct {
    mu      sync.RWMutex
    rmap    map[string]string
}

// GET take string and return string or nil based on given key
func (db *rdb) GET(key string) string {
    db.mu.RLock()
    defer db.mu.RUnlock()

    return db.rmap[key]
}

// SET take key value string and insert/update new item
// using given key and value pair
func (db *rdb) SET(key,value string) {
    db.mu.Lock()
    defer db.mu.Unlock()
    db.rmap[key] = value
}

// DEL take string key and return false if item in given key is not exist
// otherwise delete item in db with given key and return true
func (db *rdb) DEL(key string) bool {
    db.mu.Lock()
    defer db.mu.Unlock()

    lengthBefore := len(db.rmap)
    delete(db.rmap, key)

    if lengthBefore == len(db.rmap) {
        return false
    }

    return true
}


func CreateAndSeed() *rdb {
    data := map[string]string {
        "HP 1": "Harry Potter and the Philosopher's Stone",
        "HP 2": "Harry Potter and the Chamber of Secrets",
        "HP 3": "Harry Potter and the Prisoner Of Azkabat",
        "HP 4": "Harry Potter and the Goblet of Fire",
        "HP 5": "Harry Potter and the Order of the Phoenix",
        "HP 6": "Harry Potter and the Half-Blood Prince",
        "HP 7": "Harry Potter and the Deathly Hallows - Part 1",
        "HP 8": "Harry Potter and the Deathly Hallows - Part 2",
    }

    db := rdb{}
    db.rmap = data

    return &db
}
