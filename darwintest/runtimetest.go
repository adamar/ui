// 27 february 2014
package main

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include <objc/message.h>
// #include <objc/objc.h>
// #include <objc/runtime.h>
// /* TODO <objc/NSObjCRuntime.h not found?!?! */
// /* TODO this HAS to be unsafe */
// typedef unsigned long NSUInteger;
// /* avoid depending on Objective-C */
// #include <CoreGraphics/CGGeometry.h>
// /* cgo doesn't handle ... */
// id objc_msgSend_noargs(id obj, SEL sel) { return objc_msgSend(obj, sel); }
// id objc_msgSend_strarg(id obj, SEL sel, char *a) { return objc_msgSend(obj, sel, a); }
// id objc_msgSend_NSRect_uint_uint_bool(id obj, SEL sel, CGRect a, NSUInteger b, NSUInteger c, BOOL d) { return objc_msgSend(obj, sel, a, b, c, d); }
// id objc_msgSend_id(id obj, SEL sel, id a) { return objc_msgSend(obj, sel, a); }
// Class NilClass = Nil; /* for newtypes.go */
import "C"

func objc_getClass(class string) C.id {
	cclass := C.CString(class)
	defer C.free(unsafe.Pointer(cclass))

	return C.objc_getClass(cclass)
}

func sel_getUid(sel string) C.SEL {
	csel := C.CString(sel)
	defer C.free(unsafe.Pointer(csel))

	return C.sel_getUid(csel)
}

func main() {
	sel := sel_getUid("ourMethod")
	C.objc_msgSend_noargs(mk("hello", sel),
		sel)
}

const (
	NSBorderlessWindowMask = 0
	NSTitledWindowMask = 1 << 0
	NSClosableWindowMask = 1 << 1
	NSMiniaturizableWindowMask = 1 << 2
	NSResizableWindowMask = 1 << 3
	NSTexturedBackgroundWindowMask = 1 << 8
)

const (
//	NSBackingStoreRetained = 0			// "You should not use this mode."
//	NSBackingStoreNonretained = 1		// "You should not use this mode."
	NSBackingStoreBuffered = 2
)

var alloc = sel_getUid("alloc")

func wintest() {
	NSWindow := objc_getClass("NSWindow")
	NSWindowinit :=
		sel_getUid("initWithContentRect:styleMask:backing:defer:")
	makeKeyAndOrderFront := sel_getUid("makeKeyAndOrderFront:")

	// need an NSApplication first - see https://github.com/TooTallNate/NodObjC/issues/21
	NSApplication := objc_getClass("NSApplication")
	sharedApplication := sel_getUid("sharedApplication")
	C.objc_msgSend_noargs(NSApplication, sharedApplication)

	rect := C.CGRect{
		origin:	C.CGPoint{100, 100},
		size:		C.CGSize{320, 240},
	}
	style := C.NSUInteger(NSTitledWindowMask | NSClosableWindowMask)
	backing := C.NSUInteger(NSBackingStoreBuffered)
	deferx := C.BOOL(C.YES)
	window := C.objc_msgSend_noargs(NSWindow, alloc)
	window = C.objc_msgSend_NSRect_uint_uint_bool(window, NSWindowinit, rect, style, backing, deferx)
	C.objc_msgSend_id(window, makeKeyAndOrderFront, window)
	select{}
}

func helloworld() {
	_hello := C.CString("hello, world\n")
	defer C.free(unsafe.Pointer(_hello))

	NSString := objc_getClass("NSString")
	stringWithUTF8String :=
		sel_getUid("stringWithUTF8String:")
	str := C.objc_msgSend_strarg(NSString,
		stringWithUTF8String,
		_hello)
	UTF8String := sel_getUid("UTF8String")
	res := C.objc_msgSend_noargs(str,
			UTF8String)
	cres := (*C.char)(unsafe.Pointer(res))
	fmt.Printf("%s", C.GoString(cres))
}