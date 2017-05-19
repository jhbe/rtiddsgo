package rtiddsgo

import (
	"errors"
	"unsafe"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I/opt/rti_connext_dds-5.2.3/include -I/opt/rti_connext_dds-5.2.3/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L/opt/rti_connext_dds-5.2.3/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
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
// DataWriter.DW won't work outside rtiddsgo, in particular in whatever
// package the generated type code reside.
func (dw DataWriter) GetUnsafe() unsafe.Pointer {
	return unsafe.Pointer(dw.DW)
}
