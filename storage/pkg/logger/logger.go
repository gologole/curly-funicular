package logger

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
)

/*Структурированные логи: Для каждого лога добавляются поля:

    module — указывает на компонент системы (например, api, db, processor).
    action — действие, которое выполняется в момент логирования (например, fetch_data, connect,
process_data).
    request_id — уникальный идентификатор запроса.
    user_id — уникальный идентификатор пользователя.

Дополнительные поля: К примеру, в logger.Warn добавлено поле retry_after,
указывающее время повторной попытки.

Заполнение полей: Эти поля позволяют создавать более структурированные и детализированные логи,
которые удобнее анализировать через Kibana.*/

// Функция для отправки логов в Logstash
func sendLogToLogstash(logstashURL string, logEntry []byte) error {
	req, err := http.NewRequest("POST", os.Getenv("ELK_DOMAIN"), // logstash:5044
		bytes.NewBuffer(logEntry))
	if err != nil {
		log.Printf("НЕТ СВЯЗИ С ELK")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

// Создание кастомного core для отправки логов в Logstash
func logstashCore(logstashURL string) zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp" // Форматирование времени для ELK
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logstashWriteSyncer := zapcore.AddSync(&logstashWriter{url: logstashURL})
	return zapcore.NewCore(encoder, logstashWriteSyncer, zapcore.DebugLevel)
}

// Реализация zapcore.WriteSyncer для отправки в Logstash
type logstashWriter struct {
	url string
}

func (w *logstashWriter) Write(p []byte) (n int, err error) {
	err = sendLogToLogstash(w.url, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *logstashWriter) Sync() error {
	return nil
}

// Функция для настройки zap логгера
func ConfigureLogger(logstashURL string) *zap.Logger {
	// Создаем конфигурацию для вывода в консоль (Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriteSyncer, zap.WarnLevel)

	// Создаем core для отправки всех логов в Logstash
	logstashCore := logstashCore(logstashURL)

	// Комбинируем оба core
	combinedCore := zapcore.NewTee(consoleCore, logstashCore)

	return zap.New(combinedCore)
}
