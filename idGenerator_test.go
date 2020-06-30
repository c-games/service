package service

import (
	"reflect"
	"testing"
)

func Test_newIDGenerator(t *testing.T) {


	g,err:=newIDGenerator(1)
	if err!=nil{
		t.Fatalf("err %s",err.Error())
	}
	id:= g.GenerateInt64ID()

	if reflect.TypeOf(id).String()!="int64"{
		t.Fatalf("type error got %s want %s\n",reflect.TypeOf(id),"int64")
	}




}
