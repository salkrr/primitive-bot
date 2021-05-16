package main

import (
	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

var shapeNames = map[primitive.Shape]string{
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
	combo             = "Все"
	triangles         = "Треугольники"
	rectangles        = "Прямоугольники"
	rotatedRectangles = "Повёрнутые прямоугольники"
	circles           = "Круги"
	ellipses          = "Эллипсы"
	rotatedEllipses   = "Повёрнутые эллипсы"
	quadrilaterals    = "Четырёхугольники"
	bezierCurves      = "Кривые Безье"
)

const (
	helpMessage     = "Отправь мне какую-нибудь фотографию."
	errorMessage    = "Что-то пошло не так! Попробуй снова через пару минут."
	enqueuedMessage = "Добавил в очередь.\nФигуры: %s\nИтерации: %d\nПовторения: %d\nАльфа-канал: %d\nРасширение: %s\nРазмеры: %d"
)

const (
	rootMenu     = "Меню:"
	settingsMenu = "Настройки:"
	shapesMenu   = "Выбери фигуры, из которых будет выстраиваться изображение:"
	iterMenu     = "Выбери количество итераций - шагов, на каждом из которых будет отрисовываться фигуры:"
	repMenu      = "Выбери сколько фигур будет отрисовываться на каждой итерации:"
	alphaMenu    = "Выбери значение альфа-канала каждой отрисовываемой фигуры:"
	extMenu      = "Выбери расширение файла:"
	sizeMenu     = "Выбери размер для большей стороны изображения (соотношение сторон будет сохранено):"
)

const (
	startButton    = "Начать"
	settingsButton = "Настройки"
	backButton     = "Назад"
	shapesButton   = "Фигуры"
	iterButton     = "Итерации"
	repButton      = "Повторения"
	alphaButton    = "Альфа"
	extButton      = "Расширение"
	sizeButton     = "Размеры"
)

var (
	rootKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: startButton, CallbackData: "/start"},
			},
			{
				{Text: shapesButton, CallbackData: "/settings/shape"},
				{Text: iterButton, CallbackData: "/settings/iter"},
				{Text: repButton, CallbackData: "/settings/rep"},
			},
			{
				{Text: alphaButton, CallbackData: "/settings/alpha"},
				{Text: extButton, CallbackData: "/settings/ext"},
				{Text: sizeButton, CallbackData: "/settings/size"},
			},
		},
	}

	shapesKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{{Text: shapeNames[primitive.ShapeAny], CallbackData: "/settings/shape/0"}},
			{{Text: triangles, CallbackData: "/settings/shape/1"}},
			{{Text: rectangles, CallbackData: "/settings/shape/2"}},
			{{Text: ellipses, CallbackData: "/settings/shape/3"}},
			{{Text: circles, CallbackData: "/settings/shape/4"}},
			{{Text: rotatedRectangles, CallbackData: "/settings/shape/5"}},
			{{Text: bezierCurves, CallbackData: "/settings/shape/6"}},
			{{Text: rotatedEllipses, CallbackData: "/settings/shape/7"}},
			{{Text: quadrilaterals, CallbackData: "/settings/shape/8"}},
			{{Text: backButton, CallbackData: "/"}},
		},
	}

	iterKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "100", CallbackData: "/settings/iter/100"},
				{Text: "200", CallbackData: "/settings/iter/200"},
				{Text: "400", CallbackData: "/settings/iter/400"},
			},
			{
				{Text: "800", CallbackData: "/settings/iter/800"},
				{Text: "1000", CallbackData: "/settings/iter/1000"},
				{Text: "2000", CallbackData: "/settings/iter/2000"},
			},
			{
				{Text: backButton, CallbackData: "/"},
			},
		},
	}

	repKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "1", CallbackData: "/settings/rep/1"},
				{Text: "2", CallbackData: "/settings/rep/2"},
				{Text: "3", CallbackData: "/settings/rep/3"},
			},
			{
				{Text: "4", CallbackData: "/settings/rep/4"},
				{Text: "5", CallbackData: "/settings/rep/5"},
				{Text: "6", CallbackData: "/settings/rep/6"},
			},
			{
				{Text: backButton, CallbackData: "/"},
			},
		},
	}

	alphaKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "Автоматически", CallbackData: "/settings/alpha/0"},
			},
			{
				{Text: "32", CallbackData: "/settings/alpha/32"},
				{Text: "64", CallbackData: "/settings/alpha/64"},
				{Text: "128", CallbackData: "/settings/alpha/128"},
				{Text: "255", CallbackData: "/settings/alpha/255"},
			},
			{
				{Text: backButton, CallbackData: "/"},
			},
		},
	}

	extKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "jpg", CallbackData: "/settings/ext/jpg"},
				{Text: "png", CallbackData: "/settings/ext/png"},
				{Text: "svg", CallbackData: "/settings/ext/svg"},
				{Text: "gif", CallbackData: "/settings/ext/gif"},
			},
			{
				{Text: backButton, CallbackData: "/"},
			},
		},
	}

	sizeKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "256", CallbackData: "/settings/size/256"},
				{Text: "512", CallbackData: "/settings/size/512"},
				{Text: "720", CallbackData: "/settings/size/720"},
			},
			{
				{Text: "1024", CallbackData: "/settings/size/1024"},
				{Text: "1280", CallbackData: "/settings/size/1280"},
				{Text: "1920", CallbackData: "/settings/size/1920"},
			},
			{
				{Text: backButton, CallbackData: "/"},
			},
		},
	}
)
