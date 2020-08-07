package main

// #cgo LDFLAGS: -L./ -lccache -lz -lm
//
// #include <stdlib.h>
// static void* alloc_string_slice(int len){
// return malloc(sizeof(char*)*len);
// }
//
// int ccache_main(int argc, char *argv[]);
import "C"

import (
	log "github.com/sirupsen/logrus"
	"os"
	"unsafe"
)

const(
	maxArgsLen = 0xfff
)

func cMain(args []string){
	argc:=C.int(len(args))

	log.Debugf("Got %v args: %v\n", argc, args)

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