#!/usr/bin/env bash

ROOT_DIR=$(dirname $(realpath $0))/..

CCACHE_URL=https://github.com/ccache/ccache/releases/download/v3.7.11/ccache-3.7.11.tar.gz

CCACHE_DIR=$ROOT_DIR/ccache-c

if [ ! -f $CCACHE_DIR/configure ]; then
  echo "========= downloading ccache ========="
  mkdir -p $CCACHE_DIR
  curl -sL $CCACHE_URL | tar xz  -C $CCACHE_DIR --strip-components=1
fi

cd $CCACHE_DIR
./configure

cat >> Makefile <<"EOF"
.PHONY: libccache
libccache: src/libccache.a src/zlib/libz.a
src/libccache.a: $(ccache_objs)
	$(if $(quiet),@echo "  AR       $@")
	$(Q)$(AR) cr $@ $(ccache_objs)
	$(if $(quiet),@echo "  RANLIB   $@")
	$(Q)$(RANLIB) $@
EOF

make libccache