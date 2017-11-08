package netdicom

import (
	"fmt"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomtag"
	"github.com/grailbio/go-dicom/dicomuid"
	"github.com/grailbio/go-netdicom/dimse"
	"v.io/x/lib/vlog"
)

// Helper function used by C-{STORE,GET,MOVE} to send a dataset using C-STORE
// over an already-established association.
func runCStoreOnAssociation(upcallCh chan upcallEvent, downcallCh chan stateEvent,
	cm *contextManager,
	messageID uint16,
	ds *dicom.DataSet) error {
	var getElement = func(tag dicomtag.Tag) (string, error) {
		elem, err := ds.FindElementByTag(tag)
		if err != nil {
			return "", fmt.Errorf("C-STORE data lacks %s: %v", tag.String(), err)
		}
		s, err := elem.GetString()
		if err != nil {
			return "", err
		}
		return s, nil
	}
	sopInstanceUID, err := getElement(dicomtag.MediaStorageSOPInstanceUID)
	if err != nil {
		return fmt.Errorf("C-STORE data lacks SOPInstanceUID: %v", err)
	}
	sopClassUID, err := getElement(dicomtag.MediaStorageSOPClassUID)
	if err != nil {
		return fmt.Errorf("C-STORE data lacks MediaStorageSOPClassUID: %v", err)
	}
	vlog.VI(1).Infof("DICOM abstractsyntax: %s, sopinstance: %s", dicomuid.UIDString(sopClassUID), sopInstanceUID)
	context, err := cm.lookupByAbstractSyntaxUID(sopClassUID)
	if err != nil {
		vlog.Errorf("C-STORE: sop class %v not found in context %v", sopClassUID, err)
		return err
	}
	vlog.VI(1).Infof("C-STORE: using transfersyntax %s to send sop class %s, instance %s",
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
		vlog.Errorf("C-STORE: body encoder failed: %v", err)
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
		vlog.Infof("Start reading resp w/ messageID:%v", messageID)
		event, ok := <-upcallCh
		if !ok {
			return fmt.Errorf("Connection closed while waiting for C-STORE response")
		}
		vlog.VI(1).Infof("C-STORE resp event: %v", event.command)
		doassert(event.eventType == upcallEventData)
		doassert(event.command != nil)
		resp, ok := event.command.(*dimse.CStoreRsp)
		doassert(ok) // TODO(saito)
		if resp.Status.Status != 0 {
			return fmt.Errorf("C_STORE failed: %v", resp.String())
		}
		return nil
	}
}
