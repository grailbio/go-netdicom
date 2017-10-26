
package dimse

// Code generated from generate_dimse_messages.py. DO NOT EDIT.

import (
        "fmt"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
)
        
type CStoreRq struct {
	AffectedSOPClassUID string
	MessageID uint16
	Priority uint16
	CommandDataSetType uint16
	AffectedSOPInstanceUID string
	MoveOriginatorApplicationEntityTitle string
	MoveOriginatorMessageID uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CStoreRq) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(1))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageID, v.MessageID)
	encodeField(e, dicom.TagPriority, v.Priority)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	encodeField(e, dicom.TagAffectedSOPInstanceUID, v.AffectedSOPInstanceUID)
	if v.MoveOriginatorApplicationEntityTitle != "" {
		encodeField(e, dicom.TagMoveOriginatorApplicationEntityTitle, v.MoveOriginatorApplicationEntityTitle)
	}
	if v.MoveOriginatorMessageID != 0 {
		encodeField(e, dicom.TagMoveOriginatorMessageID, v.MoveOriginatorMessageID)
	}
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CStoreRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CStoreRq) CommandField() int {
	return 1
}

func (v* CStoreRq) GetMessageID() uint16 {
	return v.MessageID
}

func (v* CStoreRq) String() string {
	return fmt.Sprintf("CStoreRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v MoveOriginatorApplicationEntityTitle:%v MoveOriginatorMessageID:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.MoveOriginatorApplicationEntityTitle, v.MoveOriginatorMessageID)
}

func decodeCStoreRq(d *messageDecoder) *CStoreRq {
	v := &CStoreRq{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicom.TagMessageID, requiredElement)
	v.Priority = d.getUInt16(dicom.TagPriority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.AffectedSOPInstanceUID = d.getString(dicom.TagAffectedSOPInstanceUID, requiredElement)
	v.MoveOriginatorApplicationEntityTitle = d.getString(dicom.TagMoveOriginatorApplicationEntityTitle, optionalElement)
	v.MoveOriginatorMessageID = d.getUInt16(dicom.TagMoveOriginatorMessageID, optionalElement)
	v.Extra = d.unparsedElements()
	return v
}
type CStoreRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo uint16
	CommandDataSetType uint16
	AffectedSOPInstanceUID string
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CStoreRsp) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32769))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	encodeField(e, dicom.TagAffectedSOPInstanceUID, v.AffectedSOPInstanceUID)
	encodeStatus(e, v.Status)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CStoreRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CStoreRsp) CommandField() int {
	return 32769
}

func (v* CStoreRsp) GetMessageID() uint16 {
	return v.MessageIDBeingRespondedTo
}

func (v* CStoreRsp) String() string {
	return fmt.Sprintf("CStoreRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.Status)
}

func decodeCStoreRsp(d *messageDecoder) *CStoreRsp {
	v := &CStoreRsp{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicom.TagMessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.AffectedSOPInstanceUID = d.getString(dicom.TagAffectedSOPInstanceUID, requiredElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}
type CFindRq struct {
	AffectedSOPClassUID string
	MessageID uint16
	Priority uint16
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CFindRq) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageID, v.MessageID)
	encodeField(e, dicom.TagPriority, v.Priority)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CFindRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CFindRq) CommandField() int {
	return 32
}

func (v* CFindRq) GetMessageID() uint16 {
	return v.MessageID
}

func (v* CFindRq) String() string {
	return fmt.Sprintf("CFindRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCFindRq(d *messageDecoder) *CFindRq {
	v := &CFindRq{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicom.TagMessageID, requiredElement)
	v.Priority = d.getUInt16(dicom.TagPriority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}
type CFindRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo uint16
	CommandDataSetType uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CFindRsp) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32800))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	encodeStatus(e, v.Status)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CFindRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CFindRsp) CommandField() int {
	return 32800
}

func (v* CFindRsp) GetMessageID() uint16 {
	return v.MessageIDBeingRespondedTo
}

func (v* CFindRsp) String() string {
	return fmt.Sprintf("CFindRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCFindRsp(d *messageDecoder) *CFindRsp {
	v := &CFindRsp{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicom.TagMessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}
type CGetRq struct {
	AffectedSOPClassUID string
	MessageID uint16
	Priority uint16
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CGetRq) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(16))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageID, v.MessageID)
	encodeField(e, dicom.TagPriority, v.Priority)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CGetRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CGetRq) CommandField() int {
	return 16
}

func (v* CGetRq) GetMessageID() uint16 {
	return v.MessageID
}

func (v* CGetRq) String() string {
	return fmt.Sprintf("CGetRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCGetRq(d *messageDecoder) *CGetRq {
	v := &CGetRq{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicom.TagMessageID, requiredElement)
	v.Priority = d.getUInt16(dicom.TagPriority, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}
type CGetRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo uint16
	CommandDataSetType uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations uint16
	NumberOfWarningSuboperations uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CGetRsp) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32784))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	if v.NumberOfRemainingSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations)
	}
	if v.NumberOfCompletedSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations)
	}
	if v.NumberOfFailedSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfFailedSuboperations, v.NumberOfFailedSuboperations)
	}
	if v.NumberOfWarningSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfWarningSuboperations, v.NumberOfWarningSuboperations)
	}
	encodeStatus(e, v.Status)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CGetRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CGetRsp) CommandField() int {
	return 32784
}

func (v* CGetRsp) GetMessageID() uint16 {
	return v.MessageIDBeingRespondedTo
}

func (v* CGetRsp) String() string {
	return fmt.Sprintf("CGetRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCGetRsp(d *messageDecoder) *CGetRsp {
	v := &CGetRsp{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicom.TagMessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.NumberOfRemainingSuboperations = d.getUInt16(dicom.TagNumberOfRemainingSuboperations, optionalElement)
	v.NumberOfCompletedSuboperations = d.getUInt16(dicom.TagNumberOfCompletedSuboperations, optionalElement)
	v.NumberOfFailedSuboperations = d.getUInt16(dicom.TagNumberOfFailedSuboperations, optionalElement)
	v.NumberOfWarningSuboperations = d.getUInt16(dicom.TagNumberOfWarningSuboperations, optionalElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}
type CMoveRq struct {
	AffectedSOPClassUID string
	MessageID uint16
	Priority uint16
	MoveDestination string
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CMoveRq) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(33))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageID, v.MessageID)
	encodeField(e, dicom.TagPriority, v.Priority)
	encodeField(e, dicom.TagMoveDestination, v.MoveDestination)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CMoveRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CMoveRq) CommandField() int {
	return 33
}

func (v* CMoveRq) GetMessageID() uint16 {
	return v.MessageID
}

func (v* CMoveRq) String() string {
	return fmt.Sprintf("CMoveRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v MoveDestination:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.MoveDestination, v.CommandDataSetType)
}

func decodeCMoveRq(d *messageDecoder) *CMoveRq {
	v := &CMoveRq{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageID = d.getUInt16(dicom.TagMessageID, requiredElement)
	v.Priority = d.getUInt16(dicom.TagPriority, requiredElement)
	v.MoveDestination = d.getString(dicom.TagMoveDestination, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}
type CMoveRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo uint16
	CommandDataSetType uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations uint16
	NumberOfWarningSuboperations uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CMoveRsp) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32801))
	encodeField(e, dicom.TagAffectedSOPClassUID, v.AffectedSOPClassUID)
	encodeField(e, dicom.TagMessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	if v.NumberOfRemainingSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations)
	}
	if v.NumberOfCompletedSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations)
	}
	if v.NumberOfFailedSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfFailedSuboperations, v.NumberOfFailedSuboperations)
	}
	if v.NumberOfWarningSuboperations != 0 {
		encodeField(e, dicom.TagNumberOfWarningSuboperations, v.NumberOfWarningSuboperations)
	}
	encodeStatus(e, v.Status)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CMoveRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CMoveRsp) CommandField() int {
	return 32801
}

func (v* CMoveRsp) GetMessageID() uint16 {
	return v.MessageIDBeingRespondedTo
}

func (v* CMoveRsp) String() string {
	return fmt.Sprintf("CMoveRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCMoveRsp(d *messageDecoder) *CMoveRsp {
	v := &CMoveRsp{}
	v.AffectedSOPClassUID = d.getString(dicom.TagAffectedSOPClassUID, requiredElement)
	v.MessageIDBeingRespondedTo = d.getUInt16(dicom.TagMessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.NumberOfRemainingSuboperations = d.getUInt16(dicom.TagNumberOfRemainingSuboperations, optionalElement)
	v.NumberOfCompletedSuboperations = d.getUInt16(dicom.TagNumberOfCompletedSuboperations, optionalElement)
	v.NumberOfFailedSuboperations = d.getUInt16(dicom.TagNumberOfFailedSuboperations, optionalElement)
	v.NumberOfWarningSuboperations = d.getUInt16(dicom.TagNumberOfWarningSuboperations, optionalElement)
	v.Status = d.getStatus()
	v.Extra = d.unparsedElements()
	return v
}
type CEchoRq struct {
	MessageID uint16
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CEchoRq) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(48))
	encodeField(e, dicom.TagMessageID, v.MessageID)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CEchoRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CEchoRq) CommandField() int {
	return 48
}

func (v* CEchoRq) GetMessageID() uint16 {
	return v.MessageID
}

func (v* CEchoRq) String() string {
	return fmt.Sprintf("CEchoRq{MessageID:%v CommandDataSetType:%v}}", v.MessageID, v.CommandDataSetType)
}

func decodeCEchoRq(d *messageDecoder) *CEchoRq {
	v := &CEchoRq{}
	v.MessageID = d.getUInt16(dicom.TagMessageID, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
	v.Extra = d.unparsedElements()
	return v
}
type CEchoRsp struct {
	MessageIDBeingRespondedTo uint16
	CommandDataSetType uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v* CEchoRsp) Encode(e *dicomio.Encoder) {
	encodeField(e, dicom.TagCommandField, uint16(32816))
	encodeField(e, dicom.TagMessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
	encodeField(e, dicom.TagCommandDataSetType, v.CommandDataSetType)
	encodeStatus(e, v.Status)
	for _, elem := range v.Extra {
		dicom.WriteElement(e, elem)
	}
}

func (v* CEchoRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v* CEchoRsp) CommandField() int {
	return 32816
}

func (v* CEchoRsp) GetMessageID() uint16 {
	return v.MessageIDBeingRespondedTo
}

func (v* CEchoRsp) String() string {
	return fmt.Sprintf("CEchoRsp{MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCEchoRsp(d *messageDecoder) *CEchoRsp {
	v := &CEchoRsp{}
	v.MessageIDBeingRespondedTo = d.getUInt16(dicom.TagMessageIDBeingRespondedTo, requiredElement)
	v.CommandDataSetType = d.getUInt16(dicom.TagCommandDataSetType, requiredElement)
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
func decodeMessageForType(d* messageDecoder, commandField uint16) Message {
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
