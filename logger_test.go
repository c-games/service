package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {

	l, _ := newLogger("log", "info")
	fmt.Printf("defalut callReprter: \n")
	l.SetFormatter(FormatterTypeText, true,true,nil)

	l.WithField(LevelInfo, "text", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "text", 1, "a", "b")

	fmt.Printf("custom callReprter: \n")

	l.SetFormatter(FormatterTypeJSON, true,true, l.CallerParser(l.ReturnCaller()))
	l.WithField(LevelInfo, "json", "1", "a", "b")
	_ = l.WithFieldFile(LevelInfo, "json", 1, "a", "b")
	_ = l.WithFieldsFile(LevelInfo, Fields{"1": 1, "2": 3}, "a", "b")
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
