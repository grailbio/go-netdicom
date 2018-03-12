package dimse

// Code generated from generate_dimse_messages.py. DO NOT EDIT.

import (
	"fmt"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomtag"
)

type CStoreRq struct {
	AffectedSOPClassUID                  string
	MessageID                            MessageID
	Priority                             uint16
	CommandDataSetType                   uint16
	AffectedSOPInstanceUID               string
	MoveOriginatorApplicationEntityTitle string
	MoveOriginatorMessageID              MessageID
	Extra                                []*dicom.Element // Unparsed elements
}

func (v *CStoreRq) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(1)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageID, v.MessageID))
	elems = append(elems, newElement(dicomtag.Priority, v.Priority))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, newElement(dicomtag.AffectedSOPInstanceUID, v.AffectedSOPInstanceUID))
	if v.MoveOriginatorApplicationEntityTitle != "" {
		elems = append(elems, newElement(dicomtag.MoveOriginatorApplicationEntityTitle, v.MoveOriginatorApplicationEntityTitle))
	}
	if v.MoveOriginatorMessageID != 0 {
		elems = append(elems, newElement(dicomtag.MoveOriginatorMessageID, v.MoveOriginatorMessageID))
	}
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CStoreRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CStoreRq) CommandField() int {
	return 1
}

func (v *CStoreRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CStoreRq) GetStatus() *Status {
	return nil
}

func (v *CStoreRq) String() string {
	return fmt.Sprintf("CStoreRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v MoveOriginatorApplicationEntityTitle:%v MoveOriginatorMessageID:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.MoveOriginatorApplicationEntityTitle, v.MoveOriginatorMessageID)
}

func decodeCStoreRq(d *messageDecoder) *CStoreRq {
	v := &CStoreRq{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicomtag.MessageID, requiredElement)
	v.Priority = d.getUInt16(dicomtag.Priority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.AffectedSOPInstanceUID = d.getString(dicomtag.AffectedSOPInstanceUID, requiredElement)
	v.MoveOriginatorApplicationEntityTitle = d.getString(dicomtag.MoveOriginatorApplicationEntityTitle, optionalElement)
	v.MoveOriginatorMessageID = d.getUInt16(dicomtag.MoveOriginatorMessageID, optionalElement)
	v.Extra = d.unparsedElements()
	return v
}

type CStoreRsp struct {
	AffectedSOPClassUID       string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType        uint16
	AffectedSOPInstanceUID    string
	Status                    Status
	Extra                     []*dicom.Element // Unparsed elements
}

func (v *CStoreRsp) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32769)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, newElement(dicomtag.AffectedSOPInstanceUID, v.AffectedSOPInstanceUID))
	elems = append(elems, newStatusElements(v.Status)...)
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CStoreRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CStoreRsp) CommandField() int {
	return 32769
}

func (v *CStoreRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CStoreRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CStoreRsp) String() string {
	return fmt.Sprintf("CStoreRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.Status)
}

func decodeCStoreRsp(d *messageDecoder) *CStoreRsp {
	v := &CStoreRsp{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicomtag.MessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.AffectedSOPInstanceUID = d.getString(dicomtag.AffectedSOPInstanceUID, requiredElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}

type CFindRq struct {
	AffectedSOPClassUID string
	MessageID           MessageID
	Priority            uint16
	CommandDataSetType  uint16
	Extra               []*dicom.Element // Unparsed elements
}

func (v *CFindRq) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageID, v.MessageID))
	elems = append(elems, newElement(dicomtag.Priority, v.Priority))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CFindRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CFindRq) CommandField() int {
	return 32
}

func (v *CFindRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CFindRq) GetStatus() *Status {
	return nil
}

func (v *CFindRq) String() string {
	return fmt.Sprintf("CFindRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCFindRq(d *messageDecoder) *CFindRq {
	v := &CFindRq{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicomtag.MessageID, requiredElement)
	v.Priority = d.getUInt16(dicomtag.Priority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}

type CFindRsp struct {
	AffectedSOPClassUID       string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType        uint16
	Status                    Status
	Extra                     []*dicom.Element // Unparsed elements
}

func (v *CFindRsp) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32800)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, newStatusElements(v.Status)...)
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CFindRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CFindRsp) CommandField() int {
	return 32800
}

func (v *CFindRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CFindRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CFindRsp) String() string {
	return fmt.Sprintf("CFindRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCFindRsp(d *messageDecoder) *CFindRsp {
	v := &CFindRsp{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicomtag.MessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}

type CGetRq struct {
	AffectedSOPClassUID string
	MessageID           MessageID
	Priority            uint16
	CommandDataSetType  uint16
	Extra               []*dicom.Element // Unparsed elements
}

func (v *CGetRq) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(16)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageID, v.MessageID))
	elems = append(elems, newElement(dicomtag.Priority, v.Priority))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CGetRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CGetRq) CommandField() int {
	return 16
}

func (v *CGetRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CGetRq) GetStatus() *Status {
	return nil
}

func (v *CGetRq) String() string {
	return fmt.Sprintf("CGetRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCGetRq(d *messageDecoder) *CGetRq {
	v := &CGetRq{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicomtag.MessageID, requiredElement)
	v.Priority = d.getUInt16(dicomtag.Priority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}

type CGetRsp struct {
	AffectedSOPClassUID            string
	MessageIDBeingRespondedTo      MessageID
	CommandDataSetType             uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations    uint16
	NumberOfWarningSuboperations   uint16
	Status                         Status
	Extra                          []*dicom.Element // Unparsed elements
}

func (v *CGetRsp) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32784)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	if v.NumberOfRemainingSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations))
	}
	if v.NumberOfCompletedSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations))
	}
	if v.NumberOfFailedSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfFailedSuboperations, v.NumberOfFailedSuboperations))
	}
	if v.NumberOfWarningSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfWarningSuboperations, v.NumberOfWarningSuboperations))
	}
	elems = append(elems, newStatusElements(v.Status)...)
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CGetRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CGetRsp) CommandField() int {
	return 32784
}

func (v *CGetRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CGetRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CGetRsp) String() string {
	return fmt.Sprintf("CGetRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCGetRsp(d *messageDecoder) *CGetRsp {
	v := &CGetRsp{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicomtag.MessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.NumberOfRemainingSuboperations = d.getUInt16(dicomtag.NumberOfRemainingSuboperations, optionalElement)
	v.NumberOfCompletedSuboperations = d.getUInt16(dicomtag.NumberOfCompletedSuboperations, optionalElement)
	v.NumberOfFailedSuboperations = d.getUInt16(dicomtag.NumberOfFailedSuboperations, optionalElement)
	v.NumberOfWarningSuboperations = d.getUInt16(dicomtag.NumberOfWarningSuboperations, optionalElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}

type CMoveRq struct {
	AffectedSOPClassUID string
	MessageID           MessageID
	Priority            uint16
	MoveDestination     string
	CommandDataSetType  uint16
	Extra               []*dicom.Element // Unparsed elements
}

func (v *CMoveRq) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(33)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageID, v.MessageID))
	elems = append(elems, newElement(dicomtag.Priority, v.Priority))
	elems = append(elems, newElement(dicomtag.MoveDestination, v.MoveDestination))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CMoveRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CMoveRq) CommandField() int {
	return 33
}

func (v *CMoveRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CMoveRq) GetStatus() *Status {
	return nil
}

func (v *CMoveRq) String() string {
	return fmt.Sprintf("CMoveRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v MoveDestination:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.MoveDestination, v.CommandDataSetType)
}

func decodeCMoveRq(d *messageDecoder) *CMoveRq {
	v := &CMoveRq{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicomtag.MessageID, requiredElement)
	v.Priority = d.getUInt16(dicomtag.Priority, requiredElement)
	v.MoveDestination = d.getString(dicomtag.MoveDestination, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}

type CMoveRsp struct {
	AffectedSOPClassUID            string
	MessageIDBeingRespondedTo      MessageID
	CommandDataSetType             uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations    uint16
	NumberOfWarningSuboperations   uint16
	Status                         Status
	Extra                          []*dicom.Element // Unparsed elements
}

func (v *CMoveRsp) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32801)))
	elems = append(elems, newElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID))
	elems = append(elems, newElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	if v.NumberOfRemainingSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations))
	}
	if v.NumberOfCompletedSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations))
	}
	if v.NumberOfFailedSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfFailedSuboperations, v.NumberOfFailedSuboperations))
	}
	if v.NumberOfWarningSuboperations != 0 {
		elems = append(elems, newElement(dicomtag.NumberOfWarningSuboperations, v.NumberOfWarningSuboperations))
	}
	elems = append(elems, newStatusElements(v.Status)...)
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CMoveRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CMoveRsp) CommandField() int {
	return 32801
}

func (v *CMoveRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CMoveRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CMoveRsp) String() string {
	return fmt.Sprintf("CMoveRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCMoveRsp(d *messageDecoder) *CMoveRsp {
	v := &CMoveRsp{}
	v.AffectedSOPClassUID = d.getString(dicomtag.AffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicomtag.MessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.NumberOfRemainingSuboperations = d.getUInt16(dicomtag.NumberOfRemainingSuboperations, optionalElement)
	v.NumberOfCompletedSuboperations = d.getUInt16(dicomtag.NumberOfCompletedSuboperations, optionalElement)
	v.NumberOfFailedSuboperations = d.getUInt16(dicomtag.NumberOfFailedSuboperations, optionalElement)
	v.NumberOfWarningSuboperations = d.getUInt16(dicomtag.NumberOfWarningSuboperations, optionalElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}

type CEchoRq struct {
	MessageID          MessageID
	CommandDataSetType uint16
	Extra              []*dicom.Element // Unparsed elements
}

func (v *CEchoRq) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(48)))
	elems = append(elems, newElement(dicomtag.MessageID, v.MessageID))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CEchoRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CEchoRq) CommandField() int {
	return 48
}

func (v *CEchoRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CEchoRq) GetStatus() *Status {
	return nil
}

func (v *CEchoRq) String() string {
	return fmt.Sprintf("CEchoRq{MessageID:%v CommandDataSetType:%v}}", v.MessageID, v.CommandDataSetType)
}

func decodeCEchoRq(d *messageDecoder) *CEchoRq {
	v := &CEchoRq{}
	v.MessageID = d.getUInt16(dicomtag.MessageID, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}

type CEchoRsp struct {
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType        uint16
	Status                    Status
	Extra                     []*dicom.Element // Unparsed elements
}

func (v *CEchoRsp) Encode(e *dicomio.Encoder) {
	elems := []*dicom.Element{}
	elems = append(elems, newElement(dicomtag.CommandField, uint16(32816)))
	elems = append(elems, newElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo))
	elems = append(elems, newElement(dicomtag.CommandDataSetType, v.CommandDataSetType))
	elems = append(elems, newStatusElements(v.Status)...)
	elems = append(elems, v.Extra...)
	encodeElements(e, elems)
}

func (v *CEchoRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CEchoRsp) CommandField() int {
	return 32816
}

func (v *CEchoRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CEchoRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CEchoRsp) String() string {
	return fmt.Sprintf("CEchoRsp{MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCEchoRsp(d *messageDecoder) *CEchoRsp {
	v := &CEchoRsp{}
	v.MessageIDBeingRespondedTo = d.getUInt16(dicomtag.MessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicomtag.CommandDataSetType, requiredElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}

const CommandFieldCStoreRq = 1
const CommandFieldCStoreRsp = 32769
const CommandFieldCFindRq = 32
const CommandFieldCFindRsp = 32800
const CommandFieldCGetRq = 16
const CommandFieldCGetRsp = 32784
const CommandFieldCMoveRq = 33
const CommandFieldCMoveRsp = 32801
const CommandFieldCEchoRq = 48
const CommandFieldCEchoRsp = 32816

func decodeMessageForType(d *messageDecoder, commandField uint16) Message {
	switch commandField {
	case 0x1:
		return decodeCStoreRq(d)
	case 0x8001:
		return decodeCStoreRsp(d)
	case 0x20:
		return decodeCFindRq(d)
	case 0x8020:
		return decodeCFindRsp(d)
	case 0x10:
		return decodeCGetRq(d)
	case 0x8010:
		return decodeCGetRsp(d)
	case 0x21:
		return decodeCMoveRq(d)
	case 0x8021:
		return decodeCMoveRsp(d)
	case 0x30:
		return decodeCEchoRq(d)
	case 0x8030:
		return decodeCEchoRsp(d)
	default:
		d.setError(fmt.Errorf("Unknown DIMSE command 0x%x", commandField))
		return nil
	}
}
