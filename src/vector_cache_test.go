package polecalc

import "testing"

// When a value is set in the cache, can that value be found and retrieved?
func TestVectorSetSticks(t *testing.T) {
	key := Vector2{0.0, 0.0}
	value := Vector2{10.0, 10.0}
	cache := NewVectorCache()
	cache.Set(key, value)
	if !cache.Contains(key) {
		t.Fatalf("VectorCache does not contain value added to it")
	}
	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("VectorCache failed to find value given to it")
	}
	gotVector, ok := got.(Vector2)
	if !ok {
		t.Fatalf("VectorCache returned a value with the wrong type")
	}
	if !value.Equals(gotVector) {
		t.Fatalf("VectorCache returned the wrong value")
	}
}
