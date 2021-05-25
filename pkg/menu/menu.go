// Package menu implements functionality of a telegram menu.
package menu

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

// View represents menu 'window' that is
// made with the telegram's inline keyboard.
type View struct {
	Text     string
	Keyboard tg.InlineKeyboardMarkup
}

// Menu represents menu made of View instances.
type Menu struct {
	RootView   View
	ShapesView View
	IterView   View
	RepView    View
	AlphaView  View
	ExtView    View
	SizeView   View
}

// New initializes instance of Menu.
func New(c primitive.Config) Menu {
	shapesCallback := fmt.Sprintf("%s/%d", ShapesViewCallback, c.Shape)
	iterCallback := fmt.Sprintf("%s/%d", IterViewCallback, c.Iterations)
	repCallback := fmt.Sprintf("%s/%d", RepViewCallback, c.Repeat)
	alphaCallback := fmt.Sprintf("%s/%d", AlphaViewCallback, c.Alpha)
	extCallback := fmt.Sprintf("%s/%s", ExtViewCallback, c.Extension)
	sizeCallback := fmt.Sprintf("%s/%d", SizeViewCallback, c.OutputSize)

	return Menu{
		RootView:   NewMenuView(RootViewTmpl, ""),
		ShapesView: NewMenuView(ShapesViewTmpl, shapesCallback),
		IterView:   NewMenuView(IterViewTmpl, iterCallback),
		RepView:    NewMenuView(RepViewTmpl, repCallback),
		AlphaView:  NewMenuView(AlphaViewTmpl, alphaCallback),
		ExtView:    NewMenuView(ExtViewTmpl, extCallback),
		SizeView:   NewMenuView(SizeViewTmpl, sizeCallback),
	}
}
