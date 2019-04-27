package rtiddsgo

import (
	"errors"
	"log"
	"sync"
	"unsafe"
)

// #include <ndds/ndds_c.h>
import "C"

type Processor interface {
	DataAvailable() error
}

type DataReader struct {
	dr  *C.DDS_DataReader
	ws  *C.DDS_WaitSet
	sc  *C.DDS_Condition
	gc  *C.DDS_Condition
	t   Topic
	sub Subscriber

	done sync.WaitGroup
}

// DO NOT USE! Use the NewFOOTYPEDataReader instead.
func CreateDataReader(sub Subscriber, t Topic, qosLibraryName, qosProfileName string, p Processor) (*DataReader, error) {
	dr := &DataReader{
		t:   t,
		sub: sub,
		ws:  C.DDS_WaitSet_new(),
		gc:  (*C.DDS_Condition)(C.DDS_GuardCondition_new()),
	}

	if len(qosLibraryName) == 0 {
		dr.dr = C.DDS_Subscriber_create_datareader(
			dr.sub.sub,
			dr.t.t._as_TopicDescription,
			&C.DDS_DATAREADER_QOS_DEFAULT,
			nil,
			C.DDS_DATA_AVAILABLE_STATUS)
	} else {
		dr.dr = C.DDS_Subscriber_create_datareader_with_profile(
			dr.sub.sub,
			dr.t.t._as_TopicDescription,
			C.CString(qosLibraryName),
			C.CString(qosProfileName),
			nil,
			C.DDS_DATA_AVAILABLE_STATUS)
	}
	if dr.dr == nil {
		return dr, errors.New("Failed to create a datareader")
	}

	dr.sc = (*C.DDS_Condition)(C.DDS_Entity_get_statuscondition((*C.struct_DDS_EntityImpl)(dr.dr)))

	if rc := C.DDS_WaitSet_attach_condition(dr.ws, dr.sc); rc != C.DDS_RETCODE_OK {
		return dr, errors.New("Failed to attach statuscondition to datareader")
	}
	if rc := C.DDS_WaitSet_attach_condition(dr.ws, dr.gc); rc != C.DDS_RETCODE_OK {
		return dr, errors.New("Failed to attach guardcondition to datareader")
	}

	dr.done.Add(1)
	go dr.process(p)

	return dr, nil
}

func (dr *DataReader) Free() {
	C.DDS_GuardCondition_set_trigger_value((*C.DDS_GuardCondition)(dr.gc), C.DDS_BOOLEAN_TRUE);
	dr.done.Wait()

	C.DDS_GuardCondition_delete((*C.DDS_GuardCondition)(dr.gc))
	dr.gc = nil
	C.DDS_WaitSet_delete(dr.ws)
	dr.ws = nil
	C.DDS_Subscriber_delete_datareader(dr.sub.sub, dr.dr)
	dr.dr = nil
}

func (dr *DataReader) process(p Processor) {
	for {
		var activeConditions C.struct_DDS_ConditionSeq
		C.DDS_ConditionSeq_initialize(&activeConditions)

		rc := C.DDS_WaitSet_wait(dr.ws, &activeConditions, &C.DDS_DURATION_INFINITE)
		if rc != C.DDS_RETCODE_OK {
			log.Fatal("WaitSet_wait failed")
		}

		for i := 0; i < int(C.DDS_ConditionSeq_get_length(&activeConditions)); i++ {
			condition := C.DDS_ConditionSeq_get(&activeConditions, C.DDS_Long(i))
			if condition == dr.gc {
				dr.done.Done()
				return
			}
			if condition == dr.sc {
				dr.processStatusChange(p)
			}
		}
	}

}

func (dr *DataReader) processStatusChange(p Processor) {
	statusChangesMask := C.DDS_Entity_get_status_changes((*C.struct_DDS_EntityImpl)(dr.dr))

	// Read Conditions
	if statusChangesMask & C.DDS_DATA_AVAILABLE_STATUS != 0 {
		//log.Println("DDS_DATA_AVAILABLE_STATUS")
		p.DataAvailable()
	}
	if statusChangesMask & C.DDS_DATA_ON_READERS_STATUS != 0 {
		//log.Println("DDS_DATA_ON_READERS_STATUS")

	}

	// Status Conditions
	if statusChangesMask & C.DDS_INCONSISTENT_TOPIC_STATUS != 0 {
		var its C.struct_DDS_InconsistentTopicStatus
		C.DDS_Topic_get_inconsistent_topic_status(dr.t.t, &its)
		//log.Println("DDS_INCONSISTENT_TOPIC_STATUS")
	}
/*
	if statusChangesMask & C.DDS_OFFERED_DEADLINE_MISSED_STATUS != 0 {
		var its C.struct_DDS_InconsistentTopicStatus
		C.DDS_DataWriter_get_offered_deadline_missed_status()
	}
	if statusChangesMask & C.DDS_OFFERED_INCOMPATIBLE_QOS_STATUS != 0 {
		C.DDS_DataWriter_get_offered_incompatible_qos_status()
	}
	if statusChangesMask & C.DDS_LIVELINESS_LOST_STATUS != 0 {
		C.DDS_DataWriter_get_liveliness_lost_status()
	}
	if statusChangesMask & C.DDS_PUBLICATION_MATCHED_STATUS != 0 {
		C.DDS_DataWriter_get_publication_matched_status()
	}
	if statusChangesMask & C.DDS_RELIABLE_WRITER_CACHE_CHANGED_STATUS != 0 {
		C.DDS_DataWriter_get_reliable_writer_cache_changed_status()
	}
	if statusChangesMask & C.DDS_RELIABLE_READER_ACTIVITY_CHANGED_STATUS != 0 {
		C.DDS_DataWriter_get_reliable_reader_activity_changed_status()
	}
	if statusChangesMask & C.DDS_DATA_WRITER_CACHE_STATUS != 0 {
		C.DDS_DataWriter_get_datawriter_cache_status()
	}
	if statusChangesMask & C.DDS_DATA_WRITER_PROTOCOL_STATUS != 0 {
		C.DDS_DataWriter_get_datawriter_protocol_status()
	}
*/
	if statusChangesMask & C.DDS_REQUESTED_DEADLINE_MISSED_STATUS != 0 {
		var status C.struct_DDS_RequestedDeadlineMissedStatus
		C.DDS_DataReader_get_requested_deadline_missed_status(dr.dr, &status)
		//log.Println("DDS_REQUESTED_DEADLINE_MISSED_STATUS")
	}
	if statusChangesMask & C.DDS_REQUESTED_INCOMPATIBLE_QOS_STATUS != 0 {
		var status C.struct_DDS_RequestedIncompatibleQosStatus
		C.DDS_DataReader_get_requested_incompatible_qos_status(dr.dr, &status)
		//log.Println("DDS_REQUESTED_INCOMPATIBLE_QOS_STATUS")
	}
	if statusChangesMask & C.DDS_SAMPLE_LOST_STATUS != 0 {
		var status C.struct_DDS_SampleLostStatus
		C.DDS_DataReader_get_sample_lost_status(dr.dr, &status)
		//log.Println("DDS_SAMPLE_LOST_STATUS")
	}
	if statusChangesMask & C.DDS_SAMPLE_REJECTED_STATUS != 0 {
		var status C.struct_DDS_SampleRejectedStatus
		C.DDS_DataReader_get_sample_rejected_status(dr.dr, &status)
		//log.Println("DDS_SAMPLE_REJECTED_STATUS")
	}
	if statusChangesMask & C.DDS_LIVELINESS_CHANGED_STATUS != 0 {
		var status C.struct_DDS_LivelinessChangedStatus
		C.DDS_DataReader_get_liveliness_changed_status(dr.dr, &status)
		//log.Println("DDS_LIVELINESS_CHANGED_STATUS")
	}
	if statusChangesMask & C.DDS_SUBSCRIPTION_MATCHED_STATUS != 0 {
		var status C.struct_DDS_SubscriptionMatchedStatus
		C.DDS_DataReader_get_subscription_matched_status(dr.dr, &status)
		//log.Println("DDS_SUBSCRIPTION_MATCHED_STATUS")
	}
	if statusChangesMask & C.DDS_DATA_READER_CACHE_STATUS != 0 {
		var status C.struct_DDS_DataReaderCacheStatus
		C.DDS_DataReader_get_datareader_cache_status(dr.dr, &status)
		//log.Println("DDS_DATA_READER_CACHE_STATUS")
	}
	if statusChangesMask & C.DDS_DATA_READER_PROTOCOL_STATUS != 0 {
		var status C.struct_DDS_DataReaderProtocolStatus
		C.DDS_DataReader_get_datareader_protocol_status(dr.dr, &status)
		//log.Println("DDS_DATA_READER_PROTOCOL_STATUS")
	}

	// Note: DATA_WRITER_APPLICATION_ACKNOWLEDGMENT_STATUS does not trigger a StatusCondition
	// Note: DATA_WRITER_INSTANCE_REPLACED_STATUS has no get_xxx_status() function. Not implemented?
}

// GetUnsafe returns a pointer to the data reader as an unsafe pointer.
// C types cannot be used in other packages, so directly referencing
// DataReader.dr won't work outside rtiddsgo, in particular in whatever
// package the generated type code reside.
func (dr DataReader) GetUnsafe() unsafe.Pointer {
	return unsafe.Pointer(dr.dr)
}
