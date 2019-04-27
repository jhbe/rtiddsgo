package rtiddsgo

import (
	"errors"
)

// #include <ndds/ndds_c.h>
import "C"

type Publisher struct {
	pub *C.DDS_Publisher
	p   Participant
}

// CreatePublisher returns a new publisher with "qosProfileName" from
// "qosLibraryName". Default QoS is used if "qosLibraryName" is an empty string.
// Invoke p.Free() when done with the publisher.
func (p Participant) CreatePublisher(qosLibraryName, qosProfileName string) (Publisher, error) {
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

func (pub Publisher) Free() {
	C.DDS_DomainParticipant_delete_publisher(pub.p.p, pub.pub)
	pub.pub = nil
}
