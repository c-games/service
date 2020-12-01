package service

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_newLogger(t *testing.T) {

	p := &LoggerParameter{
		false,
		"log",
		"error",
		"113.35.204.62:5000",
		"local",
		"srv",
	}

	l, _ := newLogger(p)

	entry := l.NewEntry()
	//lv,_:=logrus.ParseLevel("info")
	entry.WithField("key", "value").Log(logrus.WarnLevel, "data")

	//fmt.Printf("defalut callReprter: \n")
	l.SetFormatter(FormatterTypeText, true, true, nil)

	l.WithField(LevelInfo, "text", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "text", 1, "a", "b")

	//fmt.Printf("custom callReprter: \n")

	//方法一，設定成formatter
	//但行數會錯，除非每次都設定
	l.SetFormatter(FormatterTypeText, true, true, l.CallerParser(l.ReturnCaller()))
	l.WithField(LevelInfo, "json", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "json", 1, "a", "b")
	_ = l.WithFieldsFile(LevelInfo, Fields{"1": 1, "2": 3}, "a", "b")

	//方法二，關閉reportCaller, 自己呼叫放到field
	//fmt.Printf("caller formatted: \n")
	l.SetFormatter(FormatterTypeText, true, false, nil)
	l.WithField(LevelInfo, "caller", l.ReturnCallerFormatted(), "a", "b")

	//hook
	ent := l.NewEntry()
	l.Out(ioutil.Discard)
	l.WithFieldsHook(ent, LevelFatal, logrus.Fields{"cmd.command": "testcmd0"}, "args")

}

func Test_newLoggerWithLogstash(t *testing.T) {

	p := &LoggerParameter{
		true,
		"log",
		"error",
		"113.35.204.62:5000",
		"local",
		"test",
	}

	l, _ := newLogger(p)

	//hook
	ent := l.NewEntry()
	l.Out(ioutil.Discard)
	l.WithFieldsHook(ent, LevelFatal, logrus.Fields{"cmd.command": "testcmd0"}, "args")

}


