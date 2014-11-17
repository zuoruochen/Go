package main

import "fmt"

func modify1(arr [10]int) {
	arr[0] = 10
	fmt.Println("in modify1(),arr values:",arr)
}

func modify2(arr *[10]int) {
	(*arr)[0] = 10
	fmt.Println("in modify2(),arr values:",*arr)
}



func main() {
	arr := [10]int {1,2,3,4,5}
	modify1(arr)
	fmt.Println("in main,arr values:",arr)
	modify2(&arr)
	fmt.Println("in main,arr values:",arr)
}
