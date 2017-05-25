package rtiddsgo

import (
	"errors"
	"unsafe"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I/usr/dds/rti_connext_dds-5.2.3/include -I/usr/dds/rti_connext_dds-5.2.3/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L/usr/dds/rti_connext_dds-5.2.3/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
// #include <ndds/ndds_c.h>
import "C"

type Participant struct {
	p *C.DDS_DomainParticipant
}

// New returns a new participant on "domain" with "qosProfileName" from
// "qosLibraryName". Default QoS is used in "qosLibraryName" is an empty string.
// Invoke p.Free() when done with the participant.
func New(domain int, qosLibraryName, qosProfileName string) (Participant, error) {
	p := Participant{}
	if len(qosLibraryName) == 0 {
		p.p = C.DDS_DomainParticipantFactory_create_participant(
			C.DDS_DomainParticipantFactory_get_instance(),
			C.DDS_DomainId_t(domain),
			&C.DDS_PARTICIPANT_QOS_DEFAULT,
			nil,
			C.DDS_STATUS_MASK_NONE)
	} else {
		p.p = C.DDS_DomainParticipantFactory_create_participant_with_profile(
			C.DDS_DomainParticipantFactory_get_instance(),
			C.DDS_DomainId_t(domain),
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_STATUS_MASK_NONE)
	}
	if p.p == nil {
		return p, errors.New("Failed to create a participant")
	}
	return p, nil
}

// Free deletes the participant.
func (p Participant) Free() {
	C.DDS_DomainParticipant_delete_contained_entities(p.p)
	C.DDS_DomainParticipantFactory_delete_participant(C.DDS_DomainParticipantFactory_get_instance(), p.p)
	p.p = nil
}

func (p Participant) Get() *C.DDS_DomainParticipant {
	return p.p
}

// GetUnsafe returns a pointer to the participant as an unsafe pointer.
// C types cannot be used in other packages, so Get() won't work outside
// rtiddsgo, in particular in whatever package the generated type code
// reside.
func (p Participant) GetUnsafe() unsafe.Pointer {
	return unsafe.Pointer(p.Get())
}
