package dht

import (
	"testing"
)

func TestChord(t *testing.T) {
	var port int32 = 1099
	chord, err := NewChord("node1", port)
	if err != nil {
		t.Fatalf("Failed to create new Chord: %v", err)
	}
	println(chord.IsFirst())
	chord2, err := JoinChord("node2", "node1", port)
	if err != nil {
		t.Fatalf("Failed to join Chord: %v", err)
	}
	println(chord.IsFirst())

	// Test Set
	err = chord.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Test Get
	value, err := chord2.Get("key1")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if value != "value1" {
		t.Fatalf("Expected value 'value1', got '%s'", value)
	}

	// Test Delete
	err = chord.Delete("key1")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Test GetAllKeys
	keys, err := chord.GetAllKeys()
	if err != nil {
		t.Fatalf("Failed to get all keys: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("Expected 0 keys, got %d", len(keys))
	}

	// Test IsFirst
	isFirst, err := chord.IsFirst()
	if err != nil {
		t.Fatalf("Failed to check if first: %v", err)
	}
	if !isFirst {
		t.Fatalf("Expected node to be first")
	}
}
