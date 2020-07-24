package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {

	l, _ := newLogger("log", "info","114.34.35.125:8100","local","srv")

	fmt.Printf("defalut callReprter: \n")
	l.SetFormatter(FormatterTypeText, true,true,nil)

	l.WithField(LevelInfo, "text", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "text", 1, "a", "b")

	fmt.Printf("custom callReprter: \n")

	//方法一，設定成formatter
	//但行數會錯，除非每次都設定
	l.SetFormatter(FormatterTypeText, true,true, l.CallerParser(l.ReturnCaller()))
	l.WithField(LevelInfo, "json", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "json", 1, "a", "b")
	_ = l.WithFieldsFile(LevelInfo, Fields{"1": 1, "2": 3}, "a", "b")

	//方法二，關閉reportCaller, 自己呼叫放到field
	fmt.Printf("caller formatted: \n")
	l.SetFormatter(FormatterTypeText, true,false, nil)
	l.WithField(LevelInfo, "caller", l.ReturnCallerFormatted(), "a", "b")

	//
	ent:=l.NewEntry()
	lv,_:=logrus.ParseLevel("info")
	ent.WithField("key","value").Log(lv,"data")

	//
	ent.WithField("key","value").Error("error")
}

func TestLogger_fromService(t *testing.T){
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "service") {
		wd = filepath.Dir(wd)
	}
	//t.Logf("wd = %s", wd)

	app, err := NewService("serviceConfig.json")
	if err != nil {
		t.Fatalf("err=%s\n", err.Error())
	}

	l,_:=app.GetLogger()
	l.SetFormatter(FormatterTypeText, true,false,nil)

	l.WithField(LevelInfo, "text", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "text", 1, "a", "b")
	_ = l.WithFieldFile(LevelInfo, "", "", "a", "b")

	l.SetFormatter(FormatterTypeJSON, true,false,nil)
	l.WithField(LevelInfo, "json", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "json", 1, "a", "b")
	_ = l.WithFieldsFile(LevelInfo, Fields{"1": 1, "2": 3}, "a", "b")
}
