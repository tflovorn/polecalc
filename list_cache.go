package polecalc

// To be used when the expected number of keys is small and the cost of
// calling reflect.DeepEqual on keys isn't too high.

import "reflect"

type ListCache struct {
	keys, values []interface{}
}

func NewListCache() *ListCache {
	ls := new(ListCache)
	ls.keys = []interface{}{}
	ls.values = []interface{}{}
	return ls
}

func (ls *ListCache) Contains(key interface{}) bool {
	if ls.indexOf(key) > -1 {
		return true
	}
	return false
}

func (ls *ListCache) Get(key interface{}) (interface{}, bool) {
	if !ls.Contains(key) {
		return nil, false
	}
	i := ls.indexOf(key)
	return ls.values[i], true
}

func (ls *ListCache) Set(key interface{}, value interface{}) {
	if !ls.Contains(key) {
		// add key to cache
		ls.keys = append(ls.keys, key)
		ls.values = append(ls.values, value)
		return
	}
	// if we get here, key is already in the cache
	i := ls.indexOf(key)
	ls.values[i] = value
}

func (ls *ListCache) indexOf(key interface{}) int {
	for i, cachedKey := range ls.keys {
		if reflect.DeepEqual(key, cachedKey) {
			return i
		}
	}
	return -1
}
