package math

import (
	"testing"
)
// 1 2 3 4 5
// 1 1 2 3 5
func TestFib(t *testing.T){
	var (
		inputN int32=4
		expected int32=3
	)
	actual:=Fib(4)
	if actual!=expected{
		t.Errorf("Fib(%d) = %d \n expected=%d ",inputN,actual,expected)
	}
}

func TestFib2(t *testing.T){
	var table=[]struct{
		input int32
		expected int32
	}{
		{1,1},
		{1,1},
		{2,3},
		{3,5},
		{4,8},
	}
	for _,v:=range table{
		actual:=Fib(v.input)
		if actual!=v.expected{
			t.Errorf("Fib(%d) = %d \n expected=%d ",v.input,actual,v.expected)
		}
	}

}
