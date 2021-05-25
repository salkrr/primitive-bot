package menu

import "fmt"

// Callbacks that are sent by menu buttons.
var (
	RootViewCallback     = "/"
	CreateButtonCallback = "/create"

	ShapesViewCallback   = "/shape"
	ShapesButtonCallback = fmt.Sprintf("%s/([0-8])", ShapesViewCallback)

	IterViewCallback   = "/iter"
	IterButtonCallback = fmt.Sprintf("%s/([0-9]+)", IterViewCallback)
	IterInputCallback  = "/iter/input"

	RepViewCallback   = "/rep"
	RepButtonCallback = fmt.Sprintf("%s/([1-6])", RepViewCallback)

	AlphaViewCallback   = "/alpha"
	AlphaButtonCallback = fmt.Sprintf("%s/([0-9]+)", AlphaViewCallback)
	AlphaInputCallback  = "/alpha/input"

	ExtViewCallback   = "/ext"
	ExtButtonCallback = fmt.Sprintf("%s/(jpg|png|svg|gif)", ExtViewCallback)

	SizeViewCallback   = "/size"
	SizeButtonCallback = fmt.Sprintf("%s/([0-9]+)", SizeViewCallback)
	SizeInputCallback  = "/size/input"
)
