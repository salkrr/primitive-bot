package menu

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestInitTextInitializesShapeNames(t *testing.T) {
	InitText(message.NewPrinter(language.English))

	if len(ShapeNames) != 9 {
		t.Errorf("ShapeNames length is %d; want %d", len(ShapeNames), 9)
	}
	for i, v := range ShapeNames {
		if v == "" {
			t.Errorf("ShapeNames[%v] is empty string.", i)
		}
	}
}

func TestInitTextInitializesViewTemplates(t *testing.T) {
	InitText(message.NewPrinter(language.English))

	tests := []struct {
		name     string
		template View
	}{
		{
			name:     "Initializes RootViewTmpl",
			template: RootViewTmpl,
		},
		{
			name:     "Initializes ShapeViewTmpl",
			template: ShapesViewTmpl,
		},
		{
			name:     "Initializes IterViewTmpl",
			template: IterViewTmpl,
		},
		{
			name:     "Initializes RepViewTmpl",
			template: RepViewTmpl,
		},
		{
			name:     "Initializes AlphaViewTmpl",
			template: AlphaViewTmpl,
		},
		{
			name:     "Initializes ExtViewTmpl",
			template: ExtViewTmpl,
		},
		{
			name:     "Initializes SizeViewTmpl",
			template: SizeViewTmpl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.template.Text == "" {
				t.Error("Template's text is empty.")
			}

			for _, row := range tt.template.Keyboard.InlineKeyboard {
				for _, button := range row {
					if button.Text == "" {
						t.Errorf("Button with the callback %v has empty text", button.CallbackData)
					}
				}
			}
		})
	}
}
