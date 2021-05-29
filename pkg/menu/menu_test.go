package menu

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
)

func TestNew(t *testing.T) {
	InitText(message.NewPrinter(language.English))
	c := primitive.Config{
		Shape:      primitive.ShapePolygon,
		Iterations: 1000,
		Repeat:     2,
		Alpha:      255,
		Extension:  "png",
		OutputSize: 256,
	}
	ShapesView := NewMenuView(ShapesViewTmpl, fmt.Sprintf("%s/%d", ShapesViewCallback, c.Shape))
	IterView := NewMenuView(IterViewTmpl, fmt.Sprintf("%s/%d", IterViewCallback, c.Iterations))
	RepView := NewMenuView(RepViewTmpl, fmt.Sprintf("%s/%d", RepViewCallback, c.Repeat))
	AlphaView := NewMenuView(AlphaViewTmpl, fmt.Sprintf("%s/%d", AlphaViewCallback, c.Alpha))
	ExtView := NewMenuView(ExtViewTmpl, fmt.Sprintf("%s/%s", ExtViewCallback, c.Extension))
	SizeView := NewMenuView(SizeViewTmpl, fmt.Sprintf("%s/%d", SizeViewCallback, c.OutputSize))

	menu := New(c)

	switch {
	case !reflect.DeepEqual(menu.RootView, RootViewTmpl):
		t.Errorf("RootView: %+v;\n want: %+v", menu.RootView, RootViewTmpl)
	case !reflect.DeepEqual(menu.ShapesView, ShapesView):
		t.Errorf("ShapesView: %+v;\n want: %+v", menu.ShapesView, ShapesView)
	case !reflect.DeepEqual(menu.IterView, IterView):
		t.Errorf("IterView: %+v;\n want: %+v", menu.IterView, IterView)
	case !reflect.DeepEqual(menu.RepView, RepView):
		t.Errorf("RepView: %+v;\n want: %+v", menu.RepView, RepView)
	case !reflect.DeepEqual(menu.AlphaView, AlphaView):
		t.Errorf("AlphaView: %+v;\n want: %+v", menu.AlphaView, AlphaView)
	case !reflect.DeepEqual(menu.ExtView, ExtView):
		t.Errorf("ExtView: %+v;\n want: %+v", menu.ExtView, ExtView)
	case !reflect.DeepEqual(menu.SizeView, SizeView):
		t.Errorf("SizeView: %+v;\n want: %+v", menu.SizeView, SizeView)
	}
}
