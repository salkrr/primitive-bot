package menu

import "fmt"

var (
	RootActivityCallback = "/"
	CreateButtonCallback = "/create"

	ShapesActivityCallback = "/shape"
	ShapesButtonCallback   = fmt.Sprintf("%s/([0-8])", ShapesActivityCallback)

	IterActivityCallback = "/iter"
	IterButtonCallback   = fmt.Sprintf("%s/([0-9]+)", IterActivityCallback)
	IterInputCallback    = "/iter/input"

	RepActivityCallback = "/rep"
	RepButtonCallback   = fmt.Sprintf("%s/([1-6])", RepActivityCallback)

	AlphaActivityCallback = "/alpha"
	AlphaButtonCallback   = fmt.Sprintf("%s/([0-9]+)", AlphaActivityCallback)
	AlphaInputCallback    = "/alpha/input"

	ExtActivityCallback = "/ext"
	ExtButtonCallback   = fmt.Sprintf("%s/(jpg|png|svg|gif)", ExtActivityCallback)

	SizeActivityCallback = "/size"
	SizeButtonCallback   = fmt.Sprintf("%s/([0-9]+)", SizeActivityCallback)
	SizeInputCallback    = "/size/input"
)
