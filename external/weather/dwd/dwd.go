package dwd

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L/usr/libs/x86_64-linux-gnu -leccodes
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "eccodes.h"
*/
import "C"
import (
	"unsafe"
)

type CodesValues struct {
	cCodesValues *C.codes_values
}

type CodesKeyValueList struct {
	cCodesKeyValueList *C.codes_key_value_list
}

type CodesHandle struct {
	charptr      unsafe.Pointer
	fileptr      *C.FILE
	cCodesHandle *C.codes_handle
}

type CodesMultiHandle struct {
	cCodesMultiHandle *C.codes_multi_handle
}

type CodesContext struct {
	cCodesContext *C.codes_context
}

type CodesNearest struct {
	cCodesNearest *C.codes_nearest
}

type CodesNearestValue struct {
	outlat   float32
	outlon   float32
	value    float32
	distance float32
	index    int
}

type CodesNearestValues []CodesNearestValue

type ProductKind int

// codes_handle* codes_grib_handle_new_from_file(codes_context* c, FILE* f, int* error);
// codes_handle* codes_bufr_handle_new_from_file(codes_context* c, FILE* f, int* error);
// codes_handle* codes_handle_new_from_message(codes_context* c, const void* data, size_t data_len);
// codes_handle* codes_grib_handle_new_from_multi_message(codes_context* c, void** data,
// codes_handle* codes_handle_new_from_message_copy(codes_context* c, const void* data, size_t data_len);
// codes_handle* codes_grib_handle_new_from_samples(codes_context* c, const char* sample_name);
// codes_handle* codes_bufr_handle_new_from_samples(codes_context* c, const char* sample_name);
// codes_handle* codes_handle_new_from_samples(codes_context* c, const char* sample_name);
// codes_handle* codes_handle_clone(const codes_handle* h);

func CodesHandleNewFromFile(c *CodesContext, f string, product ProductKind) (*CodesHandle, int) {
	var fp *C.FILE
	cfilename := C.CString(f)
	defer C.free(unsafe.Pointer(cfilename))
	cmode := C.CString("r")
	defer C.free(unsafe.Pointer(cmode))
	fp = C.fopen(cfilename, cmode)
	p := C.ProductKind(product)
	var cerr C.int
	codesHandle := C.codes_handle_new_from_file(c.cCodesContext, fp, p, &cerr)
	return &CodesHandle{
		fileptr:      fp,
		cCodesHandle: codesHandle,
	}, int(cerr)
}

func CodesHandleNewFromMessage(c *CodesContext, data []byte) *CodesHandle {
	p := C.malloc(C.size_t(len(data)))
	cBuf := (*[1 << 30]byte)(p)
	copy(cBuf[:], data)

	return &CodesHandle{
		charptr: p,
		cCodesHandle: C.codes_handle_new_from_message(
			c.cCodesContext,
			p,
			C.size_t(len(data)),
		),
	}
}

func CodesContextGetDefault() *CodesContext {
	return &CodesContext{
		cCodesContext: C.codes_context_get_default(),
	}
}

func CodesGetLong(h *CodesHandle, field string) (int64, int) {
	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))
	long := C.long(0)
	rtn := C.codes_get_long(h.cCodesHandle, cField, &long)
	return int64(long), int(rtn)
}

func CodesGetString(h *CodesHandle, field string) (string, int) {
	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))
	cMesg := make([]C.char, 20)
	cSize := C.size_t(0)
	rtn := C.codes_get_string(h.cCodesHandle, cField, &cMesg[0], &cSize)
	return C.GoStringN(&cMesg[0], C.int(cSize)), int(rtn)
}

func CodesGribNearestNew(h *CodesHandle) (*CodesNearest, int) {
	cErr := C.int(0)
	codesNearest := CodesNearest{
		cCodesNearest: C.codes_grib_nearest_new(h.cCodesHandle, &cErr),
	}
	return &codesNearest, int(cErr)
}

func CodesGribNearestFind(nearest *CodesNearest, h *CodesHandle, inlat float32, inlon float32, flags uint64) (CodesNearestValues, int) {
	cInlat := C.double(inlat)
	cInlon := C.double(inlon)
	cFlags := C.ulong(flags)
	cOutlats := make([]C.double, 4)
	cOutlons := make([]C.double, 4)
	cValues := make([]C.double, 4)
	cDistances := make([]C.double, 4)
	cIndices := make([]C.int, 4)
	cSize := C.size_t(4)

	rtn := int(C.codes_grib_nearest_find(nearest.cCodesNearest, h.cCodesHandle, cInlat, cInlon, cFlags, &cOutlats[0], &cOutlons[0], &cValues[0], &cDistances[0], &cIndices[0], &cSize))
	// This is fine.
	codesNearestValues := CodesNearestValues{}
	for n := 0; n < int(cSize); n++ {
		codesNearestValues = append(codesNearestValues, CodesNearestValue{
			outlat:   float32(cOutlats[n]),
			outlon:   float32(cOutlons[n]),
			value:    float32(cValues[n]),
			distance: float32(cDistances[n]),
			index:    int(cIndices[n]),
		})
	}

	return codesNearestValues, rtn
}

func CodesHandleDelete(h *CodesHandle) int {
	if h.charptr != nil {
		defer C.free(h.charptr)
	}
	if h.fileptr != nil {
		defer C.fclose(h.fileptr)
	}
	return int(C.codes_handle_delete(h.cCodesHandle))
}
