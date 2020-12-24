package main

/*
#include <stdio.h>

void sayHello() {
    printf("hello world");
}
 */
import "C"

func main() {
	C.sayHello()
}