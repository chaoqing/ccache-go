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
	"path/filepath"
	"regexp"
	"strings"
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


type PatternSet struct {
	Values       []string
	IsRegex      bool
	regexPattern []*regexp.Regexp
}

func (t *PatternSet) compile() {
	if t.IsRegex{
		for _, s := range t.Values {
			t.regexPattern = append(t.regexPattern, regexp.MustCompile(s))
		}
	}
}

func (t *PatternSet) match(v string) bool{
	if t.IsRegex{
		for _, p := range t.regexPattern{
			if p.MatchString(v) {
				return true
			}
		}
		return false
	}else{
		for _, p := range t.Values{
			if p==v {
				return true
			}
		}
		return false
	}
}

type Rule struct {
	command *PatternSet
	arg *PatternSet
}

func(r *Rule) match(command, arg string) bool{
	return r.arg.match(arg) && (len(r.command.Values)==0 || r.command.match(command))
}

const (
	extraOption = "--ccache-skip"
)

var(
	rules = []*Rule{
		&Rule{
			command: &PatternSet{Values: []string{"icp?c"}, IsRegex: true},
			arg: &PatternSet{Values: []string{"-xAVX", "-xCORE_AVX2"}},
			},
	}
)

func PrepareArgs(origArgs []string) []string {
	args := make([]string, 0, len(origArgs)*2)
	args = append(args, origArgs[0])

	command := filepath.Base(origArgs[0])
	if command=="ccache" && len(origArgs)>1 && !strings.HasPrefix(origArgs[1], "-"){
		command = filepath.Base(origArgs[1])
	}
	for _, arg := range origArgs[1:]{
		for _, r := range rules{
			if r.match(command, arg){
				args = append(args, extraOption)
				break
			}
		}
		args = append(args, arg)
	}
	return args
}

func init(){
	for _, r := range rules{
		r.command.compile()
		r.arg.compile()
	}
}


func main(){
	args := PrepareArgs(os.Args)
	log.Debugln(args)

	cMain(args)
}



