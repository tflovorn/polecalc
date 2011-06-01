package main

import "fmt"
import "./mesh2d"

func main() {
    cmesh := mesh2d.Square(16)
    for i := 0; i < 33; i++ {
        fmt.Println(<-cmesh)
    }
}
