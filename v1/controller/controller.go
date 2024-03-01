package controller

import (
	"fmt"
	"bufio"
	// "os/signal"
	// "syscall"
	"os"
	"bytes"
	"path/filepath"
	"os/exec"
	// "regexp"
	"io/ioutil"
	"strings"
	"runtime"
	utils "github.com/0187773933/IRController/v1/utils"
	types "github.com/0187773933/IRController/v1/types"
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
	Config *types.ConfigFile `yaml:"-"`
	Remote string `yaml:"remote"`
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
		fmt.Println( "Error reading /dev directory:" , err )
		return
	}
	for _ , file := range files {
		file_name := file.Name()
		if strings.HasPrefix( file_name , "lirc" ) {
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
func NewLinux( config *types.ConfigFile ) ( result IRController ) {
	result.Config = config
	result.Remote = result.Config.DefaultRemote
	result.DevicePaths = LinuxGetLIRCDevices()
	result.DevicePath = result.DevicePaths[ 0 ]
	result.DeviceType = LinuxGetDeviceType( result.DevicePath )
	if err := os.MkdirAll( result.Config.KeySaveFileBasePath , 0755 ); err != nil {
		fmt.Println( "Failed to create directory: %s" , err )
	}
	return
}

// 2.) Get Features of USB Controller
// ir-ctl -f
func New( config *types.ConfigFile ) ( result IRController ) {
	switch os := runtime.GOOS; os {
		case "linux":
			result = NewLinux( config )
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

func ( irc *IRController ) ScanLinux() {
	cmd := exec.Command( "sudo" , "ir-keytable" , "-v" , "-t" ,
		"-p", "rc-5,rc-5-sz,jvc,samsung,sony,nec,sanyo,mce_kbd,rc-6,sharp,xmp", "-s", "rc0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting Cmd: %v\n", err)
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for Cmd: %v\n", err)
		return
	}
}

// func ( irc *IRController ) ScanRawLinux( save_path string ) {
// 	cmd := exec.Command( "ir-ctl" , "-d" , irc.DevicePath , "--receive" , save_path )
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Start(); err != nil {
// 		fmt.Printf("Error starting Cmd: %v\n", err)
// 		return
// 	}
// 	if err := cmd.Wait(); err != nil {
// 		fmt.Printf("Error waiting for Cmd: %v\n", err)
// 		return
// 	}
// }


// ir-ctl -d /dev/lirc0 --receive=samnsung_power.key
func ( irc *IRController ) ScanRawLinux( save_path string ) {
	fmt.Println( "Press Single Button on IR Remote Once , then space or enter on computer keyboard to stop recording" )
	key_name := fmt.Sprintf( "%s.key" , save_path )
	key_path := filepath.Join( irc.Config.KeySaveFileBasePath , key_name )
	cmd := exec.Command( "ir-ctl" , "-d" , irc.DevicePath , "-1" , "-r" , "--mode2" , fmt.Sprintf( "--receive=%s" , key_path ) )

	// Directly attach command's stdout and stderr to the os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command asynchronously.
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting Cmd: %v\n", err)
		return
	}

	// Set up channel and goroutine for handling user interrupt.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			char , _ , err := reader.ReadRune()
			if err != nil {
				fmt.Printf( "Error reading rune: %v\n" , err )
				return
			}

			// Exit on 'Enter' or 'Space' key press.
			if char == '\r' || char == '\n' || char == ' ' {
				fmt.Println( "Stopping IR capture..." )
				if err := cmd.Process.Signal( os.Interrupt ); err != nil {
					fmt.Printf( "Error sending interrupt: %v\n" , err )
				}
				return
			}
		}
	}()

	// Wait for the command to complete or to be interrupted.
	if err := cmd.Wait(); err != nil {
		fmt.Printf( "IR capture stopped: %v\n" , err )
		return
	}
}

func ( irc *IRController ) TransmitRawLinux( save_path string ) {
	key_path := fmt.Sprintf( "%s.key" , save_path )
	cmd := exec.Command( "ir-ctl" , "-d" , irc.DevicePath , fmt.Sprintf( "--send=%s" , key_path ) )
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error sending IR command: %v\n", err)
	} else {
		fmt.Println("IR command sent successfully.")
	}
}

func (irc *IRController) PressKeyLinux(key_name string) {
    remote := irc.Config.Remotes[irc.Remote] // Access the Remote instance
    key, exists := remote.Keys[key_name]     // Access the Key from the Remote's Keys map

    if !exists {
        fmt.Println("Key does not exist:", key_name)
        return
    }

    if key.Code != "" {
        irc.TransmitLinux(key.Code)
    } else {
        // If Code is empty, assume KeyPath should be used (adjust logic as needed)
        fmt.Println("then we need to send as a key file")
        // Here you would presumably use key.KeyPath, but it's not shown in this snippet
    }
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

func ( irc *IRController ) PressKey( key string ) {
	switch os := runtime.GOOS; os {
		case "linux":
			irc.PressKeyLinux( key )
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
			irc.ScanLinux()
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

func ( irc *IRController ) ScanRaw( key string ) {
	switch os := runtime.GOOS; os {
		case "linux":
			irc.ScanRawLinux( key )
			break;
		case "windows":
			fmt.Println( "scan raw not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "scan raw not implemented for mac osx" )
			break;
		default:
			fmt.Println( "scan raw not implemented for" , os )
			break;
	}
}

func ( irc *IRController ) TransmitRaw( key string ) {
	switch os := runtime.GOOS; os {
		case "linux":
			irc.TransmitRawLinux( key )
			break;
		case "windows":
			fmt.Println( "scan raw not implemented for windows" )
			break;
		case "darwin":
			fmt.Println( "scan raw not implemented for mac osx" )
			break;
		default:
			fmt.Println( "scan raw not implemented for" , os )
			break;
	}
}