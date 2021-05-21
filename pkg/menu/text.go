package menu

import "github.com/lazy-void/primitive-bot/pkg/primitive"

var ShapeNames = map[primitive.Shape]string{
	primitive.ShapeAny:              "Все",
	primitive.ShapeTriangle:         "Треугольники",
	primitive.ShapeRectangle:        "Прямоугольники",
	primitive.ShapeRotatedRectangle: "Повёрнутые прямоугольники",
	primitive.ShapeCircle:           "Круги",
	primitive.ShapeEllipse:          "Эллипсы",
	primitive.ShapeRotatedEllipse:   "Повёрнутые эллипсы",
	primitive.ShapePolygon:          "Четырёхугольники",
	primitive.ShapeBezier:           "Кривые Безье",
}

const (
	rootMenuText   = "Меню:"
	shapesMenuText = "Выбери фигуры, из которых будет выстраиваться изображение:"
	iterMenuText   = "Выбери количество итераций - шагов, на каждом из которых будет отрисовываться фигуры:"
	repMenuText    = "Выбери сколько фигур будет отрисовываться на каждой итерации:"
	alphaMenuText  = "Выбери значение альфа-канала каждой отрисовываемой фигуры:"
	extMenuText    = "Выбери расширение файла:"
	sizeMenuText   = "Выбери размер для большей стороны изображения (соотношение сторон будет сохранено):"
)

const (
	createButtonText = "Начать"
	backButtonText   = "Назад"
	shapesButtonText = "Фигуры"
	iterButtonText   = "Итерации"
	repButtonText    = "Повторения"
	alphaButtonText  = "Альфа"
	extButtonText    = "Расширение"
	sizeButtonText   = "Размеры"
	autoButtonText   = "Автоматически"
	OtherButtonText  = "Другое"
)
