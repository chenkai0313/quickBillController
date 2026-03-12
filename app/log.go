package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"quickBillController/config"
	"quickBillController/utils/color"
)

var ZapLog *zap.Logger
var nowLogDir = ""

func InitLogger() {
	timeStr := time.Now().Format("2006-01-02")
	nowLogDir = timeStr
	logPath := "./logs/" + timeStr + "/"
	logFileName := fmt.Sprintf("%v/%v.log", logPath, "server")

	createLogFile(logPath, logFileName)
	ZapLog = initLog(logFileName)

	go CronDeleteLogOlderFile("./logs/", config.GetCfg().Log.MaxAge)

}

func initLog(fileName string) *zap.Logger {
	lumberJackLogger := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.GetCfg().Log.MaxSize,
		MaxBackups: config.GetCfg().Log.MaxBackups,
		MaxAge:     config.GetCfg().Log.MaxAge,
		Compress:   true,
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.MessageKey = "category"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	atomicLevel := zap.NewAtomicLevel()

	logLevel := zap.DebugLevel

	switch config.GetCfg().Db.LogLevel {
	case "debug":
		logLevel = zap.DebugLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "info":
		logLevel = zap.InfoLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.DebugLevel
	}
	atomicLevel.SetLevel(logLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(Writer{}), zapcore.AddSync(&lumberJackLogger)),
		atomicLevel,
	)

	return zap.New(core, zap.AddCaller())
}

func createLogFile(logPath string, logName string) {
	existBool, _ := isFileExist(logName)
	if !existBool {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			panic("Log folder creation failed")
		}
		f, err := os.Create(logName)
		if err != nil {
			panic("Log file creation failed")
		}
		f.Close()
	}
}

func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

type Writer struct {
}

func (w Writer) Write(p []byte) (n int, err error) {
	defer func() {
		if time.Now().Format("2006-01-02") != nowLogDir {
			InitLogger()
		}
	}()
	log := make(map[string]interface{})
	_ = json.Unmarshal(p, &log)

	level, _ := log["level"]
	switch level {
	case "error":
		fmt.Println(color.Red(string(p)))
	case "panic":
		fmt.Println(color.Red(string(p)))
	case "warn":
		fmt.Println(color.Yellow(string(p)))
	case "fatal":
		fmt.Println(color.Red(string(p)))
	default:
		fmt.Println(color.Green(string(p)))
	}
	return
}

func CronDeleteLogOlderFile(directoryPath string, maxAge int) {
	if maxAge == 0 {
		maxAge = 7
	}

	clearFun := func(directoryPath string, maxAge int) {
		sevenDaysAgo := time.Now().AddDate(0, 0, -maxAge)
		layout := "2006-01-02"

		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

		_ = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				dirName := info.Name()

				match := re.FindString(dirName)
				if match != "" {
					date, err := time.Parse(layout, match)
					if err == nil && date.Before(sevenDaysAgo) {
						_ = os.RemoveAll(path)
					}
				}
			}

			return nil
		})
	}

	go clearFun(directoryPath, maxAge)
	ticker := time.NewTicker(3600 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			clearFun(directoryPath, maxAge)
		}
	}
}
