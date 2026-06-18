package generator

// Context holds all input parsed from the CLI for a generation request.
type Context struct {
	// Name is the file or project name passed as the first argument.
	Name string
	// Template is the --template flag value (empty for file generation).
	Template string
	// OutputDir is where output should be written (defaults to cwd).
	OutputDir string
	// Vars holds any additional key=value pairs passed via --var flags.
	Vars map[string]string
}

// Generator is implemented by anything that can produce files from a Context.
type Generator interface {
	// Name returns the unique identifier for this generator.
	Name() string
	// Detect returns true if this generator should handle the given Context.
	Detect(ctx *Context) bool
	// Generate executes the generation and returns any error.
	Generate(ctx *Context) error
}
