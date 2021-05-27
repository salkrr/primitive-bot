package menu

import (
	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"golang.org/x/text/message"
)

// ShapeNames contains mapping of shapes to their string representation.
var ShapeNames map[primitive.Shape]string

var (
	rootMenuText   string
	shapesMenuText string
	iterMenuText   string
	repMenuText    string
	alphaMenuText  string
	extMenuText    string
	sizeMenuText   string
)

// Text of buttons in the menu.
var (
	createButtonText string
	backButtonText   string
	shapesButtonText string
	iterButtonText   string
	repButtonText    string
	alphaButtonText  string
	extButtonText    string
	sizeButtonText   string
	autoButtonText   string
	OtherButtonText  string
)

// InitText initializes all global variables that contain text
// and also keyboard and view templates.
func InitText(p message.Printer) {
	ShapeNames = map[primitive.Shape]string{
		primitive.ShapeAny:              p.Sprintf("All"),
		primitive.ShapeTriangle:         p.Sprintf("Triangles"),
		primitive.ShapeRectangle:        p.Sprintf("Rectangles"),
		primitive.ShapeRotatedRectangle: p.Sprintf("Rotated Rectangles"),
		primitive.ShapeCircle:           p.Sprintf("Circles"),
		primitive.ShapeEllipse:          p.Sprintf("Ellipses"),
		primitive.ShapeRotatedEllipse:   p.Sprintf("Rotated Ellipses"),
		primitive.ShapePolygon:          p.Sprintf("Quadrilaterals"),
		primitive.ShapeBezier:           p.Sprintf("Bezier Curves"),
	}

	createButtonText = p.Sprintf("Create")
	backButtonText = p.Sprintf("Back")
	shapesButtonText = p.Sprintf("Shapes")
	iterButtonText = p.Sprintf("Steps")
	repButtonText = p.Sprintf("Repetitions")
	alphaButtonText = p.Sprintf("Alpha")
	extButtonText = p.Sprintf("Extension")
	sizeButtonText = p.Sprintf("Size")
	autoButtonText = p.Sprintf("Auto")
	OtherButtonText = p.Sprintf("Other")

	rootMenuText = p.Sprintf("Menu:")
	shapesMenuText = p.Sprintf("Select the shapes to be used to create the image:")
	iterMenuText = p.Sprintf("Select the number of steps. Shapes will be drawn at each step:")
	repMenuText = p.Sprintf("Select the number of shapes to draw in each step:")
	alphaMenuText = p.Sprintf("Select an alpha-channel value for the shapes:")
	extMenuText = p.Sprintf("Select an extension of the resulting image:")
	sizeMenuText = p.Sprintf("Select a size for the larger side of the resulting image (the aspect ratio will be preserved):")

	initKeyboardTemplates()
	initViewTemplates()
}
