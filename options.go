package scrapegoat

type Options struct {
	MaxRecursionDepth int
	EnableConcurrency bool
	EnableLogging     bool
}

var DefaultOptions = Options{
	MaxRecursionDepth: 3,
	EnableConcurrency: false,
	EnableLogging:     false,
}
