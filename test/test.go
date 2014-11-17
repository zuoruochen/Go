package main

import "fmt"

func main(){
	xs := []float64{1.1,2.2,3.3}
	var sum float64 = 1.2
	var avg float64
	avg = sum/float64(len(xs))
	fmt.Print(avg)

/*
	for i := 1; i <= 100; i++ {
		switch{
			case i%3 == 0 && i%5 == 0 :
				fmt.Println("FizzBuzz")
			case i%3 == 0:
				fmt.Println("Fizz")
			case i%5 == 0 :
				fmt.Println("Buzz")
			default:
				fmt.Printf("%d\n",i)
		}
	}
*/
}

func myprintf(args... interface{}) {
	for _, arg := range args {
		switch arg.(type) {
			case int:
				fmt.Println(arg, "is an int value.")
			case string:
				fmt.Println(arg, "is a string value")
			case int32:
				fmt.Println(arg, "is an int32 value")
			case float64:
				fmt.Println(arg, "is an float64 value")
			default:
				fmt.Println(arg, "is an unknow type")
		}
	}
}
