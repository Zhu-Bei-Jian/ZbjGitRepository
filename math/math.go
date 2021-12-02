package math

func Add(a int32,b int32)int32{
	return a+b
}
func Sub(a int32,b int32)int32{
	return a-b
}
func Fib(n int32)int32{
	if n==1 ||n==2 {
		return 1
	}
	return Fib(n-1)+Fib(n-2)
}
