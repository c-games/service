package main

import (
	"fmt"
	"os"

	"github.com/c-games/service"
)

func main() {
	app, err := service.NewService("serviceConfig.json")
	if err != nil {
		fmt.Printf("err=%s\n", err.Error())
		os.Exit(1)
	}

	l, err := app.GetLogger()
	l.IsFullPath(true)
	if err != nil {
		fmt.Printf("err=%s\n", err.Error())
		os.Exit(10)
	}

	err = l.WithFieldsFile(service.LevelError, service.Fields{"caller": l.ReturnCallerFormatted()}, "msg")
	if err != nil {
		fmt.Printf("err=%s\n", err.Error())
		os.Exit(20)
	}

	fmt.Printf("service start...\n")
	app.Hold()
}
