#ifndef __MAIN_H__
#define __MAIN_H__

#include <stdlib.h>
static void* alloc_string_slice(int len){
    return malloc(sizeof(char*)*len);
}

int ccache_main(int argc, char *argv[]);

#endif