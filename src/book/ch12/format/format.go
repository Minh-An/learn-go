package format

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

//Any ..
func Any(val interface{}) string {
	return formatAtom(reflect.ValueOf(val))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Ptr, reflect.Func:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value "
	}
}

func main() {
	var a []string
	b := 34
	c := "hello?"
	d := false
	var e io.Writer
	e = os.Stdout
	var f reflect.Value
	fmt.Println(Any(a))
	fmt.Println(Any(b))
	fmt.Println(Any(c))
	fmt.Println(Any(d))
	fmt.Println(Any(e))
	fmt.Println(Any(f))
}
