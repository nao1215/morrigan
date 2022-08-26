//go:build linux
// +build linux

// Package gocrypt provides wrappers around functions available in crypt.h
//
// It wraps around the GNU specific extension (crypt) when the reentrant version
// (crypt_r) is unavailable. The non-reentrant version is guarded by a global lock
// so as to be safely callable from concurrent goroutines.
package gocrypt

/*
#include <features.h>
#ifdef __GLIBC__
#include <gnu/libc-version.h>
unsigned int get_glibc_minor_version(void) {
	return __GLIBC_MINOR__;
}
#else
unsigned int get_glibc_minor_version(void) {
	return 0;
}
#endif
*/
import "C"

// This function is specific to the tests. It basically checks if
// we're running with glibc and if the used glibc is new enough
func checkGlibCVersion() bool {
	c_minor := C.get_glibc_minor_version()
	return c_minor >= 17
}
