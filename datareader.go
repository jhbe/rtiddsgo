package rtiddsgo

import (
	"errors"
	"opendds/callbacks"
	"unsafe"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I/opt/rti_connext_dds-5.2.3/include -I/opt/rti_connext_dds-5.2.3/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L/opt/rti_connext_dds-5.2.3/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
// #include <ndds/ndds_c.h>
//
// struct CallbackInfo {
//   int onDataAvailableIndex;
// };
//
// void on_data_available(void* listener_data, DDS_DataReader* dataReader) {
//   OnDataAvailable(((struct CallbackInfo *)listener_data)->onDataAvailableIndex);
// }
//
import "C"

var onDataAvailableCallbacks = callbacks.New()

type DataReader struct {
	dr  *C.DDS_DataReader
	t   Topic
	sub Subscriber

	callbackInfo *C.struct_CallbackInfo
}
// DO NOT USE! Use the NewFOOTYPEDataReader instead.
func CreateDataReader(sub Subscriber, t Topic, qosLibraryName, qosProfileName string, onDataAvailable func()) (DataReader, error) {
	dr := DataReader{
		t:            t,
		sub:          sub,
		callbackInfo: (*C.struct_CallbackInfo)(C.malloc(C.sizeof_struct_CallbackInfo)),
	}
	dr.callbackInfo.onDataAvailableIndex = C.int(onDataAvailableCallbacks.Add(onDataAvailable))

	listener := C.struct_DDS_DataReaderListener{
		as_listener:       C.struct_DDS_Listener{listener_data: unsafe.Pointer(dr.callbackInfo)},
		on_data_available: C.DDS_DataReaderListener_DataAvailableCallback(C.on_data_available),
	}
	if len(qosLibraryName) == 0 {
		dr.dr = C.DDS_Subscriber_create_datareader(
			dr.sub.sub,
			dr.t.t._as_TopicDescription,
			&C.DDS_DATAREADER_QOS_DEFAULT,
			&listener,
			C.DDS_DATA_AVAILABLE_STATUS)
	} else {
		dr.dr = C.DDS_Subscriber_create_datareader_with_profile(
			dr.sub.sub,
			dr.t.t._as_TopicDescription,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			&listener,
			C.DDS_DATA_AVAILABLE_STATUS)
	}
	if dr.dr == nil {
		return dr, errors.New("Failed to create a datareader")
	}
	return dr, nil
}

func (dr DataReader) Free() {
	C.DDS_Subscriber_delete_datareader(dr.sub.sub, dr.dr)
	dr.dr = nil
	C.free(unsafe.Pointer(dr.callbackInfo))
}

// GetUnsafe returns a pointer to the data reader as an unsafe pointer.
// C types cannot be used in other packages, so directly referencing
// DataReader.dr won't work outside rtiddsgo, in particular in whatever
// package the generated type code reside.
func (dr DataReader) GetUnsafe() unsafe.Pointer {
	return unsafe.Pointer(dr.dr)
}
