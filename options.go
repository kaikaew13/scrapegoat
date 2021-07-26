package scrapegoat

type Options struct {
	MaxRecursionDepth int
	EnableConcurrency bool
}

var DefaultOptions = Options{
	MaxRecursionDepth: 3,
	EnableConcurrency: false,
}
