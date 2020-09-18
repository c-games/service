package service

import (
	"reflect"
	"testing"
)

func Test_newIDGenerator(t *testing.T) {

	g, err := newIDGenerator(1)
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	id := g.GenerateInt64ID()

	if reflect.TypeOf(id).String() != "int64" {
		t.Fatalf("type error got %s want %s\n", reflect.TypeOf(id), "int64")
	}

}

func TestParseTime(t *testing.T) {

	//node 0
	var node0 int64 =0
	g0, _ := newIDGenerator(node0)

	var node99 int64 =99
	g99, _ := newIDGenerator(node99)

	var node1023 int64 =1023
	g1023, _ := newIDGenerator(node1023)

	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"0",
			args{g0.GenerateInt64ID()},
			node0,
		},
		{
			"1",
			args{g99.GenerateInt64ID()},
			node99,
		},
		{
			"0",
			args{g1023.GenerateInt64ID()},
			node1023,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IDParseNode(tt.args.id); got != tt.want {
				t.Errorf("ParseNode() = %d, want %d", got, tt.want)
			}
		})
	}
}
