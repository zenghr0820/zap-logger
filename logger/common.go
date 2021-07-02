package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
	FatalLevel Level = "FATAL"

	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

type Level string

func (logLevel Level) String() string {
	return strings.ToUpper(string(logLevel))
}

// 把字符串转换为日志级别（数字）
func convertLevel(logLevel Level) zapcore.Level {
	// 不区分大小写
	var level zapcore.Level
	switch logLevel {
	case DebugLevel:
		level = zap.DebugLevel
	case InfoLevel:
		level = zap.InfoLevel
	case WarnLevel:
		level = zap.WarnLevel
	case ErrorLevel:
		level = zap.ErrorLevel
	case FatalLevel:
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

func getWriter(filename string, maxAge, rotationTime uint) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每天分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y-%m-%d", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(rotationTime)),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func enableConsoleOut() (info, warn zapcore.WriteSyncer) {
	writeInfoConsole := zapcore.AddSync(os.Stdout)
	writeWarnConsole := zapcore.AddSync(os.Stderr)

	return writeInfoConsole, writeWarnConsole
}

func enableFileOut(infoWriter, warnWriter io.Writer) (info, warn zapcore.WriteSyncer) {
	writeInfoFile := zapcore.AddSync(infoWriter)
	writeWarnFile := zapcore.AddSync(warnWriter)

	return writeInfoFile, writeWarnFile
}

// 检查文件是否存在并且创建文件
func CheckFileAndCreate(path string) (err error) {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		fmt.Println("path not exists ", path)
		err := os.MkdirAll(path, 0711)

		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return err
		}
		return nil
	}
}
