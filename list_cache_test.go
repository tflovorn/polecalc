package polecalc

import "testing"

func TestListSetSticks(t *testing.T) {
	key := Vector2{0.0, 0.0}
	value := Vector2{10.0, 10.0}
	cache := NewListCache()
	cache.Set(key, value)
	if !cache.Contains(key) {
		t.Fatalf("ListCache does not contain value added to it")
	}
	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("ListCache failed to find value given to it")
	}
	gotVector, ok := got.(Vector2)
	if !ok {
		t.Fatalf("ListCache returned a value with the wrong type")
	}
	if !value.Equals(gotVector) {
		t.Fatalf("ListCache returned the wrong value")
	}
}

func TestListBadSearch(t *testing.T) {
	key := Vector2{0.0, 0.0}
	cache := NewListCache()
	if cache.Contains(key) {
		t.Fatalf("ListCache reports it contains a key not given to it")
	}
}
