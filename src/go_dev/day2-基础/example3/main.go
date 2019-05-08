package main

import (
    "fmt"
)
func modify(a int){
    a=10
    return 
}

func modify1(a *int){
    *a=10
}

func main(){
    var a = 100
    var b chan int = make( chan int,1)

    fmt.Println("a=", a)
    fmt.Println("b=", b)

    modify(a)
    fmt.Println("a=",a)

    modify1(&a)
    fmt.Println("a=",a)
}
