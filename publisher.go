package rtiddsgo

import (
	"errors"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I/opt/rti_connext_dds-5.2.3/include -I/opt/rti_connext_dds-5.2.3/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L/opt/rti_connext_dds-5.2.3/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
// #include <ndds/ndds_c.h>
import "C"

type Publisher struct {
	pub *C.DDS_Publisher
	p Participant
}

func (p Participant)CreatePublisher(qosLibraryName, qosProfileName string) (Publisher, error) {
	pub := Publisher{p: p}
	if len(qosLibraryName) == 0 {
		pub.pub = C.DDS_DomainParticipant_create_publisher(
			pub.p.p,
			&C.DDS_PUBLISHER_QOS_DEFAULT,
			nil,
			C.DDS_STATUS_MASK_NONE)
	} else {
		pub.pub = C.DDS_DomainParticipant_create_publisher_with_profile(
			pub.p.p,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_STATUS_MASK_NONE)
	}
	if pub.pub == nil {
		return pub, errors.New("Failed to create a publisher")
	}
	return pub, nil
}

func (pub Publisher)Free() {
	C.DDS_DomainParticipant_delete_publisher(pub.p.p, pub.pub)
	pub.pub = nil
}
