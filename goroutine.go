package main

import (
	"math/rand"
	"fmt"
)

//same func goroutine shares variable?

func goroutine(){
	a := rand.Int()
	fmt.Println(a)
}

func main() {
	for{
		go goroutine()
	}

}