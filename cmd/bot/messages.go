package main

const (
	helpMessage            = "Отправь мне какую-нибудь фотографию."
	errorMessage           = "Что-то пошло не так! Попробуй снова через пару минут."
	inputMessage           = "Неверное значение!\nВведи число от %d до %d:"
	statusMessage          = "%d место в очереди.\n\nФигуры: %s\nИтерации: %d\nПовторения: %d\nАльфа-канал: %d\nРасширение: %s\nРазмеры: %d"
	statusEmptyMessage     = "Нету активных операций."
	operationsLimitMessage = "Вы не можете добавить больше операций в очередь."
)

const (
	enqueuedLogMessage = "Enqueued: user id %d | input %s | iterations=%d, shape=%d, alpha=%d, repeat=%d, resolution=%d, extension=%s"
	creatingLogMessage = "Creating: user id %d | input %s | output %s | iterations=%d, shape=%d, alpha=%d, repeat=%d, resolution=%d, extension=%s"
	finishedLogMessage = "Finished: user id %d | input %s | output %s | %.1f seconds"
	sentLogMessage     = "Sent: user id %d | output %s"
)
