package main

import (
	"fmt"
	"os"
	"path/filepath"
	// "time"
	ir "github.com/0187773933/IRController/v1/controller"
	utils "github.com/0187773933/IRController/v1/utils"
)

func main() {
	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )
	fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )

	x := ir.New( &config )
	// x.Scan()
	// x.PressKey( "power" )
	// x.Transmit( "necx:0x70702" )
	// time.Sleep( 10 * time.Second )
	// x.SaveKeyFile( "test_raw" )
	// x.TransmitKeyFile( "test_raw" )
	x.PressKey( "test_raw" )
}