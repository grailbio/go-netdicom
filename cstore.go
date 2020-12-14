package netdicom

import (
	"fmt"

	"github.com/BTsykaniuk/go-netdicom/dimse"
	"github.com/apaladiychuk/go-dicom"
	"github.com/apaladiychuk/go-dicom/dicomio"
	"github.com/apaladiychuk/go-dicom/dicomlog"
	"github.com/apaladiychuk/go-dicom/dicomtag"
	"github.com/apaladiychuk/go-dicom/dicomuid"
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
	dicomlog.Vprintf(1, "dicom.cstore(%s): DICOM abstractsyntax: %s, sopinstance: %s", cm.label, dicomuid.UIDString(sopClassUID), sopInstanceUID)
	context, err := cm.lookupByAbstractSyntaxUID(sopClassUID)
	if err != nil {
		dicomlog.Vprintf(0, "dicom.cstore(%s): sop class %v not found in context %v", cm.label, sopClassUID, err)
		return err
	}
	dicomlog.Vprintf(1, "dicom.cstore(%s): using transfersyntax %s to send sop class %s, instance %s",
		cm.label,
		dicomuid.UIDString(context.transferSyntaxUID),
		dicomuid.UIDString(sopClassUID),
		sopInstanceUID)
	bodyEncoder := dicomio.NewBytesEncoderWithTransferSyntax(context.transferSyntaxUID)
	for _, elem := range ds.Elements {
		if elem.Tag.Group == dicomtag.MetadataGroup {
			continue
		}
		dicom.WriteElement(bodyEncoder, elem)
	}
	if err := bodyEncoder.Error(); err != nil {
		dicomlog.Vprintf(0, "dicom.cstore(%s): body encoder failed: %v", cm.label, err)
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
		dicomlog.Vprintf(0, "dicom.cstore(%s): Start reading resp w/ messageID:%v", cm.label, messageID)
		event, ok := <-upcallCh
		if !ok {
			return fmt.Errorf("dicom.cstore(%s): Connection closed while waiting for C-STORE response", cm.label)
		}
		dicomlog.Vprintf(1, "dicom.cstore(%s): resp event: %v", cm.label, event.command)
		doassert(event.eventType == upcallEventData)
		doassert(event.command != nil)
		resp, ok := event.command.(*dimse.CStoreRsp)
		doassert(ok) // TODO(saito)
		if resp.Status.Status != 0 {
			return fmt.Errorf("dicom.cstore(%s): failed: %v", cm.label, resp.String())
		}
		return nil
	}
}
