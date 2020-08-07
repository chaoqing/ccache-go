package main

// #cgo CPPFLAGS: -I./
// #cgo LDFLAGS: -L./ -lccache -lz -lm
// #include <main.h>
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

const(
	maxArgsLen = 0xfff
)

func cMain(args []string){
	argc:=C.int(len(args))

	fmt.Printf("Got %v args: %v\n", argc, args)

	argv:=(*[maxArgsLen]*C.char)(C.alloc_string_slice(argc))
	defer C.free(unsafe.Pointer(argv))

	for i, arg := range args{
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	C.ccache_main(argc, (**C.char)(unsafe.Pointer(argv)))
}

func main(){
	args:=os.Args

	cMain(args)
}