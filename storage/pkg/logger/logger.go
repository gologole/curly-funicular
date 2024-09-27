package pkg

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

// Функция для отправки логов в Logstash
func sendLogToLogstash(logstashURL string, logEntry []byte) error {
	req, err := http.NewRequest("POST", logstashURL, bytes.NewBuffer(logEntry))
	if err != nil {
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
func configureLogger(logstashURL string) *zap.Logger {
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
