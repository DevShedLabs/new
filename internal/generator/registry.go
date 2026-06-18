package generator

import "fmt"

var registered []Generator

// Register adds a generator to the global registry.
// Call this from each generator package's init() or explicitly in main.
func Register(g Generator) {
	registered = append(registered, g)
}

// Resolve walks the registry in order and returns the first Generator
// whose Detect method returns true for the given Context.
// Returns an error if no generator matches.
func Resolve(ctx *Context) (Generator, error) {
	for _, g := range registered {
		if g.Detect(ctx) {
			return g, nil
		}
	}
	return nil, fmt.Errorf("no generator found for %q", ctx.Name)
}
