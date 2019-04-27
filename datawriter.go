package rtiddsgo

import (
	"errors"
	"unsafe"
)

// #include <ndds/ndds_c.h>
import "C"

type DataWriter struct {
	DW  *C.DDS_DataWriter
	t   Topic
	pub Publisher
}

// DO NOT USE! Use the NewFOOTYPEDataWriter instead.
func CreateDataWriter(pub Publisher, t Topic, qosLibraryName, qosProfileName string) (DataWriter, error) {
	dw := DataWriter{t: t, pub: pub}
	if len(qosLibraryName) == 0 {
		dw.DW = C.DDS_Publisher_create_datawriter(
			dw.pub.pub,
			dw.t.t,
			&C.DDS_DATAWRITER_QOS_DEFAULT,
			nil,
			C.DDS_STATUS_MASK_NONE)
	} else {
		dw.DW = C.DDS_Publisher_create_datawriter_with_profile(
			dw.pub.pub,
			dw.t.t,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_STATUS_MASK_NONE)
	}
	if dw.DW == nil {
		return dw, errors.New("Failed to create a datawriter")
	}
	return dw, nil
}

func (dw DataWriter) Free() {
	C.DDS_Publisher_delete_datawriter(dw.pub.pub, dw.DW)
	dw.DW = nil
}

// GetUnsafe returns a pointer to the data writer as an unsafe pointer.
// C types cannot be used in other packages, so directly referencing
// DataWriter.DW won't work outside godds, in particular in whatever
// package the generated type code reside.
func (dw DataWriter) GetUnsafe() unsafe.Pointer {
	return unsafe.Pointer(dw.DW)
}
