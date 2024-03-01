package controller

import (
	"fmt"
	// "bufio"
	"bytes"
	"os/exec"
	// "regexp"
	"io/ioutil"
	"strings"
	"runtime"
	utils "github.com/0187773933/IRController/v1/utils"
	gousb "github.com/google/gousb"
)

const TRANSMIT_AND_RECEIVE = 1
const TRANSMIT_ONLY = 2
const RECEIVE_ONLY = 3

// 1.) Find Special Character File for USB Controller
// lsusb
// lsusb | awk '{print $2 ":" $4}' | sed 's/://g'
// ls -lah /dev/
// ls -lah /dev/input/by-path/
// ls -lah /dev/input/by-id
// mode2 -l
// lsusb -s 001:018 -v
// ir-ctl -d /dev/lirc0 --features

type IRController struct {
	DevicePaths []string `yaml:"device_paths"`
	DevicePath string `yaml:"device_path"`
	DeviceType int `yaml:"device_type"`
}

type USBDevice struct {
	Manufacturer string `yaml:"manufacturer"`
	Product string `yaml:"product"`
	SerialNumber string `yaml:"serial_number"`
	VendorID string `yaml:"vendor_id"`
	ProductID string `yaml:"product_id"`
	BusNumber string `yaml:"bus_number"`
	AddressNumber string `yaml:"address_number"`
}

func LinuxFindDevices() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices , _ := ctx.OpenDevices( func( desc *gousb.DeviceDesc ) bool {
		// fmt.Printf( "vid=%s pid=%s\n" , desc.Vendor , desc.Product )
		return true
	})
	var parsed_devices []USBDevice
	for _ , device := range devices {
		manufacturer , _ := device.Manufacturer()
		product , _ := device.Product()
		serial_number , _ := device.SerialNumber()
		device_string := device.String()
		vid := strings.Split( device_string , "vid=" )[ 1 ][ :4 ]
		pid := strings.Split( strings.Split( device_string , "pid=" )[ 1 ] , "," )[ 0 ]
		bus := strings.Split( strings.Split( device_string , "bus=" )[ 1 ] , "," )[ 0 ]
		addr := strings.Split( strings.Split( device_string , "addr=" )[ 1 ] , "," )[ 0 ]
		// fmt.Printf( "%s:%s === %s:%s === %s === %s === %s\n" , vid , pid , bus , addr , manufacturer , product , serial_number )
		device.Close()
		parsed_devices = append( parsed_devices , USBDevice{
			Manufacturer: manufacturer ,
			Product: product ,
			SerialNumber: serial_number ,
			VendorID: vid ,
			ProductID: pid ,
			BusNumber: bus ,
			AddressNumber: addr ,
		})
	}
	for _ , device := range parsed_devices {
		utils.PrettyPrint( device )
	}
}

func LinuxGetLIRCDevices() ( results []string ) {
	files, err := ioutil.ReadDir( "/dev" )
	if err != nil {
		fmt.Println("Error reading /dev directory:", err)
		return
	}
	for _ , file := range files {
		file_name := file.Name()
		if strings.HasPrefix( file_name, "lirc" ) {
			results = append( results , "/dev/" + file_name )
		}
	}
	return
}

func LinuxGetDeviceType( device_path string ) ( result int ) {
	cmd := exec.Command( "ir-ctl" , "-d" , device_path , "--features" )
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println( "Error executing command:" , err )
		return 0 // You might want to handle errors differently
	}
	output := out.String()
	has_receive := strings.Contains( output , "Device can receive" )
	has_send := strings.Contains( output , "Device can send" )
	switch {
		case has_receive && has_send:
			return TRANSMIT_AND_RECEIVE
		case has_send:
			return TRANSMIT_ONLY
		case has_receive:
			return RECEIVE_ONLY
		default:
			return 0 // Unknown or no capabilities detected
	}
}

// https://github.com/google/gousb/tree/master?tab=readme-ov-file
// https://github.com/libusb/libusb/wiki
// sudo apt-get install libusb-1.0-0-dev -y
// https://pkg.go.dev/github.com/google/gousb#DeviceDesc
func NewLinux() ( result IRController ) {
	result.DevicePaths = LinuxGetLIRCDevices()
	result.DevicePath = result.DevicePaths[ 0 ]
	result.DeviceType = LinuxGetDeviceType( result.DevicePath )
	return
}

// 2.) Get Features of USB Controller
// ir-ctl -f
func New() ( result IRController ) {
	switch os := runtime.GOOS; os {
		case "linux":
			result = NewLinux()
			break;
		case "windows":
			fmt.Println( "not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "not implemented for mac osx" )
			break;
		default:
			fmt.Println( "not implemented for" , os )
			break;
	}
	return
}

func ( irc *IRController ) TransmitLinux( code string ) {
	cmd := exec.Command( "ir-ctl" , "-d" , irc.DevicePath , "-S" , code )
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil { fmt.Println( "Error executing command:" , err ) }
	fmt.Println( out.String() )
}

func ( irc *IRController ) Transmit( code string ) {
	switch os := runtime.GOOS; os {
		case "linux":
			irc.TransmitLinux( code )
			break;
		case "windows":
			fmt.Println( "transmit not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "transmitnot implemented for mac osx" )
			break;
		default:
			fmt.Println( "transmit not implemented for" , os )
			break;
	}
}

func ( irc *IRController ) Receive() {
	switch os := runtime.GOOS; os {
		case "linux":
			fmt.Println( "receive not implemented for linux" )
			break;
		case "windows":
			fmt.Println( "receive not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "receive not implemented for mac osx" )
			break;
		default:
			fmt.Println( "receive not implemented for" , os )
			break;
	}
}

func ( irc *IRController ) Scan() {
	switch os := runtime.GOOS; os {
		case "linux":
			fmt.Println( "scan not implemented for linux" )
			break;
		case "windows":
			fmt.Println( "scan not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "scan not implemented for mac osx" )
			break;
		default:
			fmt.Println( "scan not implemented for" , os )
			break;
	}
}