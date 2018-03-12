#!/usr/bin/env python3.6

import enum
from typing import IO, List, NamedTuple

Field = NamedTuple('Field', [('name', str),
                             ('type', str),
                             ('required', bool)])
class Type(enum.Enum):
    REQUEST = 1
    RESPONSE = 2

Message = NamedTuple('Message',
                     [('name', str),
                      ('type', Type),
                      ('command_field', int),
                      ('fields', List[Field])])

MESSAGES = [
    # P3.7 9.3.1.1
    Message('CStoreRq',
            Type.REQUEST, 1,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageID', 'MessageID', True),
             Field('Priority', 'uint16', True),
             Field('CommandDataSetType', 'uint16', True),
             Field('AffectedSOPInstanceUID', 'string', True),
	     Field('MoveOriginatorApplicationEntityTitle', 'string', False),
	     Field('MoveOriginatorMessageID', 'MessageID', False)]),
    # P3.7 9.3.1.2
    Message('CStoreRsp',
            Type.RESPONSE, 0x8001,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageIDBeingRespondedTo', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True),
             Field('AffectedSOPInstanceUID', 'string', True),
	     Field('Status', 'Status', True)]),
    # P3.7 9.1.2.1
    Message('CFindRq',
            Type.REQUEST, 0x20,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageID', 'MessageID', True),
             Field('Priority', 'uint16', True),
             Field('CommandDataSetType', 'uint16', True)]),
    Message('CFindRsp',
            Type.RESPONSE, 0x8020,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageIDBeingRespondedTo', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True),
	     Field('Status', 'Status', True)]),
    # P3.7 9.1.2.1
    Message('CGetRq',
            Type.REQUEST, 0x10,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageID', 'MessageID', True),
             Field('Priority', 'uint16', True),
             Field('CommandDataSetType', 'uint16', True)]),
    Message('CGetRsp',
            Type.RESPONSE, 0x8010,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageIDBeingRespondedTo', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True),
             Field('NumberOfRemainingSuboperations', 'uint16', False),
             Field('NumberOfCompletedSuboperations', 'uint16', False),
             Field('NumberOfFailedSuboperations', 'uint16', False),
             Field('NumberOfWarningSuboperations', 'uint16', False),
	     Field('Status', 'Status', True)]),
    # P3.7 9.3.4.1
    Message('CMoveRq',
            Type.REQUEST, 0x21,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageID', 'MessageID', True),
             Field('Priority', 'uint16', True),
             Field('MoveDestination', 'string', True),
             Field('CommandDataSetType', 'uint16', True)]),
    Message('CMoveRsp',
            Type.RESPONSE, 0x8021,
            [Field('AffectedSOPClassUID', 'string', True),
             Field('MessageIDBeingRespondedTo', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True),
             Field('NumberOfRemainingSuboperations', 'uint16', False),
             Field('NumberOfCompletedSuboperations', 'uint16', False),
             Field('NumberOfFailedSuboperations', 'uint16', False),
             Field('NumberOfWarningSuboperations', 'uint16', False),
	     Field('Status', 'Status', True)]),
    # P3.7 9.3.5
    Message('CEchoRq',
            Type.REQUEST, 0x30,
            [Field('MessageID', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True)]),
    Message('CEchoRsp',
            Type.RESPONSE, 0x8030,
            [Field('MessageIDBeingRespondedTo', 'MessageID', True),
             Field('CommandDataSetType', 'uint16', True),
	     Field('Status', 'Status', True)])
]

def generate_go_definition(m: Message, out: IO[str]):
    print(f'type {m.name} struct {{', file=out)
    for f in m.fields:
        print(f'	{f.name} {f.type}', file=out)
    print(f'	Extra []*dicom.Element  // Unparsed elements', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) Encode(e *dicomio.Encoder) {{', file=out)
    print('    elems := []*dicom.Element{}', file=out)
    print(f'	elems = append(elems, newElement(dicomtag.CommandField, uint16({m.command_field})))', file=out)
    for f in m.fields:
        if not f.required:
            if f.type == 'string':
                zero = '""'
            else:
                zero = '0'
            print(f'	if v.{f.name} != {zero} {{', file=out)
            print(f'		elems = append(elems, newElement(dicomtag.{f.name}, v.{f.name}))', file=out)
            print(f'	}}', file=out)
        elif f.type == 'Status':
            print(f'	elems = append(elems, newStatusElements(v.{f.name})...)', file=out)
        else:
            print(f'	elems = append(elems, newElement(dicomtag.{f.name}, v.{f.name}))', file=out)
    print('	elems = append(elems, v.Extra...)', file=out)
    print('	encodeElements(e, elems)', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) HasData() bool {{', file=out)
    print(f'	return v.CommandDataSetType != CommandDataSetTypeNull', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) CommandField() int {{', file=out)
    print(f'	return {m.command_field}', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) GetMessageID() MessageID {{', file=out)
    if m.type == Type.REQUEST:
        print(f'	return v.MessageID', file=out)
    else:
        print(f'	return v.MessageIDBeingRespondedTo', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) GetStatus() *Status {{', file=out)
    if m.type == Type.REQUEST:
        print(f'	return nil', file=out)
    else:
        print(f'	return &v.Status', file=out)
    print('}', file=out)

    print('', file=out)
    print(f'func (v *{m.name}) String() string {{', file=out)
    i = 0
    fmt = f'{m.name}{{'
    args = ''
    for f in m.fields:
        space = ''
        if i > 0:
            fmt += ' '
            args += ', '
        fmt += f'{f.name}:%v'
        args += f'v.{f.name}'
        i += 1
    fmt += "}}"
    print(f'	return fmt.Sprintf("{fmt}", {args})', file=out)
    print('}', file=out)


    print('', file=out)
    print(f'func decode{m.name}(d *messageDecoder) *{m.name} {{', file=out)
    print(f'	v := &{m.name}{{}}', file=out)
    for f in m.fields:
        if f.type == 'Status':
            print(f'	v.{f.name} = d.getStatus()', file=out)
        else:
            if f.type == 'string':
                decoder = 'String'
            elif f.type in ('uint16', 'MessageID'):
                decoder = 'UInt16'
            elif f.type == 'uint32':
                decoder = 'UInt32'
            else:
                raise Exception(f)
            if f.required:
                required = 'requiredElement'
            else:
                required = 'optionalElement'
            print(f'	v.{f.name} = d.get{decoder}(dicomtag.{f.name}, {required})', file=out)
    print(f'	v.Extra = d.unparsedElements()', file=out)
    print(f'	return v', file=out)
    print('}', file=out)

def main():
    with open('dimse_messages.go', 'w') as out:
        print("""
package dimse

// Code generated from generate_dimse_messages.py. DO NOT EDIT.

import (
	"fmt"

	"github.com/grailbio/go-dicom"
	"github.com/grailbio/go-dicom/dicomio"
	"github.com/grailbio/go-dicom/dicomtag"
)

        """, file=out)
        for m in MESSAGES:
            generate_go_definition(m, out)

        for m in MESSAGES:
            print(f'const CommandField{m.name} = {m.command_field}', file=out)

        print('func decodeMessageForType(d* messageDecoder, commandField uint16) Message {', file=out)
        print('	switch commandField {', file=out)
        for m in MESSAGES:
            print('	case 0x%x:' % (m.command_field, ), file=out)
            print(f'		return decode{m.name}(d)', file=out)
        print('	default:', file=out)
        print('		d.setError(fmt.Errorf("Unknown DIMSE command 0x%x", commandField))', file=out)
        print('		return nil', file=out)
        print('	}', file=out)
        print('}', file=out)

main()
