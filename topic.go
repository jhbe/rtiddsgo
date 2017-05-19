package rtiddsgo

import (
	"errors"
	"unsafe"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I/opt/rti_connext_dds-5.2.3/include -I/opt/rti_connext_dds-5.2.3/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L/opt/rti_connext_dds-5.2.3/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
import "C"

type Topic struct {
	t *C.DDS_Topic
	p Participant
	name, typeName *C.char
}

func (p Participant)CreateTopic(name, typeName string, qosLibraryName, qosProfileName string) (Topic, error)  {
	t := Topic{p: p, name: C.CString(name), typeName: C.CString(typeName)}
	if len(qosLibraryName) == 0 {
		t.t = C.DDS_DomainParticipant_create_topic(
			t.p.p,
			t.name,
			t.typeName,
			&C.DDS_TOPIC_QOS_DEFAULT,
			nil,
			C.DDS_STATUS_MASK_NONE)
	} else {
		t.t = C.DDS_DomainParticipant_create_topic_with_profile(
			t.p.p,
			t.name,
			t.typeName,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_STATUS_MASK_NONE)
	}
	if t.t == nil {
		return t, errors.New("Failed to create a topic")
	}
	return t, nil
}

func (t Topic)Free() {
	C.DDS_DomainParticipant_delete_topic(t.p.p, t.t)
	t.t = nil

	C.free(unsafe.Pointer(t.name))
	C.free(unsafe.Pointer(t.typeName))
}

func (t Topic)description() *C.DDS_TopicDescription {
	if t.t == nil {
		return nil
	}
	return t.t._as_TopicDescription
}