package service

import (
	"context"
	"fmt"
	gas "github.com/firstrow/goautosocket"
	logrustash "github.com/idtksrv/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level uint32

const (
	LevelFatal = iota //exit after logging
	LevelPanic        //throw panic after logging
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
	FormatterTypeText
	FormatterTypeJSON
)

type LoggerParameter struct {
	ELKEnable   bool
	LogFileName string
	Level       string
	Address     string
	Environment string
	Service     string
}

type Logger struct {
	level          Level
	logFileName    string
	activeLog      bool
	activeLogFile  bool
	fullPath       bool
	logrus         *logrus.Logger
	forceColor     bool //當輸出到stdout, 是否使用彩色
	formatterType  int
	environment    string
	service        string
	logstashEnable bool
	logstashReady  bool
	logstashEntry  *logrus.Entry
	logstashChan   chan *logstashData
	logstashCtx    context.Context
	logstashCancel context.CancelFunc
}

type logstashData struct {
	level  Level
	fields logrus.Fields
	args   []interface{}
}

//func newLogger(elkEnable bool, fileNamePrefix, level, address, env, service string) (*Logger, error) {
func newLogger(param *LoggerParameter) (*Logger, error) {

	if param.LogFileName == "" {
		return nil, fmt.Errorf("logger file name is empty")
	}

	l := &Logger{
		level:          0,
		logFileName:    param.LogFileName,
		activeLog:      true,
		activeLogFile:  true,
		fullPath:       false,
		logrus:         logrus.New(),
		environment:    param.Environment,
		service:        param.Service,
		logstashEnable: param.ELKEnable,
		logstashReady:  false,
	}

	l.level = l.GetLevel(param.Level)
	l.SetFormatter(FormatterTypeText, true, false, nil)
	l.logrus.SetReportCaller(false) //logrus內部的列印呼叫的func和檔案、行數，沒用

	if param.ELKEnable {

		l.logstashEntry = l.NewEntry()
		l.logstashCtx, l.logstashCancel = context.WithCancel(context.Background())
		l.logstashChan = make(chan *logstashData)

		conn, err := gas.Dial("tcp", param.Address)
		if err != nil {
			l.logrus.Error("newLogger failed at gas.Dail", err)
		} else {
			l.logrus.Hooks.Add(logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{})))
			l.logstashReady = true
		}

		go l.logstashWorker(l.logstashCtx, l.logstashChan)
	}

	return l, nil
}

func (l *Logger) NewEntry() *logrus.Entry {
	ent := logrus.NewEntry(l.logrus).WithFields(logrus.Fields{
		"requestID": uuidGetV4(),
		"env":       l.environment,
		"service":   l.service,
	})

	return ent
}

//Out sets log output
//disable stdout, pass in ioutil.Discard
//l.logrus.Out = ioutil.Discard
//enable stdout, pass in os.Stdout
//l.logrus.Out=os.Stdout
func (l *Logger) Out(out io.Writer) {
	l.logrus.Out = out
}

func (l *Logger) WithFieldsHook(entry *logrus.Entry, level Level, fields logrus.Fields, args ...interface{}) {

	//if logstash is ready
	if l.logstashReady {
		//push into logstash queue
		l.logstashChan <-
			&logstashData{
				level,
				fields,
				args,
			}
		return
	}

	//logstash is not ready, just log
	//if level == LevelFatal {
	//	entry.WithFields(fields).Fatal(args...)
	//} else if level == LevelPanic {
	//	entry.WithFields(fields).Panic(args...)
	//} else if level == LevelError {
	//	entry.WithFields(fields).Error(args...)
	//} else if level == LevelWarning {
	//	entry.WithFields(fields).Warning(args...)
	//} else if level == LevelInfo {
	//	entry.WithFields(fields).Info(args...)
	//} else if level == LevelDebug {
	//	entry.WithFields(fields).Debug(args...)
	//} else {
	//	entry.WithFields(fields).Trace(args...)
	//}
}

func (l *Logger) logstashWorker(ctx context.Context, dataChan <-chan *logstashData) {
	for {
		select {
		case <-ctx.Done():
			//接收到取消訊號
			//		logger.Println("任務", ctx.Value(key), ":任務停止...")
			return
		case data := <-dataChan:
			if data.level == LevelFatal {
				l.logstashEntry.WithFields(data.fields).Fatal(data.args...)
			} else if data.level == LevelPanic {
				l.logstashEntry.WithFields(data.fields).Panic(data.args...)
			} else if data.level == LevelError {
				l.logstashEntry.WithFields(data.fields).Error(data.args...)
			} else if data.level == LevelWarning {
				l.logstashEntry.WithFields(data.fields).Warning(data.args...)
			} else if data.level == LevelInfo {
				l.logstashEntry.WithFields(data.fields).Info(data.args...)
			} else if data.level == LevelDebug {
				l.logstashEntry.WithFields(data.fields).Debug(data.args...)
			} else {
				l.logstashEntry.WithFields(data.fields).Trace(data.args...)
			}
		}
	}
}

func (l *Logger) GetLevel(level string) Level {
	if level == "debug" {
		return LevelDebug
	}
	if level == "info" {
		return LevelInfo
	}
	if level == "warning" {
		return LevelWarning
	}
	if level == "error" {
		return LevelError
	}
	if level == "panic" {
		return LevelPanic
	}
	if level == "fatal" {
		return LevelFatal
	}
	if level == "trace" {
		return LevelTrace
	}
	return LevelInfo
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) IsFullPath(fullPath bool) {
	l.fullPath = fullPath
}

//SetLogFileName sets file name for log file.
func (l *Logger) SetLogFileName(name string) {
	l.logFileName = name

}

//Fields is a temp type for logrus.Fields
type Fields map[string]interface{}

//SetFormatter can set logger format.
//reportCaller
//callerReporter
//if callerReporter is nil and reportCaller is true, will using default call reporter
//Parameter json for json format.
//Parameter text for text format.
func (l *Logger) SetFormatter(formatterType int, forceColor bool, reportCaller bool, callerReporter func(f *runtime.Frame) (string, string)) {
	l.forceColor = forceColor
	l.formatterType = formatterType
	l.logrus.SetReportCaller(reportCaller)

	if formatterType == FormatterTypeText {
		l.logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			ForceColors:      forceColor,
			DisableColors:    false,
			PadLevelText:     false,
			CallerPrettyfier: callerReporter,
		})
	} else {
		l.logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint:      true,
			CallerPrettyfier: callerReporter,
		})
	}
}

//CallerParser parse func name,file name, line and return func for formatter callerPrettyfier
func (l *Logger) CallerParser(funcName, fileName string, line int) func(*runtime.Frame) (string, string) {

	return func(f *runtime.Frame) (string, string) {
		return funcName, fmt.Sprintf(" %s:%d", fileName, line)
	}
}

//ReturnCaller return caller name,file name line
func (l *Logger) ReturnCaller() (funcName string, fileName string, line int) {

	pc := make([]uintptr, 5)
	n := runtime.Callers(2, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return "", "", 0
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	f, _ := frames.Next()

	filename := f.File
	if !l.fullPath {
		filename = f.File[strings.LastIndex(f.File, "/")+1:]
	}

	fn := f.Function[strings.LastIndex(f.Function, ".")+1:]
	return fn, filename, f.Line
}

func (l *Logger) ReturnCallerFormatted() string {

	pc := make([]uintptr, 5)
	n := runtime.Callers(2, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	f, _ := frames.Next()

	filename := f.File
	if !l.fullPath {
		filename = f.File[strings.LastIndex(f.File, "/")+1:]
	}

	fn := f.Function[strings.LastIndex(f.Function, ".")+1:]

	return fmt.Sprintf("%s %s:%d", fn, filename, f.Line)
}

//spaceFieldsJoin stripping all whitespace characters.
func (l *Logger) spaceFieldsJoin(str string) string {
	return strings.Join(strings.Fields(str), "")
}

func (l *Logger) log(targetLevel Level, entry *logrus.Entry, args ...interface{}) {
	if targetLevel > l.level {
		return
	}

	if targetLevel == LevelFatal {
		entry.Fatal(args...)
	} else if targetLevel == LevelPanic {
		entry.Panic(args...)
	} else if targetLevel == LevelError {
		entry.Error(args...)
	} else if targetLevel == LevelWarning {
		entry.Warning(args...)
	} else if targetLevel == LevelInfo {
		entry.Info(args...)
	} else if targetLevel == LevelDebug {
		entry.Debug(args...)
	} else {
		entry.Trace(args...)
	}
}

//Log logs with level and args
func (l *Logger) Log(targetLevel Level, args ...interface{}) {
	if targetLevel == LevelFatal {
		l.logrus.Fatal(args...)
	} else if targetLevel == LevelPanic {
		l.logrus.Panic(args...)
	} else if targetLevel == LevelError {
		l.logrus.Error(args...)
	} else if targetLevel == LevelWarning {
		l.logrus.Warning(args...)
	} else if targetLevel == LevelInfo {
		l.logrus.Info(args...)
	} else if targetLevel == LevelDebug {
		l.logrus.Debug(args...)
	} else {
		l.logrus.Trace(args...)
	}
}

func (l *Logger) LogFile(targetLevel Level, args ...interface{}) {
	filePath := l.logFileName + time.Now().Format("20060102") + ".log"
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("LogFile %s", err.Error())
		return
	}

	defer f.Close()

	//關掉 color 如果是 TextFormatter
	if l.formatterType == FormatterTypeText {
		l.logrus.Formatter.(*logrus.TextFormatter).DisableColors = true
	}

	l.logrus.SetOutput(f)
	l.Log(targetLevel, args...)
}

func (l *Logger) WithField(level Level, key string, value interface{}, args ...interface{}) {

	var entry *logrus.Entry
	l.logrus.SetOutput(os.Stdout)

	if len(l.spaceFieldsJoin(key)) > 0 {
		entry = l.logrus.WithField(key, value)
	} else {
		//no WithField ,create entry
		entry = logrus.NewEntry(l.logrus)
	}
	l.log(level, entry, args...)

}

func (l *Logger) WithFields(level Level, fields Fields, args ...interface{}) {
	var entry *logrus.Entry
	l.logrus.SetOutput(os.Stdout)

	if len(fields) > 0 {

		//local Fields to logrus.Fields
		f := logrus.Fields{}
		for k, v := range fields {
			f[k] = v
		}

		entry = l.logrus.WithFields(f)
	} else {
		entry = logrus.NewEntry(l.logrus)
	}
	l.log(level, entry, args...)
}

func (l *Logger) WithFieldFile(level Level, key string, value interface{}, args ...interface{}) error {

	filePath := l.logFileName + time.Now().Format("20060102") + ".log"
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("WithFieldFile %s", err.Error())
		return err
	}

	defer f.Close()

	var entry *logrus.Entry
	//關掉 color 如果是 TextFormatter
	if l.formatterType == FormatterTypeText {
		l.logrus.Formatter.(*logrus.TextFormatter).DisableColors = true
	}

	l.logrus.SetOutput(f)

	if len(l.spaceFieldsJoin(key)) > 0 {
		entry = l.logrus.WithField(key, value)
	} else {
		entry = logrus.NewEntry(l.logrus)
	}

	l.log(level, entry, args...)
	return nil
}

func (l *Logger) WithFieldsFile(level Level, fields Fields, args ...interface{}) error {
	filePath := l.logFileName + time.Now().Format("20060102") + ".log"
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("WithFieldsFile %s", err.Error())
		return err
	}

	defer f.Close()

	var entry *logrus.Entry

	//輸出到檔案，關閉 force color
	//關掉 color 如果是 TextFormatter
	if l.formatterType == FormatterTypeText {
		l.logrus.Formatter.(*logrus.TextFormatter).DisableColors = true
	}

	l.logrus.SetOutput(f)

	if len(fields) > 0 {

		fd := logrus.Fields{}
		for k, v := range fields {
			fd[k] = v
		}
		entry = l.logrus.WithFields(fd)
	} else {
		entry = logrus.NewEntry(l.logrus)
	}

	l.log(level, entry, args...)
	return nil

}
