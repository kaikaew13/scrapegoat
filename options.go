package scrapegoat

const defaultMaxScrapingDepth uint = 3

var defaultOptions = options{
	maxScrapingDepth:  defaultMaxScrapingDepth,
	curScrapingDepth:  0,
	enableConcurrency: false,
	enableLogging:     false,
}

type options struct {
	maxScrapingDepth  uint
	curScrapingDepth  uint
	enableConcurrency bool
	enableLogging     bool
}

type optionFunc func(opts *options)

func MaxScrapingDepth(depth uint) optionFunc {
	return func(opts *options) {
		opts.maxScrapingDepth = depth
	}
}

func EnableConcurrency(b bool) optionFunc {
	return func(opts *options) {
		opts.enableConcurrency = b
	}
}

func EnableLogging(b bool) optionFunc {
	return func(opts *options) {
		opts.enableLogging = b
	}
}
