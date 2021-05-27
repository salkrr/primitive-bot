package menu

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

var (
	rootKeyboardTmpl   tg.InlineKeyboardMarkup
	shapesKeyboardTmpl tg.InlineKeyboardMarkup
	iterKeyboardTmpl   tg.InlineKeyboardMarkup
	repKeyboardTmpl    tg.InlineKeyboardMarkup
	alphaKeyboardTmpl  tg.InlineKeyboardMarkup
	extKeyboardTmpl    tg.InlineKeyboardMarkup
	sizeKeyboardTmpl   tg.InlineKeyboardMarkup
)

// Templates for the different menu views.
var (
	RootViewTmpl   View
	ShapesViewTmpl View
	IterViewTmpl   View
	RepViewTmpl    View
	AlphaViewTmpl  View
	ExtViewTmpl    View
	SizeViewTmpl   View
)

func initKeyboardTemplates() {
	rootKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: createButtonText, CallbackData: CreateButtonCallback},
			},
			{
				{Text: shapesButtonText, CallbackData: ShapesViewCallback},
				{Text: iterButtonText, CallbackData: IterViewCallback},
			},
			{
				{Text: repButtonText, CallbackData: RepViewCallback},
				{Text: alphaButtonText, CallbackData: AlphaViewCallback},
			},
			{
				{Text: extButtonText, CallbackData: ExtViewCallback},
				{Text: sizeButtonText, CallbackData: SizeViewCallback},
			},
		},
	}

	shapesKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{
					Text:         ShapeNames[primitive.ShapeAny],
					CallbackData: fmt.Sprintf("%s/0", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeTriangle],
					CallbackData: fmt.Sprintf("%s/1", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapePolygon],
					CallbackData: fmt.Sprintf("%s/8", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRectangle],
					CallbackData: fmt.Sprintf("%s/2", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedRectangle],
					CallbackData: fmt.Sprintf("%s/5", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeEllipse],
					CallbackData: fmt.Sprintf("%s/3", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedEllipse],
					CallbackData: fmt.Sprintf("%s/7", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeCircle],
					CallbackData: fmt.Sprintf("%s/4", ShapesViewCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeBezier],
					CallbackData: fmt.Sprintf("%s/6", ShapesViewCallback),
				},
			},
			{
				{
					Text:         backButtonText,
					CallbackData: RootViewCallback,
				},
			},
		},
	}

	iterKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: "100", CallbackData: fmt.Sprintf("%s/100", IterViewCallback)},
				{Text: "200", CallbackData: fmt.Sprintf("%s/200", IterViewCallback)},
				{Text: "400", CallbackData: fmt.Sprintf("%s/400", IterViewCallback)},
			},
			{
				{Text: "800", CallbackData: fmt.Sprintf("%s/800", IterViewCallback)},
				{Text: "1000", CallbackData: fmt.Sprintf("%s/1000", IterViewCallback)},
				{Text: "2000", CallbackData: fmt.Sprintf("%s/2000", IterViewCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: IterInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootViewCallback},
			},
		},
	}

	repKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: "1", CallbackData: fmt.Sprintf("%s/1", RepViewCallback)},
				{Text: "2", CallbackData: fmt.Sprintf("%s/2", RepViewCallback)},
				{Text: "3", CallbackData: fmt.Sprintf("%s/3", RepViewCallback)},
			},
			{
				{Text: "4", CallbackData: fmt.Sprintf("%s/4", RepViewCallback)},
				{Text: "5", CallbackData: fmt.Sprintf("%s/5", RepViewCallback)},
				{Text: "6", CallbackData: fmt.Sprintf("%s/6", RepViewCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootViewCallback},
			},
		},
	}

	alphaKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: autoButtonText, CallbackData: fmt.Sprintf("%s/0", AlphaViewCallback)},
			},
			{
				{Text: "32", CallbackData: fmt.Sprintf("%s/32", AlphaViewCallback)},
				{Text: "64", CallbackData: fmt.Sprintf("%s/64", AlphaViewCallback)},
				{Text: "128", CallbackData: fmt.Sprintf("%s/128", AlphaViewCallback)},
				{Text: "255", CallbackData: fmt.Sprintf("%s/255", AlphaViewCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: AlphaInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootViewCallback},
			},
		},
	}

	extKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: "jpg", CallbackData: fmt.Sprintf("%s/jpg", ExtViewCallback)},
				{Text: "png", CallbackData: fmt.Sprintf("%s/png", ExtViewCallback)},
				{Text: "svg", CallbackData: fmt.Sprintf("%s/svg", ExtViewCallback)},
				// gifs are disabled due to performance issues
				// {Text: "gif", CallbackData: fmt.Sprintf("%s/gif", extButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootViewCallback},
			},
		},
	}

	sizeKeyboardTmpl = tg.InlineKeyboardMarkup{
		InlineKeyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: "256", CallbackData: fmt.Sprintf("%s/256", SizeViewCallback)},
				{Text: "512", CallbackData: fmt.Sprintf("%s/512", SizeViewCallback)},
				{Text: "720", CallbackData: fmt.Sprintf("%s/720", SizeViewCallback)},
			},
			{
				{Text: "1024", CallbackData: fmt.Sprintf("%s/1024", SizeViewCallback)},
				{Text: "1280", CallbackData: fmt.Sprintf("%s/1280", SizeViewCallback)},
				{Text: "1920", CallbackData: fmt.Sprintf("%s/1920", SizeViewCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: SizeInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootViewCallback},
			},
		},
	}
}

func initViewTemplates() {
	RootViewTmpl = View{
		Text:     rootMenuText,
		Keyboard: rootKeyboardTmpl,
	}

	ShapesViewTmpl = View{
		Text:     shapesMenuText,
		Keyboard: shapesKeyboardTmpl,
	}

	IterViewTmpl = View{
		Text:     iterMenuText,
		Keyboard: iterKeyboardTmpl,
	}

	RepViewTmpl = View{
		Text:     repMenuText,
		Keyboard: repKeyboardTmpl,
	}

	AlphaViewTmpl = View{
		Text:     alphaMenuText,
		Keyboard: alphaKeyboardTmpl,
	}

	ExtViewTmpl = View{
		Text:     extMenuText,
		Keyboard: extKeyboardTmpl,
	}

	SizeViewTmpl = View{
		Text:     sizeMenuText,
		Keyboard: sizeKeyboardTmpl,
	}
}

// NewMenuView creates new View from the template
// adding symbol to the option that is chosen at the moment.
func NewMenuView(
	template View,
	selectedCallback string,
	newButtonText ...string,
) View {
	checkSymbol := "👉"

	// Create new Keyboard
	newKeyboard := tg.InlineKeyboardMarkup{}
	newKeyboard.InlineKeyboard = make([][]tg.InlineKeyboardButton, len(template.Keyboard.InlineKeyboard))
	for i, row := range template.Keyboard.InlineKeyboard {
		newKeyboard.InlineKeyboard[i] = make([]tg.InlineKeyboardButton, len(row))
		for j, button := range row {
			newKeyboard.InlineKeyboard[i][j] = button

			if button.CallbackData == selectedCallback {
				newText := button.Text
				if len(newButtonText) > 0 {
					newText = newButtonText[0]
				}

				newKeyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("%s %s", checkSymbol, newText)
			}
		}
	}

	return View{
		Text:     template.Text,
		Keyboard: newKeyboard,
	}
}
