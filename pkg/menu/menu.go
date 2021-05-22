package menu

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

type Activity struct {
	Text     string
	Keyboard tg.InlineKeyboardMarkup
}

type Menu struct {
	RootActivity   Activity
	ShapesActivity Activity
	IterActivity   Activity
	RepActivity    Activity
	AlphaActivity  Activity
	ExtActivity    Activity
	SizeActivity   Activity
}

func New(c primitive.Config) Menu {
	shapesCallback := fmt.Sprintf("%s/%d", ShapesActivityCallback, c.Shape)
	iterCallback := fmt.Sprintf("%s/%d", IterActivityCallback, c.Iterations)
	repCallback := fmt.Sprintf("%s/%d", RepActivityCallback, c.Repeat)
	alphaCallback := fmt.Sprintf("%s/%d", AlphaActivityCallback, c.Alpha)
	extCallback := fmt.Sprintf("%s/%s", ExtActivityCallback, c.Extension)
	sizeCallback := fmt.Sprintf("%s/%d", SizeActivityCallback, c.OutputSize)

	return Menu{
		RootActivity:   NewMenuActivity(RootActivityTmpl, ""),
		ShapesActivity: NewMenuActivity(ShapesActivityTmpl, shapesCallback),
		IterActivity:   NewMenuActivity(IterActivityTmpl, iterCallback),
		RepActivity:    NewMenuActivity(RepActivityTmpl, repCallback),
		AlphaActivity:  NewMenuActivity(AlphaActivityTmpl, alphaCallback),
		ExtActivity:    NewMenuActivity(ExtActivityTmpl, extCallback),
		SizeActivity:   NewMenuActivity(SizeActivityTmpl, sizeCallback),
	}
}
