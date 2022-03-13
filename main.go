package main

import (
	"errors"
	"fmt"
)

func getBeginAndEnd(len int, parts int, which int) (int, int) {
	if len < 0 || parts <= 0 || which <= 0 || which > parts {
		panic(errors.New("getBeginAndEnd函数的参数不合法！"))
	}
	begin := (len * (which - 1)) / parts
	end := ((len * which) / parts) - 1
	//循环时只需要 begin <= i <= end
	return begin, end
}

func main() {

	len := 20
	parts := 3

	for i := 1; i <= parts; i++ {
		fmt.Println(getBeginAndEnd(len, parts, i))
	}

	// startTime := time.Now()
	// run.StartScan()
	// fmt.Printf("程序执行总时长为：[%s]\n", time.Since(startTime).String())

}
