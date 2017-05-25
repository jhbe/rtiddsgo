package rtiddsgo

import (
	"errors"
)

// #include <ndds/ndds_c.h>
import "C"

type Subscriber struct {
	sub *C.DDS_Subscriber
	p Participant
}

func (p Participant)CreateSubscriber(qosLibraryName, qosProfileName string) (Subscriber, error) {
	sub := Subscriber{p: p}
	if len(qosLibraryName) == 0 {
		sub.sub = C.DDS_DomainParticipant_create_subscriber(
			sub.p.p,
			&C.DDS_SUBSCRIBER_QOS_DEFAULT,
			nil,
			C.DDS_STATUS_MASK_NONE)
	} else {
		sub.sub = C.DDS_DomainParticipant_create_subscriber_with_profile(
			sub.p.p,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_STATUS_MASK_NONE)
	}
	if sub.sub == nil {
		return sub, errors.New("Failed to create a subscriber")
	}
	return sub, nil
}

func (sub Subscriber)Free() {
	C.DDS_DomainParticipant_delete_subscriber(sub.p.p, sub.sub)
	sub.sub = nil
}
