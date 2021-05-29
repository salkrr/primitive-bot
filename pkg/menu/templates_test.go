package menu

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/lazy-void/primitive-bot/pkg/tg"
)

func copyKeyboard(k tg.InlineKeyboardMarkup) tg.InlineKeyboardMarkup {
	newKeyboard := tg.InlineKeyboardMarkup{}
	newKeyboard.InlineKeyboard = make([][]tg.InlineKeyboardButton, len(k.InlineKeyboard))
	for i, row := range k.InlineKeyboard {
		newKeyboard.InlineKeyboard[i] = make([]tg.InlineKeyboardButton, len(row))
		for j, button := range row {
			newKeyboard.InlineKeyboard[i][j] = button
		}
	}

	return newKeyboard
}

func TestNewMenuView(t *testing.T) {
	InitText(message.NewPrinter(language.English))

	tests := []struct {
		name     string
		template View
	}{
		{
			name:     "RootViewTmpl",
			template: RootViewTmpl,
		},
		{
			name:     "ShapeViewTmpl",
			template: ShapesViewTmpl,
		},
		{
			name:     "IterViewTmpl",
			template: IterViewTmpl,
		},
		{
			name:     "RepViewTmpl",
			template: RepViewTmpl,
		},
		{
			name:     "AlphaViewTmpl",
			template: AlphaViewTmpl,
		},
		{
			name:     "ExtViewTmpl",
			template: ExtViewTmpl,
		},
		{
			name:     "SizeViewTmpl",
			template: SizeViewTmpl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, row := range tt.template.Keyboard.InlineKeyboard {
				for j, button := range row {
					// not changing button text
					callback := button.CallbackData
					expected := copyKeyboard(tt.template.Keyboard)
					expected.InlineKeyboard[i][j].Text = fmt.Sprintf("👉 %s", button.Text)

					res := NewMenuView(tt.template, callback)
					if res.Text != tt.template.Text {
						t.Errorf("Got menu view text: %s; want: %s", res.Text, tt.template.Text)
					}
					if !reflect.DeepEqual(res.Keyboard, expected) {
						t.Errorf("Got InlineKeyboard: %+v;\n want: %+v", res.Keyboard, expected)
					}

					// changing button text
					newText := "New Text"
					expected.InlineKeyboard[i][j].Text = fmt.Sprintf("👉 %s", newText)

					res = NewMenuView(tt.template, callback, newText)
					if res.Text != tt.template.Text {
						t.Errorf("Got menu view text: %s; want: %s", res.Text, tt.template.Text)
					}
					if !reflect.DeepEqual(res.Keyboard, expected) {
						t.Errorf("Got InlineKeyboard: %+v;\n want: %+v", res.Keyboard, expected)
					}
				}
			}
		})
	}
}

func TestNewMenuViewWhenGivenIncorrectCallback(t *testing.T) {
	InitText(message.NewPrinter(language.English))
	callback := "some random text"
	template := RootViewTmpl

	res := NewMenuView(template, callback)
	if !reflect.DeepEqual(res, template) {
		t.Errorf("Got menu view: %+v;\n want: %+v", res, template)
	}
}
