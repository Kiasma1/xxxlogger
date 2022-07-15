package xxxlogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

var logger *zap.Logger

// InitProdLogger 只写入 文件 中
// 且不输出 Debug日志
func InitProdLogger(logFilepath string){
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置 日志行 中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置 等级 大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 日志Encoder（使用JSONEncode，把 日志行 格式化为JSON格式）
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, getFileLogWriter(logFilepath), zapcore.InfoLevel),
	)

	logger = zap.New(core)
}

// InitDevLogger 同时写入 控制台 和 文件 中
// 输出 Debug日志
func InitDevLogger(logFilepath string){
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置 日志行 中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置 等级 大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 日志Encoder（使用JSONEncode，把 日志行 格式化为JSON格式）
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, getFileLogWriter(logFilepath), zapcore.DebugLevel),
	)

	logger = zap.New(core)
}

func getFileLogWriter(logFilepath string)(writeSyncer zapcore.WriteSyncer){
	lumberJackLogger := &lumberjack.Logger{
		Filename: logFilepath,
		MaxSize: 100, // 单个日志文件最大100MB
		MaxBackups: 60, // 多于60个日志文件后，清理较旧的日志
		MaxAge: 1, // 一天一切割
		Compress: false, // 不实用 gzip 压缩日志文件
	}

	return zapcore.AddSync(lumberJackLogger)
}

func Debug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Error(message, fields...)
}


func getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}
