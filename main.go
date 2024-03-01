package main

import (
	"fmt"
	ir "github.com/0187773933/IRController/v1/controller"
)

func main() {
	x := ir.New()
	fmt.Println( x )
	x.Transmit( "necx:0x70702" )
}