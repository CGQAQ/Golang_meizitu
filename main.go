package main

import (
	"meizitu/meizitu"
)

var a = 5

func main() {
	// debug can't read os.Stdin

	meizi := new(meizitu.Meizitu)
	meizi.Run()

	//var q meizitu.Pushable = new(meizitu.Queue)
	//q.Push("1", "2")
	//fmt.Println(q.Pop())
	//q.Push("1", "2")
	//fmt.Println(q.Pop())
	//fmt.Println(q.Pop())
	//fmt.Println(q.Pop())
}
