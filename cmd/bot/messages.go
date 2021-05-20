package main

import (
	"fmt"

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
	helpMessage            = "Отправь мне какую-нибудь фотографию."
	errorMessage           = "Что-то пошло не так! Попробуй снова через пару минут."
	inputMessage           = "Неверное значение!\nВведи число от %d до %d:"
	statusMessage          = "%d место в очереди.\n\nФигуры: %s\nИтерации: %d\nПовторения: %d\nАльфа-канал: %d\nРасширение: %s\nРазмеры: %d"
	statusEmptyMessage     = "Нету активных операций."
	operationsLimitMessage = "Вы не можете добавить больше операций в очередь."
)

const (
	rootMenuText     = "Меню:"
	settingsMenuText = "Настройки:"
	shapesMenuText   = "Выбери фигуры, из которых будет выстраиваться изображение:"
	iterMenuText     = "Выбери количество итераций - шагов, на каждом из которых будет отрисовываться фигуры:"
	repMenuText      = "Выбери сколько фигур будет отрисовываться на каждой итерации:"
	alphaMenuText    = "Выбери значение альфа-канала каждой отрисовываемой фигуры:"
	extMenuText      = "Выбери расширение файла:"
	sizeMenuText     = "Выбери размер для большей стороны изображения (соотношение сторон будет сохранено):"
)

const (
	createButtonText   = "Начать"
	settingsButtonText = "Настройки"
	backButtonText     = "Назад"
	shapesButtonText   = "Фигуры"
	iterButtonText     = "Итерации"
	repButtonText      = "Повторения"
	alphaButtonText    = "Альфа"
	extButtonText      = "Расширение"
	sizeButtonText     = "Размеры"
	otherButtonText    = "Другое"
	autoButtonText     = "Автоматически"
)

var (
	rootKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: createButtonText, CallbackData: createButtonCallback},
			},
			{
				{Text: shapesButtonText, CallbackData: shapesButtonCallback},
				{Text: iterButtonText, CallbackData: iterButtonCallback},
			},
			{
				{Text: repButtonText, CallbackData: repButtonCallback},
				{Text: alphaButtonText, CallbackData: alphaButtonCallback},
			},
			{
				{Text: extButtonText, CallbackData: extButtonCallback},
				{Text: sizeButtonText, CallbackData: sizeButtonCallback},
			},
		},
	}

	shapesKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{{Text: shapeNames[primitive.ShapeAny], CallbackData: fmt.Sprintf("%s/0", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeTriangle], CallbackData: fmt.Sprintf("%s/1", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeRectangle], CallbackData: fmt.Sprintf("%s/2", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeEllipse], CallbackData: fmt.Sprintf("%s/3", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeCircle], CallbackData: fmt.Sprintf("%s/4", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeRotatedRectangle], CallbackData: fmt.Sprintf("%s/5", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeBezier], CallbackData: fmt.Sprintf("%s/6", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapeRotatedEllipse], CallbackData: fmt.Sprintf("%s/7", shapesButtonCallback)}},
			{{Text: shapeNames[primitive.ShapePolygon], CallbackData: fmt.Sprintf("%s/8", shapesButtonCallback)}},
			{{Text: backButtonText, CallbackData: rootMenuCallback}},
		},
	}

	iterKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "100", CallbackData: fmt.Sprintf("%s/100", iterButtonCallback)},
				{Text: "200", CallbackData: fmt.Sprintf("%s/200", iterButtonCallback)},
				{Text: "400", CallbackData: fmt.Sprintf("%s/400", iterButtonCallback)},
			},
			{
				{Text: "800", CallbackData: fmt.Sprintf("%s/800", iterButtonCallback)},
				{Text: "1000", CallbackData: fmt.Sprintf("%s/1000", iterButtonCallback)},
				{Text: "2000", CallbackData: fmt.Sprintf("%s/2000", iterButtonCallback)},
			},
			{
				{Text: otherButtonText, CallbackData: iterInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: rootMenuCallback},
			},
		},
	}

	repKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "1", CallbackData: fmt.Sprintf("%s/1", repButtonCallback)},
				{Text: "2", CallbackData: fmt.Sprintf("%s/2", repButtonCallback)},
				{Text: "3", CallbackData: fmt.Sprintf("%s/3", repButtonCallback)},
			},
			{
				{Text: "4", CallbackData: fmt.Sprintf("%s/4", repButtonCallback)},
				{Text: "5", CallbackData: fmt.Sprintf("%s/5", repButtonCallback)},
				{Text: "6", CallbackData: fmt.Sprintf("%s/6", repButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: rootMenuCallback},
			},
		},
	}

	alphaKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: autoButtonText, CallbackData: fmt.Sprintf("%s/0", alphaButtonCallback)},
			},
			{
				{Text: "32", CallbackData: fmt.Sprintf("%s/32", alphaButtonCallback)},
				{Text: "64", CallbackData: fmt.Sprintf("%s/64", alphaButtonCallback)},
				{Text: "128", CallbackData: fmt.Sprintf("%s/128", alphaButtonCallback)},
				{Text: "255", CallbackData: fmt.Sprintf("%s/255", alphaButtonCallback)},
			},
			{
				{Text: otherButtonText, CallbackData: alphaInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: rootMenuCallback},
			},
		},
	}

	extKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "jpg", CallbackData: fmt.Sprintf("%s/jpg", extButtonCallback)},
				{Text: "png", CallbackData: fmt.Sprintf("%s/png", extButtonCallback)},
				{Text: "svg", CallbackData: fmt.Sprintf("%s/svg", extButtonCallback)},
				// gifs are disabled due to performance issues
				// {Text: "gif", CallbackData: fmt.Sprintf("%s/gif", extButtonCallback)},
			},
			{
				{Text: backButtonText, CallbackData: rootMenuCallback},
			},
		},
	}

	sizeKeyboard = telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "256", CallbackData: fmt.Sprintf("%s/256", sizeButtonCallback)},
				{Text: "512", CallbackData: fmt.Sprintf("%s/512", sizeButtonCallback)},
				{Text: "720", CallbackData: fmt.Sprintf("%s/720", sizeButtonCallback)},
			},
			{
				{Text: "1024", CallbackData: fmt.Sprintf("%s/1024", sizeButtonCallback)},
				{Text: "1280", CallbackData: fmt.Sprintf("%s/1280", sizeButtonCallback)},
				{Text: "1920", CallbackData: fmt.Sprintf("%s/1920", sizeButtonCallback)},
			},
			{
				{Text: otherButtonText, CallbackData: sizeInputCallback},
			},
			{
				{Text: backButtonText, CallbackData: rootMenuCallback},
			},
		},
	}
)

// newKeyboardFromTemplate creates new keyboard from the template
// adding symbol to the option that is chosen at the moment
func newKeyboardFromTemplate(
	template telegram.InlineKeyboardMarkup,
	optionCallback,
	newText string,
) telegram.InlineKeyboardMarkup {
	checkSymbol := "👉"

	newKeyboard := telegram.InlineKeyboardMarkup{}
	newKeyboard.InlineKeyboard = make([][]telegram.InlineKeyboardButton, len(template.InlineKeyboard))
	for i, row := range template.InlineKeyboard {
		newKeyboard.InlineKeyboard[i] = make([]telegram.InlineKeyboardButton, len(row))
		for j, button := range row {
			newKeyboard.InlineKeyboard[i][j] = button

			if button.CallbackData == optionCallback {
				if newText == "" {
					newText = button.Text
				}

				newKeyboard.InlineKeyboard[i][j].Text = fmt.Sprintf("%s %s", checkSymbol, newText)
			}
		}
	}

	return newKeyboard
}
