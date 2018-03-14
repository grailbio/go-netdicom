package netdicom

import (
	"fmt"
	"log"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomlog"
	"github.com/grailbio/go-dicom/dicomtag"
	"github.com/grailbio/go-dicom/dicomuid"
	"github.com/grailbio/go-netdicom/dimse"
)

// Helper function used by C-{STORE,GET,MOVE} to send a dataset using C-STORE
// over an already-established association.
func runCStoreOnAssociation(upcallCh chan upcallEvent, downcallCh chan stateEvent,
	cm *contextManager,
	messageID dimse.MessageID,
	ds *dicom.DataSet) error {
	var getElement = func(tag dicomtag.Tag) (string, error) {
		elem, err := ds.FindElementByTag(tag)
		if err != nil {
			return "", fmt.Errorf("dicom.cstore: data lacks %s: %v", tag.String(), err)
		}
		s, err := elem.GetString()
		if err != nil {
			return "", err
		}
		return s, nil
	}
	sopInstanceUID, err := getElement(dicomtag.MediaStorageSOPInstanceUID)
	if err != nil {
		return fmt.Errorf("dicom.cstore: data lacks SOPInstanceUID: %v", err)
	}
	sopClassUID, err := getElement(dicomtag.MediaStorageSOPClassUID)
	if err != nil {
		return fmt.Errorf("dicom.cstore: data lacks MediaStorageSOPClassUID: %v", err)
	}
	if dicomlog.Level >= 1 {
		log.Printf("dicom.cstore: DICOM abstractsyntax: %s, sopinstance: %s", dicomuid.UIDString(sopClassUID), sopInstanceUID)
	}
	context, err := cm.lookupByAbstractSyntaxUID(sopClassUID)
	if err != nil {
		log.Printf("dicom.cstore: sop class %v not found in context %v", sopClassUID, err)
		return err
	}
	if dicomlog.Level >= 1 {
		log.Printf("dicom.cstore: using transfersyntax %s to send sop class %s, instance %s",
			dicomuid.UIDString(context.transferSyntaxUID),
			dicomuid.UIDString(sopClassUID),
			sopInstanceUID)
	}
	bodyEncoder := dicomio.NewBytesEncoderWithTransferSyntax(context.transferSyntaxUID)
	for _, elem := range ds.Elements {
		if elem.Tag.Group == dicomtag.MetadataGroup {
			continue
		}
		dicom.WriteElement(bodyEncoder, elem)
	}
	if err := bodyEncoder.Error(); err != nil {
		log.Printf("dicom.cstore: body encoder failed: %v", err)
		return err
	}
	downcallCh <- stateEvent{
		event: evt09,
		dimsePayload: &stateEventDIMSEPayload{
			abstractSyntaxName: sopClassUID,
			command: &dimse.CStoreRq{
				AffectedSOPClassUID:    sopClassUID,
				MessageID:              messageID,
				CommandDataSetType:     dimse.CommandDataSetTypeNonNull,
				AffectedSOPInstanceUID: sopInstanceUID,
			},
			data: bodyEncoder.Bytes(),
		},
	}
	for {
		log.Printf("dicom.cstore: Start reading resp w/ messageID:%v", messageID)
		event, ok := <-upcallCh
		if !ok {
			return fmt.Errorf("dicom.cstore: Connection closed while waiting for C-STORE response")
		}
		if dicomlog.Level >= 1 {
			log.Printf("dicom.cstore: resp event: %v", event.command)
		}
		doassert(event.eventType == upcallEventData)
		doassert(event.command != nil)
		resp, ok := event.command.(*dimse.CStoreRsp)
		doassert(ok) // TODO(saito)
		if resp.Status.Status != 0 {
			return fmt.Errorf("dicom.cstore: failed: %v", resp.String())
		}
		return nil
	}
}
