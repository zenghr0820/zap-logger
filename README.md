# zap-logger
 基于 zap 库的日志库，根据时间分割日志文件



# 功能

- [x] 根据info/warn级别切割日志文件
- [x] 根据时间切割日志文件
- [x] 自动格式化 format
- [x] 根据运行环境，切换输出格式，是否输出控制台



# 使用

**```go get -u github.com/zenghr0820/zap-logger```**

```go
import (
	logger "github.com/zenghr0820/zap-logger"
)

func main() {
	logger := zapLogger.InitLog(&zapLogger.Config{
		Name:    "demo",
		Dir:     "",
		Level:   zapLogger.InfoLevel,
		EnvMode: "dev",
	})
	logger.Info("Info...", 1)
	logger.Warn("Warn...", 2)
	logger.Error("Error...", 3)
	logger.Debug("Debug...", 4)
}
```



# 输出

#### 一、 dev 环境输出
```bash
2020-04-11 17:35:37	INFO	Ex/main.go:14	Info...1
2020-04-11 17:35:37	WARN	Ex/main.go:15	Warn...2
2020-04-11 17:35:37	ERROR	Ex/main.go:16	Error...3
2020-04-11 17:35:37	DEBUG	Ex/main.go:17	Debug...4
```



#### 二、格式化输出

```bash
{"level":"INFO","ts":"2020-04-11 17:40:43","file":"Ex/main.go:14","msg":"Info...1"}
{"level":"DEBUG","ts":"2020-04-11 17:40:43","file":"Ex/main.go:17","msg":"Debug...4"}
```



#### 三、生成文件

```bash
common-error.log
common-error.log.2020-04-11
demo.log
demo.log.2020-04-11
```

