// Copyright 2015 The Ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore
//
package next

import (
	"fmt"
	"log"
	"reflect"
)

// genMsgCoder tries to generate an encoder and a decoder for a given message type.
func genMsgRPC(v interface{}) (string, string, error) {
	var e, d string
	var inBWrite bool
	n := fmt.Sprintf("%T", v)
	p := n[5:]
	n = n[5 : len(n)-3]
	e = fmt.Sprintf("func Marshal%v(b bytes.Buffer", n)
	d = fmt.Sprintf("func Unmarshall%v(d[]byte) (*", n)
	eParms := ""
	dRet := p + fmt.Sprintf(", error) {\n\tvar p *%v\n\tb := bytes.NewBuffer(d)\n", p)
	eCode := ""
	dCode := "\tvar u32 [4]byte\n\tvar u16 [2]byte\n\tvar l int\n"

	// Add the encoding boiler plate: 4 bytes of size to be filled in later,
	// The tag type, and the tag itself.
	eCode += "\tb.Write([]byte{0,0,0,0})\n\tb.Write([]byte{uint8(" + n + "),\n"
	inBWrite = true

	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		if !inBWrite {
			eCode += "\tb.Write([]byte{"
			inBWrite = true
		}
		f := t.Field(i)
		eParms += ", "
		n := f.Name
		eParms += fmt.Sprintf("%v %v", n, f.Type.Kind())
		switch f.Type.Kind() {
		case reflect.Uint32:
			eCode += fmt.Sprintf("\tuint8(%v>>24),uint8(%v>>16),", n, n)
			eCode += fmt.Sprintf("\tuint8(%v>>8),uint8(%v),\n", n, n)
			dCode += "\tif _, err := b.Read(u32[:]); err != nil {\n\t\treturn nil, fmt.Errorf(\"pkt too short for uint32: need 4, have %d\", b.Len())\n\t}\n"
			dCode += fmt.Sprintf("\tp.%v = uint32(u32[0])<<24|uint32(u32[1])<<16|uint32(u32[2])<<8|uint32(u32[3])\n", n)
		case reflect.Uint16:
			eCode += fmt.Sprintf("\tuint8(%v>>8),uint8(%v),\n", n, n)
			dCode += "\tif _, err := b.Read(u16[:]); err != nil {\n\t\treturn nil, fmt.Errorf(\"pkt too short for uint16: need 2, have %d\", b.Len())\n\t}\n"
			dCode += fmt.Sprintf("\tp.%v = uint16(u16[0])<<8|uint16(u16[1])\n", n)
		case reflect.String:
			eCode += fmt.Sprintf("\tuint8(len(%v)>>8),uint8(len(%v)),\n", n, n)
			if inBWrite {
				eCode += "\t})\n"
				inBWrite = false
			}
			eCode += fmt.Sprintf("\tb.Write([]byte(%v))\n", n)
			dCode += "\tif _, err := b.Read(u16[:]); err != nil {\n\t\treturn nil, fmt.Errorf(\"pkt too short for uint16: need 2, have %d\", b.Len())\n\t}\n"
			dCode += fmt.Sprintf("\tl = int(u16[0])<<8|int(u16[1])\n")
			dCode += "\tif b.Len() < l  {\n\t\treturn nil, fmt.Errorf(\"pkt too short for string: need %d, have %d\", l, b.Len())\n\t}\n"
			dCode += fmt.Sprintf("\tp.%v = b.String()\n", n)
		default:
			return "", "", fmt.Errorf("Can't encode %T.%v", v, f)
		}

	}
	if inBWrite {
		eCode += "\t})\n"
	}
	eCode += "\tl := b.Len()\n\tcopy(b.Bytes(), []byte{uint8(l>>24), uint8(l>>16), uint8(l>>8), uint8(l)})\n"
	return e + eParms + ") {\n" + eCode + "}\n", d + dRet + dCode + "\n\treturn p, nil\n}\n", nil
}

func main() {
	e, d, err := genMsgRPC(TversionPkt{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("package next\n\nimport (\n\t\"bytes\"\n\t\"fmt\"\n)\n%v \n %v \n", e, d)
}