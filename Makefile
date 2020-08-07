.PHONY: libccache FORCE all
all: ccache

ccache: clib
	go build -o $@

clib: FORCE
	@scripts/build.sh
	cp ccache-c/src/libccache.a ./
	cp ccache-c/src/zlib/libz.a ./

FORCE: