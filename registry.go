package color

import (
	"strings"
	"sync"
)

// spaceRegistry holds all registered color spaces
var spaceRegistry = struct {
	sync.RWMutex
	spaces map[string]Space
}{
	spaces: make(map[string]Space),
}

func init() {
	// Register all built-in color spaces with their primary names and aliases
	RegisterSpace("srgb", SRGBSpace)
	RegisterSpace("srgb-linear", SRGBLinearSpace)

	RegisterSpace("display-p3", DisplayP3Space)
	RegisterSpace("display-p3-d65", DisplayP3Space) // Alias

	RegisterSpace("dci-p3", DCIP3Space)
	RegisterSpace("dci-p3-d65", DCIP3Space) // Alias

	RegisterSpace("a98-rgb", A98RGBSpace)
	RegisterSpace("a98rgb", A98RGBSpace) // Alias
	RegisterSpace("adobe-rgb-1998", A98RGBSpace) // Alias

	RegisterSpace("prophoto-rgb", ProPhotoRGBSpace)
	RegisterSpace("prophoto", ProPhotoRGBSpace) // Alias

	RegisterSpace("rec2020", Rec2020Space)
	RegisterSpace("rec-2020", Rec2020Space) // Alias

	RegisterSpace("rec709", Rec709Space)
	RegisterSpace("rec-709", Rec709Space) // Alias

	RegisterSpace("oklch", OKLCHSpace)

	// LOG color spaces for professional cinema cameras
	RegisterSpace("c-log", CLogSpace)
	RegisterSpace("clog", CLogSpace) // Alias

	RegisterSpace("s-log3", SLog3Space)
	RegisterSpace("slog3", SLog3Space) // Alias

	RegisterSpace("v-log", VLogSpace)
	RegisterSpace("vlog", VLogSpace) // Alias

	RegisterSpace("arri-logc", ArriLogCSpace)
	RegisterSpace("logc", ArriLogCSpace) // Alias

	RegisterSpace("red-log3g10", RedLog3G10Space)
	RegisterSpace("log3g10", RedLog3G10Space) // Alias

	RegisterSpace("bmd-film", BMDFilmSpace)
	RegisterSpace("bmdfilm", BMDFilmSpace) // Alias
}

// RegisterSpace registers a color space with the given name(s).
// The name is case-insensitive. You can register the same space with multiple names.
//
// Example:
//   RegisterSpace("my-space", myCustomSpace)
//   RegisterSpace("my-alias", myCustomSpace)  // Register with alias
func RegisterSpace(name string, space Space) {
	spaceRegistry.Lock()
	defer spaceRegistry.Unlock()
	spaceRegistry.spaces[strings.ToLower(name)] = space
}

// GetSpace retrieves a registered color space by name.
// The name lookup is case-insensitive.
// Returns nil if the space is not found.
//
// Example:
//   space, ok := GetSpace("display-p3")
//   if ok {
//       spaceColor := NewSpaceColor(space, []float64{1, 0, 0}, 1.0)
//   }
func GetSpace(name string) (Space, bool) {
	spaceRegistry.RLock()
	defer spaceRegistry.RUnlock()
	space, ok := spaceRegistry.spaces[strings.ToLower(name)]
	return space, ok
}

// ListSpaces returns a list of all registered color space names.
// The names are returned in no particular order.
//
// Example:
//   spaces := ListSpaces()
//   for _, name := range spaces {
//       fmt.Println(name)
//   }
func ListSpaces() []string {
	spaceRegistry.RLock()
	defer spaceRegistry.RUnlock()

	names := make([]string, 0, len(spaceRegistry.spaces))
	for name := range spaceRegistry.spaces {
		names = append(names, name)
	}
	return names
}

// UnregisterSpace removes a color space from the registry.
// This is mainly useful for testing or when dynamically managing color spaces.
func UnregisterSpace(name string) {
	spaceRegistry.Lock()
	defer spaceRegistry.Unlock()
	delete(spaceRegistry.spaces, strings.ToLower(name))
}
