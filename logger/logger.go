package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Config   *Config
	zapSugar *zap.SugaredLogger
}

func New() *Logger {
	logger := &Logger{
		Config: initConfig(),
	}

	// 根据环境设置 输出格式
	if strings.EqualFold(logger.Config.envMode, "prod") {
		logger.Config.jsonFormat = true
		logger.Config.consoleOut = false
	}

	// 设置 logger
	logger.WithConfig()

	return logger
}

func (l *Logger) getFileWriter() (infoWriter io.Writer, warnWriter io.Writer) {
	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	loggerPath := l.Config.fileOut.path
	appName := l.Config.name

	if loggerPath == "" || strings.EqualFold(l.Config.envMode, "dev") {
		address, _ := os.Getwd()
		loggerPath = address + "/logs"
	}

	// 检查文件夹是否存在
	isExist := CheckFileAndCreate(loggerPath)

	if isExist != nil {
		panic("[文件夹不存在] " + loggerPath)
	}

	warnWriterFileName := loggerPath + "/common-error.log"
	if len(appName) > 0 {
		warnWriterFileName = loggerPath + "/" + appName + "-common-error.log"
	}

	infoWriter = getWriter(loggerPath+"/"+appName+".log",
		l.Config.fileOut.rotationTime,
		l.Config.fileOut.maxAge)

	warnWriter = getWriter(warnWriterFileName,
		l.Config.fileOut.rotationTime,
		l.Config.fileOut.maxAge)

	return
}

func (l *Logger) generateEncoderConfig() zapcore.EncoderConfig {
	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	encodeCaller := zapcore.FullCallerEncoder
	if !strings.EqualFold(l.Config.envMode, "prod") {
		encodeCaller = zapcore.ShortCallerEncoder // 只显示 package/file.go:line
	}

	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	return zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		CallerKey:    "file",
		EncodeTime:   customTimeEncoder,
		EncodeLevel:  customLevelEncoder,
		EncodeCaller: encodeCaller,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		EncodeName: zapcore.FullNameEncoder,
	}
}

// 根据配置文件更新 logger
func (l *Logger) WithConfig() {
	// 生成 encoderConfig
	encoderConfig := l.generateEncoderConfig()

	var encoder zapcore.Encoder
	if !l.Config.jsonFormat {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 设置级别
	logLevel := convertLevel(l.Config.level)

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	var infoWs []zapcore.WriteSyncer
	var warnWs []zapcore.WriteSyncer

	// 是否开启输出文件
	if l.Config.fileOut != nil {
		info, warn := enableFileOut(l.getFileWriter())
		infoWs = append(infoWs, info)
		warnWs = append(warnWs, warn)
	}

	// 控制台输出
	if l.Config.consoleOut {
		info, warn := enableConsoleOut()
		infoWs = append(infoWs, info)
		warnWs = append(warnWs, warn)
	}

	// 不能同时关闭两种输出方式，否则将开启控制台输出
	if l.Config.fileOut == nil && !l.Config.consoleOut {
		l.Warn("You cannot turn off both output modes at the same time!!!")
		info, warn := enableConsoleOut()
		infoWs = append(infoWs, info)
		warnWs = append(warnWs, warn)
	}

	wsInfo := zapcore.NewMultiWriteSyncer(infoWs...)
	wsWarn := zapcore.NewMultiWriteSyncer(warnWs...)
	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, wsInfo, infoLevel),
		zapcore.NewCore(encoder, wsWarn, warnLevel),
	)

	// 开启开发模式，堆栈跟踪 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	caller := zap.AddCaller()
	// 防止zap始终将包装器代码报告为调用者( 需要跳过一个级别，否则打印的文件名和行号 是封装的文件名)
	skip := zap.AddCallerSkip(l.Config.skip)

	zapLog := zap.New(core, caller, skip)
	defer zapLog.Sync()

	l.zapSugar = zapLog.Sugar().With()
}

// With adds a variadic number of fields to the logging context.
// see https://github.com/uber-go/zap/blob/v1.10.0/sugar.go#L91
func (l *Logger) With(args ...interface{}) *Logger {
	l.zapSugar = l.zapSugar.With(args...)
	return l
}

// Debug package sugar of zap
func (l *Logger) Debug(args ...interface{}) {
	l.zapSugar.Debug(args...)
}

// Debugf package sugar of zap
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.zapSugar.Debugf(template, args...)
}

// Info package sugar of zap
func (l *Logger) Info(args ...interface{}) {
	l.zapSugar.Info(args...)
}

// Infof package sugar of zap
func (l *Logger) Infof(template string, args ...interface{}) {
	l.zapSugar.Infof(template, args...)
}

// Warn package sugar of zap
func (l *Logger) Warn(args ...interface{}) {
	l.zapSugar.Warn(args...)
}

// Warnf package sugar of zap
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.zapSugar.Warnf(template, args...)
}

// Error package sugar of zap
func (l *Logger) Error(args ...interface{}) {
	l.zapSugar.Error(args...)
}

// Errorf package sugar of zap
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.zapSugar.Errorf(template, args...)
}

// Fatal package sugar of zap
func (l *Logger) Fatal(args ...interface{}) {
	l.zapSugar.Fatal(args...)
}

// Fatalf package sugar of zap
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.zapSugar.Fatalf(template, args...)
}

// Panic package sugar of zap
func (l *Logger) Panic(args ...interface{}) {
	l.zapSugar.Panic(args...)
}

// Panicf package sugar of zap
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.zapSugar.Panicf(template, args...)
}
