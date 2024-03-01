package utils

import (
	"fmt"
	"time"
	"strings"
	"unicode"
	"strconv"
	ulid "github.com/oklog/ulid/v2"
	binary "encoding/binary"
	json "encoding/json"
	// "strings"
	ioutil "io/ioutil"
	yaml "gopkg.in/yaml.v2"
	types "github.com/0187773933/IRController/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v3/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func IToB( v uint64 ) []byte {
	b := make( []byte , 8 )
	binary.BigEndian.PutUint64( b , v )
	return b
}

func GetFormattedTimeString() ( result string ) {
	location , _ := time.LoadLocation( "America/New_York" )
	time_object := time.Now().In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}
func GetFormattedTimeStringOBJ() ( result_string string , result_time time.Time ) {
	location , _ := time.LoadLocation( "America/New_York" )
	result_time = time.Now().In( location )
	month_name := strings.ToUpper( result_time.Format( "Jan" ) )
	milliseconds := result_time.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , result_time.Day() , month_name , result_time.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , result_time.Hour() , result_time.Minute() , result_time.Second() , milliseconds )
	result_string = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}
func FormatTime( input_time *time.Time ) ( result string ) {
	location , _ := time.LoadLocation( "America/New_York" )
	time_object := input_time.In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

func PrettyPrint( input interface{} ) {
	jd , _ := json.MarshalIndent( input , "" , "  " )
	fmt.Println( string( jd ) )
}

func WriteJSON( filePath string , data interface{} ) {
	file, _ := json.MarshalIndent( data , "" , " " )
	_ = ioutil.WriteFile( filePath , file , 0644 )
}

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( file , &result )
	if error != nil { panic( error ) }
	return
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	bolt_db_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tBolt DB Key === %s\n" , bolt_db_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n\n" , admin_password )
}