package dimse

//go:generate ./generate_dimse_messages.py
//go:generate stringer -type StatusCode

// Implements message types defined in P3.7.
//
// http://dicom.nema.org/medical/dicom/current/output/pdf/part07.pdf

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomtag"
	"github.com/grailbio/go-netdicom/pdu"
	"v.io/x/lib/vlog"
)

// Message defines the common interface for all DIMSE message types.
type Message interface {
	fmt.Stringer // Print human-readable description for debugging.
	Encode(*dicomio.Encoder)
	GetMessageID() uint16 // Extract the message ID field.
	CommandField() int    // Return the command field value of this message.
	HasData() bool        // Do we expact data P_DATA_TF packets after the command packets?
}

// Status represents a result of a DIMSE call.  P3.7 C defines list of status
// codes and error payloads.
type Status struct {
	// Status==StatusSuccess on success. A non-zero value on error.
	Status StatusCode

	// Optional error payloads.
	ErrorComment string // Encoded as (0000,0902)
}

// Helper class for extracting values from a list of DicomElement.
type messageDecoder struct {
	elems  []*dicom.Element
	parsed []bool // true if this element was parsed into a message field.
	err    error
}

type isOptionalElement int

const (
	requiredElement isOptionalElement = iota
	optionalElement
)

func (d *messageDecoder) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// Find an element with the given tag. If optional==OptionalElement, returns nil
// if not found.  If optional==RequiredElement, sets d.err and return nil if not found.
func (d *messageDecoder) findElement(tag dicomtag.Tag, optional isOptionalElement) *dicom.Element {
	for i, elem := range d.elems {
		if elem.Tag == tag {
			vlog.VI(3).Infof("Return %v for %s", elem, tag.String())
			d.parsed[i] = true
			return elem
		}
	}
	if optional == requiredElement {
		d.setError(fmt.Errorf("Element %s not found during DIMSE decoding", dicomtag.DebugString(tag)))
	}
	return nil
}

// Return the list of elements that did not match any of the prior getXXX calls.
func (d *messageDecoder) unparsedElements() (unparsed []*dicom.Element) {
	for i, parsed := range d.parsed {
		if !parsed {
			unparsed = append(unparsed, d.elems[i])
		}
	}
	return unparsed
}

func (d *messageDecoder) getStatus() (s Status) {
	s.Status = StatusCode(d.getUInt16(dicomtag.Status, requiredElement))
	s.ErrorComment = d.getString(dicomtag.ErrorComment, optionalElement)
	return s
}

// Find an element with "tag", and extract a string value from it. Errors are reported in d.err.
func (d *messageDecoder) getString(tag dicomtag.Tag, optional isOptionalElement) string {
	e := d.findElement(tag, optional)
	if e == nil {
		return ""
	}
	v, err := e.GetString()
	if err != nil {
		d.setError(err)
	}
	return v
}

// Find an element with "tag", and extract a uint32 from it. Errors are reported in d.err.
func (d *messageDecoder) getUInt32(tag dicomtag.Tag, optional isOptionalElement) uint32 {
	e := d.findElement(tag, optional)
	if e == nil {
		return 0
	}
	v, err := e.GetUInt32()
	if err != nil {
		d.setError(err)
	}
	return v
}

// Find an element with "tag", and extract a uint16 from it. Errors are reported in d.err.
func (d *messageDecoder) getUInt16(tag dicomtag.Tag, optional isOptionalElement) uint16 {
	e := d.findElement(tag, optional)
	if e == nil {
		return 0
	}
	v, err := e.GetUInt16()
	if err != nil {
		d.setError(err)
	}
	return v
}

// Encode a DIMSE field with the given tag, given value "v"
func encodeField(e *dicomio.Encoder, tag dicomtag.Tag, v interface{}) {
	elem := dicom.Element{
		Tag:             tag,
		VR:              "", // autodetect
		UndefinedLength: false,
		Value:           []interface{}{v},
	}
	dicom.WriteElement(e, &elem)
}

// CommandDataSetTypeNull indicates that the DIMSE message has no data payload,
// when set in dicom.TagCommandDataSetType. Any other value indicates the
// existence of a payload.
const CommandDataSetTypeNull uint16 = 0x101

// CommandDataSetTypeNonNull indicates that the DIMSE message has a data
// payload, when set in dicom.TagCommandDataSetType.
const CommandDataSetTypeNonNull uint16 = 1

// Success is an OK status for a call.
var Success = Status{Status: StatusSuccess}

// StatusCode represents a DIMSE service response code, as defined in P3.7
type StatusCode uint16

const (
	StatusSuccess               StatusCode = 0
	StatusCancel                StatusCode = 0xFE00
	StatusSOPClassNotSupported  StatusCode = 0x0112
	StatusInvalidArgumentValue  StatusCode = 0x0115
	StatusInvalidAttributeValue StatusCode = 0x0106
	StatusInvalidObjectInstance StatusCode = 0x0117
	StatusUnrecognizedOperation StatusCode = 0x0211
	StatusNotAuthorized         StatusCode = 0x0124
	StatusPending               StatusCode = 0xff00

	// C-STORE-specific status codes. P3.4 GG4-1
	CStoreOutOfResources              StatusCode = 0xa700
	CStoreCannotUnderstand            StatusCode = 0xc000
	CStoreDataSetDoesNotMatchSOPClass StatusCode = 0xa900

	// C-FIND-specific status codes.
	CFindUnableToProcess StatusCode = 0xc000

	// C-MOVE/C-GET-specific status codes.
	CMoveOutOfResourcesUnableToCalculateNumberOfMatches StatusCode = 0xa701
	CMoveOutOfResourcesUnableToPerformSubOperations     StatusCode = 0xa702
	CMoveMoveDestinationUnknown                         StatusCode = 0xa801
	CMoveDataSetDoesNotMatchSOPClass                    StatusCode = 0xa900

	// Warning codes.
	StatusAttributeValueOutOfRange StatusCode = 0x0116
	StatusAttributeListError       StatusCode = 0x0107
)

func encodeStatus(e *dicomio.Encoder, s Status) {
	encodeField(e, dicomtag.Status, uint16(s.Status))
	if s.ErrorComment != "" {
		encodeField(e, dicomtag.ErrorComment, s.ErrorComment)
	}
}

// ReadMessage constructs a typed dimse.Message object, given a set of
// dicom.Elements,
func ReadMessage(d *dicomio.Decoder) Message {
	// A DIMSE message is a sequence of Elements, encoded in implicit
	// LE.
	//
	// TODO(saito) make sure that's the case. Where the ref?
	var elems []*dicom.Element
	d.PushTransferSyntax(binary.LittleEndian, dicomio.ImplicitVR)
	defer d.PopTransferSyntax()
	for d.Len() > 0 {
		elem := dicom.ReadElement(d, dicom.ReadOptions{})
		if d.Error() != nil {
			break
		}
		elems = append(elems, elem)
	}

	// Convert elems[] into a golang struct.
	dd := messageDecoder{
		elems:  elems,
		parsed: make([]bool, len(elems)),
		err:    nil,
	}
	commandField := dd.getUInt16(dicomtag.CommandField, requiredElement)
	if dd.err != nil {
		d.SetError(dd.err)
		return nil
	}
	v := decodeMessageForType(&dd, commandField)
	if dd.err != nil {
		d.SetError(dd.err)
		return nil
	}
	return v
}

// EncodeMessage serializes the given message. Errors are reported through e.Error()
func EncodeMessage(e *dicomio.Encoder, v Message) {
	// DIMSE messages are always encoded Implicit+LE. See P3.7 6.3.1.
	subEncoder := dicomio.NewBytesEncoder(binary.LittleEndian, dicomio.ImplicitVR)
	v.Encode(subEncoder)
	if err := subEncoder.Error(); err != nil {
		e.SetError(err)
		return
	}
	bytes := subEncoder.Bytes()
	e.PushTransferSyntax(binary.LittleEndian, dicomio.ImplicitVR)
	defer e.PopTransferSyntax()
	encodeField(e, dicomtag.CommandGroupLength, uint32(len(bytes)))
	e.WriteBytes(bytes)
}

// CommandAssembler is a helper that assembles a DIMSE command message and data
// payload from a sequence of P_DATA_TF PDUs.
type CommandAssembler struct {
	contextID      byte
	commandBytes   []byte
	command        Message
	dataBytes      []byte
	readAllCommand bool

	readAllData bool
}

// AddDataPDU is to be called for each P_DATA_TF PDU received from the
// network. If the fragment is marked as the last one, AddDataPDU returns
// <SOPUID, TransferSyntaxUID, payload, nil>.  If it needs more fragments, it
// returns <"", "", nil, nil>.  On error, it returns a non-nil error.
func (a *CommandAssembler) AddDataPDU(pdu *pdu.PDataTf) (byte, Message, []byte, error) {
	for _, item := range pdu.Items {
		if a.contextID == 0 {
			a.contextID = item.ContextID
		} else if a.contextID != item.ContextID {
			return 0, nil, nil, fmt.Errorf("Mixed context: %d %d", a.contextID, item.ContextID)
		}
		if item.Command {
			a.commandBytes = append(a.commandBytes, item.Value...)
			if item.Last {
				if a.readAllCommand {
					return 0, nil, nil, fmt.Errorf("P_DATA_TF: found >1 command chunks with the Last bit set")
				}
				a.readAllCommand = true
			}
		} else {
			a.dataBytes = append(a.dataBytes, item.Value...)
			if item.Last {
				if a.readAllData {
					return 0, nil, nil, fmt.Errorf("P_DATA_TF: found >1 data chunks with the Last bit set")
				}
				a.readAllData = true
			}
		}
	}
	if !a.readAllCommand {
		return 0, nil, nil, nil
	}
	if a.command == nil {
		d := dicomio.NewBytesDecoder(a.commandBytes, nil, dicomio.UnknownVR)
		a.command = ReadMessage(d)
		if err := d.Finish(); err != nil {
			return 0, nil, nil, err
		}
	}
	if a.command.HasData() && !a.readAllData {
		return 0, nil, nil, nil
	}
	contextID := a.contextID
	command := a.command
	dataBytes := a.dataBytes
	*a = CommandAssembler{}
	return contextID, command, dataBytes, nil
	// TODO(saito) Verify that there's no unread items after the last command&data.
}

// Generate a new message ID that's unique within the "su".
var nextMessageID int32 = 123

// NewMessageID allocates a yet-unused message ID.
func NewMessageID() uint16 {
	id := atomic.AddInt32(&nextMessageID, 1)
	return uint16(id % 0x10000)
}
