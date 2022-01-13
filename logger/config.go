package logger

// Config 日志配置
type Config struct {
	name       string // 项目名称
	level      Level
	envMode    string
	skip       int  // CallerSkip次数
	jsonFormat bool // 是否输出 json
	consoleOut bool // 是否输出到控制台
	colorful   bool
	fileOut    *fileOutConfig
}

// 日志输出文件配置
type fileOutConfig struct {
	filename     string // 日志文件名称
	path         string // 文件保存路径
	maxAge       uint   //  保存几天的日志(天)
	rotationTime uint   // 日志切割时间间隔(小时)
}

func defaultConfig() *Config {
	return &Config{
		name:       "zap-logger",
		level:      InfoLevel,
		envMode:    "dev",
		skip:       2,
		jsonFormat: false,
		consoleOut: true,
		fileOut:    nil,
	}
}

func (c *Config) SetLevel(level Level) {
	c.level = level
}

func (c *Config) SetName(name string) {
	c.name = name
}

func (c *Config) SetEnvMode(envMode string) {
	c.envMode = envMode
}

func (c *Config) SetSkip(skip int) {
	c.skip = skip
}

func (c *Config) SetJsonFormat(jsonFormat bool) {
	c.jsonFormat = jsonFormat
}

func (c *Config) SetConsoleOut(consoleOut bool) {
	c.consoleOut = consoleOut
}

func (c *Config) SetColorful(colorful bool) {
	c.colorful = colorful
}

func (c *Config) EnableFileOut() {
	c.fileOut = &fileOutConfig{
		filename:     c.name,
		path:         "",
		maxAge:       7,
		rotationTime: 24,
	}
}

func (c *Config) DisableFileOut() {
	c.fileOut = nil
}

// 设置 文件输出
func (c *Config) SetFileOut(path string, maxAge, rotationTime uint) {
	c.fileOut = &fileOutConfig{
		filename:     c.name,
		path:         path,
		maxAge:       maxAge,
		rotationTime: rotationTime,
	}
}
