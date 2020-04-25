package logger

// конфигурация для подключения записи в файл
type lamberJackConfig struct {
	filename   string
	maxSize    int
	maxAge     int
	maxBackups int
	localTime  bool
	compress   bool
}
