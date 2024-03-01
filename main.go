package main

import (
	"fmt"
	// "time"
	ir "github.com/0187773933/IRController/v1/controller"
)

func main() {
	x := ir.New()
	fmt.Println( x )
	// x.Transmit( "necx:0x70702" )
	// time.Sleep( 10 * time.Second )
	// x.ScanRaw( "volume_up" )
	x.TransmitRaw( "volume_up" )
}