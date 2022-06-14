package oboparser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

type parseError struct {
	lineNum int
	err     error
}

func (p *parseError) Error() string {
	return fmt.Sprintf("%v in line %d", p.err, p.lineNum)
}

var (
	hideCharacterReg = regexp.MustCompile("/[^(\\x20-\\x7F)]*/")
)

type ParserConfig struct {
	ShowNotExistField bool
}

type parser struct {
	rd      io.Reader
	lineNum int
	numRead int
	obo     *OBO
	config  *ParserConfig
}

func New(rd io.Reader, config *ParserConfig) *parser {
	return &parser{rd: rd, config: config}
}

func (p *parser) Execute() *OBO {
	var field []byte
	var value []byte

	p.obo = newOBO()
	inHeader := true
	var stanza *Stanza
	line := make([]byte, 0)
	var prevStanzaType StanzaType

	reader := bufio.NewReader(p.rd)
	for {
		bs, isPrefix, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if stanza != nil {
					p.obo.Stanzas = append(p.obo.Stanzas, stanza)
				}
				break
			} else {
				panic(parseError{lineNum: p.lineNum, err: err})
			}
		}
		line = append(line, bs...)
		if isPrefix {
			continue
		}

		p.lineNum++
		p.numRead += len(line)

		line = bytes.Trim(line, " ")
		// remove hidden characters
		line = hideCharacterReg.ReplaceAll(line, []byte(""))
		if len(line) == 0 {
			continue
		}

		// remove comments
		line = bytes.Trim(bytes.Split(line, []byte("!"))[0], " ")
		if bytes.HasPrefix(line, []byte("[")) {
			inHeader = false
			switch string(bytes.Trim(line, "[]")) {
			case "Term":
				prevStanzaType = TermStanza
			case "Typedef":
				prevStanzaType = TypedefStanza
			case "Instance":
				prevStanzaType = InstanceStanza
			default:
				panic(parseError{lineNum: p.lineNum, err: errors.New("stanza must be one of Term, Typedef and Instance")})
			}
			if stanza != nil {
				p.obo.Stanzas = append(p.obo.Stanzas, stanza)
			}
			stanza = &Stanza{Type: prevStanzaType, Values: make(map[string]string)}
			line = line[:0]
			continue
		}

		if inHeader {
			field, value, err = fieldValue(line, p.lineNum)
			if err != nil {
				panic(err)
			}
			if p.config.ShowNotExistField && !stringIn(string(field), headerFields) {
				fmt.Println(fmt.Sprintf("[Notice] field %s is not in pre-setting header fields, in line %d", field, p.lineNum))
			}
			p.obo.SetHeader(string(field), string(value))
		}

		if stanza != nil {
			field, value, err = fieldValue(line, p.lineNum)
			if err != nil {
				panic(err)
			}
			if p.config.ShowNotExistField {
				switch stanza.Type {
				case TermStanza:
					if !stringIn(string(field), termFields) {
						fmt.Println(fmt.Sprintf("[Notice] field %s is not in pre-setting Term's fields, in line %d", field, p.lineNum))
					}
				case TypedefStanza:
					if !stringIn(string(field), typedefFields) {
						fmt.Println(fmt.Sprintf("[Notice] field %s is not in pre-setting Typedef's fields, in line %d", field, p.lineNum))
					}
				case InstanceStanza:
					if !stringIn(string(field), instanceFields) {
						fmt.Println(fmt.Sprintf("[Notice] field %s is not in pre-setting Instance's fields, in line %d", field, p.lineNum))
					}
				}
			}
			stanza.Values[string(field)] = string(value)
		}

		line = line[:0]
	}

	return p.obo
}

func fieldValue(line []byte, lineNum int) (field, value []byte, er error) {
	idx := bytes.IndexByte(line, ':')
	if idx == -1 {
		return nil, nil, &parseError{lineNum: lineNum, err: errors.New("should contains a colon symbol")}
	}
	return bytes.Trim(line[:idx], " "), bytes.Trim(line[idx+1:], " "), nil
}
