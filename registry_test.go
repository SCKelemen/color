package color

import (
	"strings"
	"testing"
)

func TestRegisterAndGetSpace(t *testing.T) {
	// Test registering and retrieving a space
	space := OKLCHSpace

	RegisterSpace("test-space", space)

	retrieved, ok := GetSpace("test-space")
	if !ok {
		t.Error("GetSpace failed to retrieve registered space")
	}
	if retrieved.Name() != space.Name() {
		t.Errorf("Retrieved space name %s, want %s", retrieved.Name(), space.Name())
	}
}

func TestGetSpaceCaseInsensitive(t *testing.T) {
	RegisterSpace("Test-Space-Case", OKLCHSpace)

	// Should work with different cases
	cases := []string{"test-space-case", "TEST-SPACE-CASE", "Test-Space-Case"}
	for _, name := range cases {
		_, ok := GetSpace(name)
		if !ok {
			t.Errorf("GetSpace(%q) failed, should be case-insensitive", name)
		}
	}
}

func TestGetSpaceNotFound(t *testing.T) {
	_, ok := GetSpace("non-existent-space-12345")
	if ok {
		t.Error("GetSpace should return false for non-existent space")
	}
}

func TestListSpaces(t *testing.T) {
	// Register some test spaces
	RegisterSpace("test1", OKLCHSpace)
	RegisterSpace("test2", DisplayP3Space)

	spaces := ListSpaces()

	// Check that registered spaces are present
	found := make(map[string]bool)
	for _, name := range spaces {
		found[strings.ToLower(name)] = true
	}

	if !found["srgb"] {
		t.Error("ListSpaces should include 'srgb'")
	}
	if !found["display-p3"] {
		t.Error("ListSpaces should include 'display-p3'")
	}
}

func TestListSpacesReturnsAll(t *testing.T) {
	spaces := ListSpaces()

	// Should have at least the core spaces
	expectedSpaces := []string{"srgb", "display-p3", "adobe-rgb", "prophoto-rgb", "rec2020"}

	spaceSet := make(map[string]bool)
	for _, s := range spaces {
		spaceSet[strings.ToLower(s)] = true
	}

	for _, expected := range expectedSpaces {
		if !spaceSet[expected] {
			t.Errorf("ListSpaces missing expected space: %s", expected)
		}
	}
}

func TestRegisterSpaceOverwrite(t *testing.T) {
	// Registering the same name twice should overwrite
	RegisterSpace("overwrite-test", OKLCHSpace)
	RegisterSpace("overwrite-test", DisplayP3Space)

	space, ok := GetSpace("overwrite-test")
	if !ok {
		t.Error("GetSpace failed after overwrite")
	}
	if space.Name() != DisplayP3Space.Name() {
		t.Error("Space was not overwritten")
	}
}

func BenchmarkGetSpace(b *testing.B) {
	RegisterSpace("bench-space", OKLCHSpace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSpace("bench-space")
	}
}

func BenchmarkListSpaces(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ListSpaces()
	}
}
