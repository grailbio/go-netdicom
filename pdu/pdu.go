package pdu

//go:generate stringer -type AbortReasonType
//go:generate stringer -type PresentationContextResult
//go:generate stringer -type RejectReasonType
//go:generate stringer -type RejectResultType
//go:generate stringer -type SourceType
//go:generate stringer -type Type

// Implements message types defined in P3.8. It sits below the DIMSE layer.
//
// http://dicom.nema.org/medical/dicom/current/output/pdf/part08.pdf
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/grailbio/go-dicom/dicomio"
)

// PDU is the interface for DUL messages like A-ASSOCIATE-AC, P-DATA-TF.
type PDU interface {
	fmt.Stringer

	// WritePayload encodes the PDU payload. The "payload" here excludes the
	// first 6 bytes that are common to all PDU types - they are encoded in
	// EncodePDU separately.
	WritePayload(*dicomio.Encoder)
}

// Type defines type of the PDU packet.
type Type byte

const (
	TypeAAssociateRq Type = 1 // A_ASSOCIATE_RQ
	TypeAAssociateAc      = 2 // A_ASSOCIATE_AC
	TypeAAssociateRj      = 3 // A_ASSOCIATE_RJ
	TypePDataTf           = 4 // P_DATA_TF
	TypeAReleaseRq        = 5 // A_RELEASE_RQ
	TypeAReleaseRp        = 6 // A_RELEASE_RP
	TypeAAbort            = 7 // A_ABORT
)

// SubItem is the interface for DUL items, such as ApplicationContextItem and
// TransferSyntaxSubItem.
type SubItem interface {
	fmt.Stringer

	// Write serializes the item.
	Write(*dicomio.Encoder)
}

// Possible Type field values for SubItem.
const (
	ItemTypeApplicationContext           = 0x10
	ItemTypePresentationContextRequest   = 0x20
	ItemTypePresentationContextResponse  = 0x21
	ItemTypeAbstractSyntax               = 0x30
	ItemTypeTransferSyntax               = 0x40
	ItemTypeUserInformation              = 0x50
	ItemTypeUserInformationMaximumLength = 0x51
	ItemTypeImplementationClassUID       = 0x52
	ItemTypeAsynchronousOperationsWindow = 0x53
	ItemTypeRoleSelection                = 0x54
	ItemTypeImplementationVersionName    = 0x55
)

func decodeSubItem(d *dicomio.Decoder) SubItem {
	itemType := d.ReadByte()
	d.Skip(1)
	length := d.ReadUInt16()
	switch itemType {
	case ItemTypeApplicationContext:
		return decodeApplicationContextItem(d, length)
	case ItemTypeAbstractSyntax:
		return decodeAbstractSyntaxSubItem(d, length)
	case ItemTypeTransferSyntax:
		return decodeTransferSyntaxSubItem(d, length)
	case ItemTypePresentationContextRequest:
		return decodePresentationContextItem(d, itemType, length)
	case ItemTypePresentationContextResponse:
		return decodePresentationContextItem(d, itemType, length)
	case ItemTypeUserInformation:
		return decodeUserInformationItem(d, length)
	case ItemTypeUserInformationMaximumLength:
		return decodeUserInformationMaximumLengthItem(d, length)
	case ItemTypeImplementationClassUID:
		return decodeImplementationClassUIDSubItem(d, length)
	case ItemTypeAsynchronousOperationsWindow:
		return decodeAsynchronousOperationsWindowSubItem(d, length)
	case ItemTypeRoleSelection:
		return decodeRoleSelectionSubItem(d, length)
	case ItemTypeImplementationVersionName:
		return decodeImplementationVersionNameSubItem(d, length)
	default:
		d.SetError(fmt.Errorf("Unknown item type: 0x%x", itemType))
		return nil
	}
}

func encodeSubItemHeader(e *dicomio.Encoder, itemType byte, length uint16) {
	e.WriteByte(itemType)
	e.WriteZeros(1)
	e.WriteUInt16(length)
}

// P3.8 9.3.2.3
type UserInformationItem struct {
	Items []SubItem // P3.8, Annex D.
}

func (v *UserInformationItem) Write(e *dicomio.Encoder) {
	itemEncoder := dicomio.NewBytesEncoder(binary.BigEndian, dicomio.UnknownVR)
	for _, s := range v.Items {
		s.Write(itemEncoder)
	}
	if err := itemEncoder.Error(); err != nil {
		e.SetError(err)
		return
	}
	itemBytes := itemEncoder.Bytes()
	encodeSubItemHeader(e, ItemTypeUserInformation, uint16(len(itemBytes)))
	e.WriteBytes(itemBytes)
}

func decodeUserInformationItem(d *dicomio.Decoder, length uint16) *UserInformationItem {
	v := &UserInformationItem{}
	d.PushLimit(int64(length))
	defer d.PopLimit()
	for !d.EOF() {
		item := decodeSubItem(d)
		if d.Error() != nil {
			break
		}
		v.Items = append(v.Items, item)
	}
	return v
}

func (v *UserInformationItem) String() string {
	return fmt.Sprintf("UserInformationItem{items: %s}",
		subItemListString(v.Items))
}

// P3.8 D.1
type UserInformationMaximumLengthItem struct {
	MaximumLengthReceived uint32
}

func (v *UserInformationMaximumLengthItem) Write(e *dicomio.Encoder) {
	encodeSubItemHeader(e, ItemTypeUserInformationMaximumLength, 4)
	e.WriteUInt32(v.MaximumLengthReceived)
}

func decodeUserInformationMaximumLengthItem(d *dicomio.Decoder, length uint16) *UserInformationMaximumLengthItem {
	if length != 4 {
		d.SetError(fmt.Errorf("UserInformationMaximumLengthItem must be 4 bytes, but found %dB", length))
	}
	return &UserInformationMaximumLengthItem{MaximumLengthReceived: d.ReadUInt32()}
}

func (v *UserInformationMaximumLengthItem) String() string {
	return fmt.Sprintf("UserInformationMaximumlengthItem{%d}",
		v.MaximumLengthReceived)
}

// PS3.7 Annex D.3.3.2.1
type ImplementationClassUIDSubItem subItemWithName

func decodeImplementationClassUIDSubItem(d *dicomio.Decoder, length uint16) *ImplementationClassUIDSubItem {
	return &ImplementationClassUIDSubItem{Name: decodeSubItemWithName(d, length)}
}

func (v *ImplementationClassUIDSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemWithName(e, ItemTypeImplementationClassUID, v.Name)
}

func (v *ImplementationClassUIDSubItem) String() string {
	return fmt.Sprintf("ImplementationClassUID{name: \"%s\"}", v.Name)
}

// PS3.7 Annex D.3.3.3.1
type AsynchronousOperationsWindowSubItem struct {
	MaxOpsInvoked   uint16
	MaxOpsPerformed uint16
}

func decodeAsynchronousOperationsWindowSubItem(d *dicomio.Decoder, length uint16) *AsynchronousOperationsWindowSubItem {
	return &AsynchronousOperationsWindowSubItem{
		MaxOpsInvoked:   d.ReadUInt16(),
		MaxOpsPerformed: d.ReadUInt16(),
	}
}

func (v *AsynchronousOperationsWindowSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemHeader(e, ItemTypeAsynchronousOperationsWindow, 2*2)
	e.WriteUInt16(v.MaxOpsInvoked)
	e.WriteUInt16(v.MaxOpsPerformed)
}

func (v *AsynchronousOperationsWindowSubItem) String() string {
	return fmt.Sprintf("AsynchronousOpsWindow{invoked: %d performed: %d}",
		v.MaxOpsInvoked, v.MaxOpsPerformed)
}

// PS3.7 Annex D.3.3.4
type RoleSelectionSubItem struct {
	SOPClassUID string
	SCURole     byte
	SCPRole     byte
}

func decodeRoleSelectionSubItem(d *dicomio.Decoder, length uint16) *RoleSelectionSubItem {
	uidLen := d.ReadUInt16()
	return &RoleSelectionSubItem{
		SOPClassUID: d.ReadString(int(uidLen)),
		SCURole:     d.ReadByte(),
		SCPRole:     d.ReadByte(),
	}
}

func (v *RoleSelectionSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemHeader(e, ItemTypeRoleSelection, uint16(2+len(v.SOPClassUID)+1*2))
	e.WriteUInt16(uint16(len(v.SOPClassUID)))
	e.WriteString(v.SOPClassUID)
	e.WriteByte(v.SCURole)
	e.WriteByte(v.SCPRole)
}

func (v *RoleSelectionSubItem) String() string {
	return fmt.Sprintf("RoleSelection{sopclassuid: %v, scu: %v, scp: %v}", v.SOPClassUID, v.SCURole, v.SCPRole)
}

// PS3.7 Annex D.3.3.2.3
type ImplementationVersionNameSubItem subItemWithName

func decodeImplementationVersionNameSubItem(d *dicomio.Decoder, length uint16) *ImplementationVersionNameSubItem {
	return &ImplementationVersionNameSubItem{Name: decodeSubItemWithName(d, length)}
}

func (v *ImplementationVersionNameSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemWithName(e, ItemTypeImplementationVersionName, v.Name)
}

func (v *ImplementationVersionNameSubItem) String() string {
	return fmt.Sprintf("ImplementationVersionName{name: \"%s\"}", v.Name)
}

// Container for subitems that this package doesnt' support
type SubItemUnsupported struct {
	Type byte
	Data []byte
}

func (item *SubItemUnsupported) Write(e *dicomio.Encoder) {
	encodeSubItemHeader(e, item.Type, uint16(len(item.Data)))
	// TODO: handle unicode properly
	e.WriteBytes(item.Data)
}

func (item *SubItemUnsupported) String() string {
	return fmt.Sprintf("SubitemUnsupported{type: 0x%0x data: %dbytes}",
		item.Type, len(item.Data))
}

type subItemWithName struct {
	// Type byte
	Name string
}

func encodeSubItemWithName(e *dicomio.Encoder, itemType byte, name string) {
	encodeSubItemHeader(e, itemType, uint16(len(name)))
	// TODO: handle unicode properly
	e.WriteBytes([]byte(name))
}

func decodeSubItemWithName(d *dicomio.Decoder, length uint16) string {
	return d.ReadString(int(length))
}

type ApplicationContextItem subItemWithName

// The app context for DICOM. The first item in the A-ASSOCIATE-RQ
const DICOMApplicationContextItemName = "1.2.840.10008.3.1.1.1"

func decodeApplicationContextItem(d *dicomio.Decoder, length uint16) *ApplicationContextItem {
	return &ApplicationContextItem{Name: decodeSubItemWithName(d, length)}
}

func (v *ApplicationContextItem) Write(e *dicomio.Encoder) {
	encodeSubItemWithName(e, ItemTypeApplicationContext, v.Name)
}

func (v *ApplicationContextItem) String() string {
	return fmt.Sprintf("ApplicationContext{name: \"%s\"}", v.Name)
}

type AbstractSyntaxSubItem subItemWithName

func decodeAbstractSyntaxSubItem(d *dicomio.Decoder, length uint16) *AbstractSyntaxSubItem {
	return &AbstractSyntaxSubItem{Name: decodeSubItemWithName(d, length)}
}

func (v *AbstractSyntaxSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemWithName(e, ItemTypeAbstractSyntax, v.Name)
}

func (v *AbstractSyntaxSubItem) String() string {
	return fmt.Sprintf("AbstractSyntax{name: \"%s\"}", v.Name)
}

type TransferSyntaxSubItem subItemWithName

func decodeTransferSyntaxSubItem(d *dicomio.Decoder, length uint16) *TransferSyntaxSubItem {
	return &TransferSyntaxSubItem{Name: decodeSubItemWithName(d, length)}
}

func (v *TransferSyntaxSubItem) Write(e *dicomio.Encoder) {
	encodeSubItemWithName(e, ItemTypeTransferSyntax, v.Name)
}

func (v *TransferSyntaxSubItem) String() string {
	return fmt.Sprintf("TransferSyntax{name: \"%s\"}", v.Name)
}

// Result of abstractsyntax/transfersyntax handshake during A-ACCEPT.  P3.8,
// 90.3.3.2, table 9-18.
type PresentationContextResult byte

const (
	PresentationContextAccepted                                    PresentationContextResult = 0
	PresentationContextUserRejection                               PresentationContextResult = 1
	PresentationContextProviderRejectionNoReason                   PresentationContextResult = 2
	PresentationContextProviderRejectionAbstractSyntaxNotSupported PresentationContextResult = 3
	PresentationContextProviderRejectionTransferSyntaxNotSupported PresentationContextResult = 4
)

// P3.8 9.3.2.2, 9.3.3.2
type PresentationContextItem struct {
	Type      byte // ItemTypePresentationContext*
	ContextID byte
	// 1 byte reserved

	// Result is meaningful iff Type=0x21, zero else.
	Result PresentationContextResult

	// 1 byte reserved
	Items []SubItem // List of {Abstract,Transfer}SyntaxSubItem
}

func decodePresentationContextItem(d *dicomio.Decoder, itemType byte, length uint16) *PresentationContextItem {
	v := &PresentationContextItem{Type: itemType}
	d.PushLimit(int64(length))
	defer d.PopLimit()
	v.ContextID = d.ReadByte()
	d.Skip(1)
	v.Result = PresentationContextResult(d.ReadByte())
	d.Skip(1)
	for !d.EOF() {
		item := decodeSubItem(d)
		if d.Error() != nil {
			break
		}
		v.Items = append(v.Items, item)
	}
	if v.ContextID%2 != 1 {
		d.SetError(fmt.Errorf("PresentationContextItem ID must be odd, but found %x", v.ContextID))
	}
	return v
}

func (v *PresentationContextItem) Write(e *dicomio.Encoder) {
	if v.Type != ItemTypePresentationContextRequest &&
		v.Type != ItemTypePresentationContextResponse {
		panic(*v)
	}

	itemEncoder := dicomio.NewBytesEncoder(binary.BigEndian, dicomio.UnknownVR)
	for _, s := range v.Items {
		s.Write(itemEncoder)
	}
	if err := itemEncoder.Error(); err != nil {
		e.SetError(err)
		return
	}
	itemBytes := itemEncoder.Bytes()
	encodeSubItemHeader(e, v.Type, uint16(4+len(itemBytes)))
	e.WriteByte(v.ContextID)
	e.WriteZeros(3)
	e.WriteBytes(itemBytes)
}

func (v *PresentationContextItem) String() string {
	itemType := "rq"
	if v.Type == ItemTypePresentationContextResponse {
		itemType = "ac"
	}
	return fmt.Sprintf("PresentationContext%s{id: %d result: %d, items:%s}",
		itemType, v.ContextID, v.Result, subItemListString(v.Items))
}

// P3.8 9.3.2.2.1 & 9.3.2.2.2
type PresentationDataValueItem struct {
	// Length: 2 + len(Value)
	ContextID byte

	// P3.8, E.2: the following two fields encode a single byte.
	Command bool // Bit 7 (LSB): 1 means command 0 means data
	Last    bool // Bit 6: 1 means last fragment. 0 means not last fragment.

	// Payload, either command or data
	Value []byte
}

func ReadPresentationDataValueItem(d *dicomio.Decoder) PresentationDataValueItem {
	item := PresentationDataValueItem{}
	length := d.ReadUInt32()
	item.ContextID = d.ReadByte()
	header := d.ReadByte()
	item.Command = (header&1 != 0)
	item.Last = (header&2 != 0)
	item.Value = d.ReadBytes(int(length - 2)) // remove contextID and header
	return item
}

func (v *PresentationDataValueItem) Write(e *dicomio.Encoder) {
	var header byte
	if v.Command {
		header |= 1
	}
	if v.Last {
		header |= 2
	}
	e.WriteUInt32(uint32(2 + len(v.Value)))
	e.WriteByte(v.ContextID)
	e.WriteByte(header)
	e.WriteBytes(v.Value)
}

func (v *PresentationDataValueItem) String() string {
	return fmt.Sprintf("PresentationDataValue{context: %d, cmd:%v last:%v value: %d bytes}", v.ContextID, v.Command, v.Last, len(v.Value))
}

// EncodePDU serializes "pdu" into []byte.
func EncodePDU(pdu PDU) ([]byte, error) {
	var pduType Type
	switch n := pdu.(type) {
	case *AAssociate:
		pduType = n.Type
	case *AAssociateRj:
		pduType = TypeAAssociateRj
	case *PDataTf:
		pduType = TypePDataTf
	case *AReleaseRq:
		pduType = TypeAReleaseRq
	case *AReleaseRp:
		pduType = TypeAReleaseRp
	case *AAbort:
		pduType = TypeAAbort
	default:
		panic(fmt.Sprintf("Unknown PDU %v", pdu))
	}
	e := dicomio.NewBytesEncoder(binary.BigEndian, dicomio.UnknownVR)
	pdu.WritePayload(e)
	if err := e.Error(); err != nil {
		return nil, err
	}
	payload := e.Bytes()
	// Reserve the header bytes. It will be filled in Finish.
	var header [6]byte // First 6 bytes of buf.
	header[0] = byte(pduType)
	header[1] = 0 // Reserved.
	binary.BigEndian.PutUint32(header[2:6], uint32(len(payload)))
	return append(header[:], payload...), nil
}

// EncodePDU reads a "pdu" from a stream. maxPDUSize defines the maximum
// possible PDU size, in bytes, accepted by the caller.
func ReadPDU(in io.Reader, maxPDUSize int) (PDU, error) {
	var pduType Type
	var skip byte
	var length uint32
	err := binary.Read(in, binary.BigEndian, &pduType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(in, binary.BigEndian, &skip)
	if err != nil {
		return nil, err
	}
	err = binary.Read(in, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	if length >= uint32(maxPDUSize)*2 {
		// Avoid using too much memory. *2 is just an arbitrary slack.
		return nil, fmt.Errorf("Invalid length %d; it's much larger than max PDU size of %d", length, maxPDUSize)
	}
	d := dicomio.NewDecoder(
		&io.LimitedReader{R: in, N: int64(length)},
		binary.BigEndian,  // PDU is always big endian
		dicomio.UnknownVR) // irrelevant for PDU parsing
	var pdu PDU
	switch pduType {
	case TypeAAssociateRq:
		fallthrough
	case TypeAAssociateAc:
		pdu = decodeAAssociate(d, pduType)
	case TypeAAssociateRj:
		pdu = decodeAAssociateRj(d)
	case TypeAAbort:
		pdu = decodeAAbort(d)
	case TypePDataTf:
		pdu = decodePDataTf(d)
	case TypeAReleaseRq:
		pdu = decodeAReleaseRq(d)
	case TypeAReleaseRp:
		pdu = decodeAReleaseRp(d)
	}
	if pdu == nil {
		err := fmt.Errorf("ReadPDU: unknown message type %d", pduType)
		return nil, err
	}
	if err := d.Finish(); err != nil {
		return nil, err
	}
	return pdu, nil
}

type AReleaseRq struct {
}

func decodeAReleaseRq(d *dicomio.Decoder) *AReleaseRq {
	pdu := &AReleaseRq{}
	d.Skip(4)
	return pdu
}

func (pdu *AReleaseRq) WritePayload(e *dicomio.Encoder) {
	e.WriteZeros(4)
}

func (pdu *AReleaseRq) String() string {
	return fmt.Sprintf("A_RELEASE_RQ(%v)", *pdu)
}

type AReleaseRp struct {
}

func decodeAReleaseRp(d *dicomio.Decoder) *AReleaseRp {
	pdu := &AReleaseRp{}
	d.Skip(4)
	return pdu
}

func (pdu *AReleaseRp) WritePayload(e *dicomio.Encoder) {
	e.WriteZeros(4)
}

func (pdu *AReleaseRp) String() string {
	return fmt.Sprintf("A_RELEASE_RP(%v)", *pdu)
}

func subItemListString(items []SubItem) string {
	buf := bytes.Buffer{}
	buf.WriteString("[")
	for i, subitem := range items {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(subitem.String())
	}
	buf.WriteString("]")
	return buf.String()
}

const CurrentProtocolVersion uint16 = 1

// Defines A_ASSOCIATE_{RQ,AC}. P3.8 9.3.2 and 9.3.3
type AAssociate struct {
	Type            Type // One of {TypeA_Associate_RQ,TypeA_Associate_AC}
	ProtocolVersion uint16
	// Reserved uint16
	CalledAETitle  string // For .._AC, the value is copied from A_ASSOCIATE_RQ
	CallingAETitle string // For .._AC, the value is copied from A_ASSOCIATE_RQ
	Items          []SubItem
}

func decodeAAssociate(d *dicomio.Decoder, pduType Type) *AAssociate {
	pdu := &AAssociate{}
	pdu.Type = pduType
	pdu.ProtocolVersion = d.ReadUInt16()
	d.Skip(2) // Reserved
	pdu.CalledAETitle = d.ReadString(16)
	pdu.CallingAETitle = d.ReadString(16)
	d.Skip(8 * 4)
	for !d.EOF() {
		item := decodeSubItem(d)
		if d.Error() != nil {
			break
		}
		pdu.Items = append(pdu.Items, item)
	}
	if pdu.CalledAETitle == "" || pdu.CallingAETitle == "" {
		d.SetError(fmt.Errorf("A_ASSOCIATE.{Called,Calling}AETitle must not be empty, in %v", pdu.String()))
	}
	return pdu
}

func (pdu *AAssociate) WritePayload(e *dicomio.Encoder) {
	if pdu.Type == 0 || pdu.CalledAETitle == "" || pdu.CallingAETitle == "" {
		panic(*pdu)
	}
	e.WriteUInt16(pdu.ProtocolVersion)
	e.WriteZeros(2) // Reserved
	e.WriteString(fillString(pdu.CalledAETitle, 16))
	e.WriteString(fillString(pdu.CallingAETitle, 16))
	e.WriteZeros(8 * 4)
	for _, item := range pdu.Items {
		item.Write(e)
	}
}

func (pdu *AAssociate) String() string {
	name := "AC"
	if pdu.Type == TypeAAssociateRq {
		name = "RQ"
	}
	return fmt.Sprintf("A_ASSOCIATE_%s{version:%v called:'%v' calling:'%v' items:%s}",
		name, pdu.ProtocolVersion,
		pdu.CalledAETitle, pdu.CallingAETitle, subItemListString(pdu.Items))
}

// P3.8 9.3.4
type AAssociateRj struct {
	Result RejectResultType
	Source SourceType
	Reason RejectReasonType
}

// Possible values for AAssociateRj.Result
type RejectResultType byte

const (
	ResultRejectedPermanent RejectResultType = 1
	ResultRejectedTransient RejectResultType = 2
)

// Possible values for AAssociateRj.Reason
type RejectReasonType byte

const (
	RejectReasonNone                               RejectReasonType = 1
	RejectReasonApplicationContextNameNotSupported RejectReasonType = 2
	RejectReasonCallingAETitleNotRecognized        RejectReasonType = 3
	RejectReasonCalledAETitleNotRecognized         RejectReasonType = 7
)

// Possible values for AAssociateRj.Source
type SourceType byte

const (
	SourceULServiceUser                 SourceType = 1
	SourceULServiceProviderACSE         SourceType = 2
	SourceULServiceProviderPresentation SourceType = 3
)

func decodeAAssociateRj(d *dicomio.Decoder) *AAssociateRj {
	pdu := &AAssociateRj{}
	d.Skip(1) // reserved
	pdu.Result = RejectResultType(d.ReadByte())
	pdu.Source = SourceType(d.ReadByte())
	pdu.Reason = RejectReasonType(d.ReadByte())
	return pdu
}

func (pdu *AAssociateRj) WritePayload(e *dicomio.Encoder) {
	e.WriteZeros(1)
	e.WriteByte(byte(pdu.Result))
	e.WriteByte(byte(pdu.Source))
	e.WriteByte(byte(pdu.Reason))
}

func (pdu *AAssociateRj) String() string {
	return fmt.Sprintf("A_ASSOCIATE_RJ{result: %v, source: %v, reason: %v}", pdu.Result, pdu.Source, pdu.Reason)
}

type AbortReasonType byte

const (
	AbortReasonNotSpecified             AbortReasonType = 0
	AbortReasonUnexpectedPDU            AbortReasonType = 2
	AbortReasonUnrecognizedPDUParameter AbortReasonType = 3
	AbortReasonUnexpectedPDUParameter   AbortReasonType = 4
	AbortReasonInvalidPDUParameterValue AbortReasonType = 5
)

type AAbort struct {
	Source SourceType
	Reason AbortReasonType
}

func decodeAAbort(d *dicomio.Decoder) *AAbort {
	pdu := &AAbort{}
	d.Skip(2)
	pdu.Source = SourceType(d.ReadByte())
	pdu.Reason = AbortReasonType(d.ReadByte())
	return pdu
}

func (pdu *AAbort) WritePayload(e *dicomio.Encoder) {
	e.WriteZeros(2)
	e.WriteByte(byte(pdu.Source))
	e.WriteByte(byte(pdu.Reason))
}

func (pdu *AAbort) String() string {
	return fmt.Sprintf("A_ABORT{source:%v reason:%v}", pdu.Source, pdu.Reason)
}

type PDataTf struct {
	Items []PresentationDataValueItem
}

func decodePDataTf(d *dicomio.Decoder) *PDataTf {
	pdu := &PDataTf{}
	for !d.EOF() {
		item := ReadPresentationDataValueItem(d)
		if d.Error() != nil {
			break
		}
		pdu.Items = append(pdu.Items, item)
	}
	return pdu
}

func (pdu *PDataTf) WritePayload(e *dicomio.Encoder) {
	for _, item := range pdu.Items {
		item.Write(e)
	}
}

func (pdu *PDataTf) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("P_DATA_TF{items: ["))
	for i, item := range pdu.Items {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(item.String())
	}
	buf.WriteString("]}")
	return buf.String()
}

// fillString pads the string with " " up to the given length.
func fillString(v string, length int) string {
	if len(v) > length {
		return v[:16]
	}
	for len(v) < length {
		v += " "
	}
	return v
}
