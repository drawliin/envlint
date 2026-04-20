package report

// Options controls how the final audit result should be rendered.
type Options struct {
	JSON           bool
	Strict         bool
	EnvFile        string
	ExampleEnvFile string
}
