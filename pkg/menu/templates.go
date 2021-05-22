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
				{Text: shapesButtonText, CallbackData: ShapesActivityCallback},
				{Text: iterButtonText, CallbackData: IterActivityCallback},
			},
			{
				{Text: repButtonText, CallbackData: RepActivityCallback},
				{Text: alphaButtonText, CallbackData: AlphaActivityCallback},
			},
			{
				{Text: extButtonText, CallbackData: ExtActivityCallback},
				{Text: sizeButtonText, CallbackData: SizeActivityCallback},
			},
		},
	}

	shapesKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{
					Text:         ShapeNames[primitive.ShapeAny],
					CallbackData: fmt.Sprintf("%s/0", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeTriangle],
					CallbackData: fmt.Sprintf("%s/1", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRectangle],
					CallbackData: fmt.Sprintf("%s/2", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeEllipse],
					CallbackData: fmt.Sprintf("%s/3", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeCircle],
					CallbackData: fmt.Sprintf("%s/4", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedRectangle],
					CallbackData: fmt.Sprintf("%s/5", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeBezier],
					CallbackData: fmt.Sprintf("%s/6", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapeRotatedEllipse],
					CallbackData: fmt.Sprintf("%s/7", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         ShapeNames[primitive.ShapePolygon],
					CallbackData: fmt.Sprintf("%s/8", ShapesActivityCallback),
				},
			},
			{
				{
					Text:         backButtonText,
					CallbackData: RootActivityCallback,
				},
			},
		},
	}

	iterKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "100", CallbackData: fmt.Sprintf("%s/100", IterActivityCallback)},
				{Text: "200", CallbackData: fmt.Sprintf("%s/200", IterActivityCallback)},
				{Text: "400", CallbackData: fmt.Sprintf("%s/400", IterActivityCallback)},
			},
			{
				{Text: "800", CallbackData: fmt.Sprintf("%s/800", IterActivityCallback)},
				{Text: "1000", CallbackData: fmt.Sprintf("%s/1000", IterActivityCallback)},
				{Text: "2000", CallbackData: fmt.Sprintf("%s/2000", IterActivityCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: IterInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootActivityCallback},
			},
		},
	}

	repKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "1", CallbackData: fmt.Sprintf("%s/1", RepActivityCallback)},
				{Text: "2", CallbackData: fmt.Sprintf("%s/2", RepActivityCallback)},
				{Text: "3", CallbackData: fmt.Sprintf("%s/3", RepActivityCallback)},
			},
			{
				{Text: "4", CallbackData: fmt.Sprintf("%s/4", RepActivityCallback)},
				{Text: "5", CallbackData: fmt.Sprintf("%s/5", RepActivityCallback)},
				{Text: "6", CallbackData: fmt.Sprintf("%s/6", RepActivityCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootActivityCallback},
			},
		},
	}

	alphaKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: autoButtonText, CallbackData: fmt.Sprintf("%s/0", AlphaActivityCallback)},
			},
			{
				{Text: "32", CallbackData: fmt.Sprintf("%s/32", AlphaActivityCallback)},
				{Text: "64", CallbackData: fmt.Sprintf("%s/64", AlphaActivityCallback)},
				{Text: "128", CallbackData: fmt.Sprintf("%s/128", AlphaActivityCallback)},
				{Text: "255", CallbackData: fmt.Sprintf("%s/255", AlphaActivityCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: AlphaInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootActivityCallback},
			},
		},
	}

	extKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "jpg", CallbackData: fmt.Sprintf("%s/jpg", ExtActivityCallback)},
				{Text: "png", CallbackData: fmt.Sprintf("%s/png", ExtActivityCallback)},
				{Text: "svg", CallbackData: fmt.Sprintf("%s/svg", ExtActivityCallback)},
				// gifs are disabled due to performance issues
				// {Text: "gif", CallbackData: fmt.Sprintf("%s/gif", extButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: RootActivityCallback},
			},
		},
	}

	sizeKeyboardTmpl = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "256", CallbackData: fmt.Sprintf("%s/256", SizeActivityCallback)},
				{Text: "512", CallbackData: fmt.Sprintf("%s/512", SizeActivityCallback)},
				{Text: "720", CallbackData: fmt.Sprintf("%s/720", SizeActivityCallback)},
			},
			{
				{Text: "1024", CallbackData: fmt.Sprintf("%s/1024", SizeActivityCallback)},
				{Text: "1280", CallbackData: fmt.Sprintf("%s/1280", SizeActivityCallback)},
				{Text: "1920", CallbackData: fmt.Sprintf("%s/1920", SizeActivityCallback)},
			},
			{
				{Text: OtherButtonText, CallbackData: SizeInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: RootActivityCallback},
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
