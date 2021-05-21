package menu

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

var (
	rootKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: createButtonText, CallbackData: CreateButtonCallback},
			},
			{
				{Text: shapesButtonText, CallbackData: ShapesButtonCallback},
				{Text: iterButtonText, CallbackData: IterButtonCallback},
			},
			{
				{Text: repButtonText, CallbackData: RepButtonCallback},
				{Text: alphaButtonText, CallbackData: AlphaButtonCallback},
			},
			{
				{Text: extButtonText, CallbackData: ExtButtonCallback},
				{Text: sizeButtonText, CallbackData: SizeButtonCallback},
			},
		},
	}

	shapesKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{
					Text:         ShapeNames[primitive.ShapeAny],
					CallbackData: fmt.Sprintf("%s/0", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeTriangle],
					CallbackData: fmt.Sprintf("%s/1", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRectangle],
					CallbackData: fmt.Sprintf("%s/2", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeEllipse],
					CallbackData: fmt.Sprintf("%s/3", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeCircle],
					CallbackData: fmt.Sprintf("%s/4", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedRectangle],
					CallbackData: fmt.Sprintf("%s/5", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeBezier],
					CallbackData: fmt.Sprintf("%s/6", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedEllipse],
					CallbackData: fmt.Sprintf("%s/7", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapePolygon],
					CallbackData: fmt.Sprintf("%s/8", ShapesButtonCallback),
				},
			},
			{
				{
					Text:         backButtonText,
					CallbackData: RootMenuCallback,
				},
			},
		},
	}

	iterKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "100", CallbackData: fmt.Sprintf("%s/100", IterButtonCallback)},
				{Text: "200", CallbackData: fmt.Sprintf("%s/200", IterButtonCallback)},
				{Text: "400", CallbackData: fmt.Sprintf("%s/400", IterButtonCallback)},
			},
			{
				{Text: "800", CallbackData: fmt.Sprintf("%s/800", IterButtonCallback)},
				{Text: "1000", CallbackData: fmt.Sprintf("%s/1000", IterButtonCallback)},
				{Text: "2000", CallbackData: fmt.Sprintf("%s/2000", IterButtonCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: IterInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootMenuCallback},
			},
		},
	}

	repKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "1", CallbackData: fmt.Sprintf("%s/1", RepButtonCallback)},
				{Text: "2", CallbackData: fmt.Sprintf("%s/2", RepButtonCallback)},
				{Text: "3", CallbackData: fmt.Sprintf("%s/3", RepButtonCallback)},
			},
			{
				{Text: "4", CallbackData: fmt.Sprintf("%s/4", RepButtonCallback)},
				{Text: "5", CallbackData: fmt.Sprintf("%s/5", RepButtonCallback)},
				{Text: "6", CallbackData: fmt.Sprintf("%s/6", RepButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootMenuCallback},
			},
		},
	}

	alphaKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: autoButtonText, CallbackData: fmt.Sprintf("%s/0", AlphaButtonCallback)},
			},
			{
				{Text: "32", CallbackData: fmt.Sprintf("%s/32", AlphaButtonCallback)},
				{Text: "64", CallbackData: fmt.Sprintf("%s/64", AlphaButtonCallback)},
				{Text: "128", CallbackData: fmt.Sprintf("%s/128", AlphaButtonCallback)},
				{Text: "255", CallbackData: fmt.Sprintf("%s/255", AlphaButtonCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: AlphaInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootMenuCallback},
			},
		},
	}

	extKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "jpg", CallbackData: fmt.Sprintf("%s/jpg", ExtButtonCallback)},
				{Text: "png", CallbackData: fmt.Sprintf("%s/png", ExtButtonCallback)},
				{Text: "svg", CallbackData: fmt.Sprintf("%s/svg", ExtButtonCallback)},
				// gifs are disabled due to performance issues
				// {Text: "gif", CallbackData: fmt.Sprintf("%s/gif", extButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootMenuCallback},
			},
		},
	}

	sizeKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "256", CallbackData: fmt.Sprintf("%s/256", SizeButtonCallback)},
				{Text: "512", CallbackData: fmt.Sprintf("%s/512", SizeButtonCallback)},
				{Text: "720", CallbackData: fmt.Sprintf("%s/720", SizeButtonCallback)},
			},
			{
				{Text: "1024", CallbackData: fmt.Sprintf("%s/1024", SizeButtonCallback)},
				{Text: "1280", CallbackData: fmt.Sprintf("%s/1280", SizeButtonCallback)},
				{Text: "1920", CallbackData: fmt.Sprintf("%s/1920", SizeButtonCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: SizeInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootMenuCallback},
			},
		},
	}
)

var (
	RootActivityTmpl = Activity{
		Text:     rootMenuText,
		Keyboard: rootKeyboardTmpl,
	}

	ShapesActivityTmpl = Activity{
		Text:     shapesMenuText,
		Keyboard: shapesKeyboardTmpl,
	}

	IterActivityTmpl = Activity{
		Text:     iterMenuText,
		Keyboard: iterKeyboardTmpl,
	}

	RepActivityTmpl = Activity{
		Text:     repMenuText,
		Keyboard: repKeyboardTmpl,
	}

	AlphaActivityTmpl = Activity{
		Text:     alphaMenuText,
		Keyboard: alphaKeyboardTmpl,
	}

	ExtActivityTmpl = Activity{
		Text:     extMenuText,
		Keyboard: extKeyboardTmpl,
	}

	SizeActivityTmpl = Activity{
		Text:     sizeMenuText,
		Keyboard: sizeKeyboardTmpl,
	}
)

// NewMenuActivity creates new Activity from the template
// adding symbol to the option that is chosen at the moment
func NewMenuActivity(
	template Activity,
	selectedCallback string,
	newButtonText ...string,
) Activity {
	checkSymbol := "👉"

	// Create new Keyboard
	newKeyboard := telegram.InlineKeyboardMarkup{}
	newKeyboard.InlineKeyboard = make([][]telegram.InlineKeyboardButton, len(template.Keyboard.InlineKeyboard))
	for i, row := range template.Keyboard.InlineKeyboard {
		newKeyboard.InlineKeyboard[i] = make([]telegram.InlineKeyboardButton, len(row))
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

	return Activity{
		Text:     template.Text,
		Keyboard: newKeyboard,
	}
}
