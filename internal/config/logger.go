package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes a production-ready Zap logger with log rotation.
func InitLogger(env string) *zap.Logger {
	// Setup Lumberjack untuk log rotation
	logWriter := &lumberjack.Logger{
		Filename:   "logs/app.log", // Nama file output (otomatis bikin folder kalau belum ada)
		MaxSize:    100,            // Maksimal 100 MB per file sebelum di-rotate
		MaxBackups: 30,             // Simpan maksimal 30 file log lama
		MaxAge:     28,             // Hapus log yang lebih tua dari 28 hari
		Compress:   true,           // Zip file log yang lama biar hemat storage
	}

	// Tentukan format log
	var encoder zapcore.Encoder
	if env == "production" {
		// Format JSON murni untuk production (gampang dibaca mesin/Elasticsearch)
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// Format Console warna-warni untuk development di terminal
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Output ganda: Tulis ke File (Lumberjack) DAN ke Terminal (Stdout)
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logWriter)),
		zap.DebugLevel,
	)

	return zap.New(core, zap.AddCaller())
}
