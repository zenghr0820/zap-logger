# zap-logger
 基于 zap 库的日志库，根据时间分割日志文件



## 功能

- [x] 根据info/warn级别切割日志文件
- [x] 根据时间切割日志文件
- [x] 自动格式化 format
- [x] 根据运行环境，切换输出格式，是否输出控制台
- [x] 提供默认的日志记录器和生成另一个日志记录器



## 使用默认的日志记录器

**```go get -u github.com/zenghr0820/zap-logger```**

```go
package main

import (
	"github.com/zenghr0820/zap-logger"
)

func main() {
	zapLogger.Info("Info...", 1)
	zapLogger.Infof("Infof -> %d", 1)
	zapLogger.Warn("Warn...", 2)
	zapLogger.Error("Error...", 3)
	zapLogger.Debug("Debug...", 4)
	// 开启 JSON 格式化输出
	zapLogger.Config.SetJsonFormat(true)
	// 更新配置
	zapLogger.WithConfig()
	zapLogger.Info("Info...", 5)
	zapLogger.Infof("Infof -> %d", 5)
	zapLogger.Warn("Warn...", 6)
	zapLogger.Error("Error...", 7)
	zapLogger.Debug("Debug...", 8)
}

```

**控制台输出**

```
[2021-07-02 15:13:29.904]	[INFO]	logger/main.go:8	Info...1
[2021-07-02 15:13:29.905]	[INFO]	logger/main.go:9	Infof -> 1
{"level":"[INFO]","ts":"[2021-07-02 15:13:29.905]","file":"logger/main.go:17","msg":"Info...5"}
{"level":"[INFO]","ts":"[2021-07-02 15:13:29.905]","file":"logger/main.go:18","msg":"Infof -> 5"}
[2021-07-02 15:13:29.905]	[WARN]	logger/main.go:10	Warn...2
[2021-07-02 15:13:29.905]	[ERROR]	logger/main.go:11	Error...3
{"level":"[WARN]","ts":"[2021-07-02 15:13:29.905]","file":"logger/main.go:19","msg":"Warn...6"}
{"level":"[ERROR]","ts":"[2021-07-02 15:13:29.905]","file":"logger/main.go:20","msg":"Error...7"}
```

**默认的日志记录器只开启控制台输出，如需要开启文件输出，可以开启默认文件输出或者自定义文件输出配置**

```go
// 开启文件输出
zapLogger.Config.EnableFileOut()
// 更新配置
zapLogger.WithConfig()
```

**自定义文件输出**

```go
// 开启文件输出
zapLogger.Config.SetFileOut("/logs", 7, 24)
// 更新配置
zapLogger.WithConfig()
```

参数说明：

- path：日志文件输出保存路径，默认保存至 `/logs`
- maxAge：保存 maxAge 天内的日志，默认保存 7 天内
- rotationTime：每 rotationTime 分割一次日志，默认 1 天分割一次

## 新建日志记录器

```go
package main

import (
	"github.com/zenghr0820/zap-logger/logger"
)

func main() {
	l := logger.New(nil)
	l.Info("Info...", 1)
	l.Infof("Infof -> %d", 1)
	l.Warn("Warn...", 2)
	l.Error("Error...", 3)
	l.Debug("Debug...", 4)
}
```

