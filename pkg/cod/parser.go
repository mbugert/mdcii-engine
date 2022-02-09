package cod

import (
	"bytes"
	fmt "fmt"
	"io/ioutil"
	"strconv"
	"strings"

	stack "github.com/siredmar/mdcii-engine/pkg/stack"
)

const (
	NUMBER_INCREMENT = "Nummer"
)

type CodIntern struct {
	variableNumbers      map[string]int
	variableNumbersArray map[string][]int
	constants            Variables
	currentObjectIndex   int
	objectStack          stack.Stack
}

type Cod struct {
	Lines   []string
	Objects Objects
	Intern  CodIntern
}

type ObjectType struct {
	Object       *Object
	numberObject bool
	spaces       int
}

func NewCod(file string, decode bool) (*Cod, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return &Cod{}, err
	}
	if decode {
		for i, v := range b {
			b[i] = -v
		}
	}
	lines, err := readLines(b)
	if err != nil {
		return &Cod{}, err
	}
	return &Cod{
		Lines: lines,
		Intern: CodIntern{
			variableNumbers:      map[string]int{},
			variableNumbersArray: map[string][]int{},
			constants:            Variables{},
			currentObjectIndex:   -1,
			objectStack:          stack.Stack{},
		},
	}, nil
}

// bytesIndex finds the start index for a subslice of slice of bytes
func bytesIndex(b []byte, s []byte) int {
	if len(s) > len(b) {
		return -1
	}

	for i := 0; i <= len(b)-len(s); i++ {
		if bytes.Equal(b[i:i+len(s)], s) {
			return i
		}
	}

	return -1
}

// readLines reads all lines from a file passed as slice of bytes and returns slice of strings.
// The end of line is marked with an '\r' followed by an '\n'.
// Lines beginning with `;' are ignored
func readLines(b []byte) ([]string, error) {
	var lines []string
	for {
		idx := bytesIndex(b, []byte{'\r', '\n'})
		if idx == -1 {
			break
		}

		line := string(b[:idx])
		if len(line) > 0 && !strings.Contains(line, ";") {
			lines = append(lines, line)
		}

		b = b[idx+2:]
	}

	return lines, nil
}

// coundFrontSpaces counts the number of leading spaces in a string.
func countFrontSpaces(s string) int {
	var count int
	for _, r := range s {
		if r != ' ' {
			break
		}
		count++
	}
	return count
}

func (c *Cod) Parse() error {
	for _, line := range c.Lines {
		spaces := countFrontSpaces(line)
		fmt.Println(spaces)
		if strings.Contains(line, "Nahrung:") || strings.Contains(line, "Soldat:") || strings.Contains(line, "Turm:") {
			continue
		}
		// contant assignments, examples:
		// KEY = Nummer
		// KEY =   VALUE
		// FOO = BAR+123
		if matches := regexSearch(`(\w+)\s*=\s*((?:\d+|\+|\w+)+)`, line); len(matches) > 0 {
			c.constantAssignment(&constantType{
				key:   matches[1],
				value: matches[2],
			})
			continue
		}
		// example: '@Pos:       +0, +42'
		if matches := regexSearch(`@(\w+):.*(,)`, line); len(matches) > 0 {
			name := matches[0]
			offsets := []int{}
			if result := regexSearch(`:\s*(.*)`, line); len(result) > 0 {
				tokens := strings.Split(result[0], ",")
				for _, e := range tokens {
					e = strings.TrimSpace(e)
					if number := regexSearch(`([+|-]?)(\d+)`, e); len(number) > 0 {
						offset, err := strconv.Atoi(number[1])
						if err != nil {
							return err
						}
						if number[2] == "-" {
							offset *= -1
						}
						offsets = append(offsets, offset)

					}
				}
			}
			index := c.existsInCurrentObject(name)
		}
	}
}

func (c *Cod) existsInCurrentObject(variableName string) int {
	if c.Intern.currentObjectIndex != -1 {
		for k, v := range c.Intern.objectStack {

		}
	}
	return -1
}

type constantType struct {
	key   string
	value string
}

// constantExists returns the index of the constant with the given key.
// returns -1 if not found.
func (c *Cod) constantExists(key string) int {
	for i, constant := range c.Intern.constants.Variable {
		if constant.Name == key {
			return i
		}
	}
	return -1
}

// constantNumberAssignment checks if the constant in constant.key already exists
// if it exits, overwrite with current `Nummer` value
// if it doesn't exist create it wiht the current `Nummer` value
func (c *Cod) constantNumberAssignment(constant *constantType) {
	// example: 'HAUSWACHS = Nummer'
	// if constant.value == NUMBER_INCREMENT {
	if _, ok := c.Intern.variableNumbers[NUMBER_INCREMENT]; ok {
		number := c.Intern.variableNumbers[NUMBER_INCREMENT]
		i := c.constantExists(constant.key)
		if i != -1 {
			variable := c.Intern.constants.Variable[i]
			variable.Name = constant.key
			variable.Value = &Variable_ValueString{
				ValueString: fmt.Sprintf("%d", number),
			}
		} else {
			variable := &Variable{}
			variable.Name = constant.key
			variable.Value = &Variable_ValueString{
				ValueString: fmt.Sprintf("%d", number),
			}
			c.Intern.constants.Variable = append(c.Intern.constants.Variable, variable)
		}
	}
	// }
}

// constantAssignment checks if the constant in constant.key already exists
// if it exits, overwrite with the provided value
// if it doesn't exist create it with the provided value
func (c *Cod) constantAssignment(constant *constantType) {
	// example: 'FOO =   12345'
	i := c.constantExists(constant.key)
	if i != -1 {
		// TODO
	} else {
		// TODO
	}
}

// func (c *Cod) constantIncrementAssignment() {

// }

func (c *Cod) getValue(key string, value string, isMath bool) (*Variable, error) {
	ret := Variable{
		Name:  key,
		Value: nil,
	}

	if isMath {
		// Searching for some characters followed by a + or - sign and some digits.
		// example: VALUE+20
		matches := regexSearch(`(\w+)(\+|\-)(\d+)`, value)
		key := matches[0]
		if len(matches) > 0 {
			oldValue := Variable{
				Name:  key,
				Value: nil,
			}
			i := c.constantExists(oldValue.Name)
			if i != -1 {
				oldValue.Value = c.Intern.constants.Variable[i].Value
			} else {
				oldValue.Value = &Variable_ValueInt{
					ValueInt: 0,
				}
			}
			switch v := oldValue.Value.(type) {

			// Special case handling for `RUINE_KONTOR_1`
			case *Variable_ValueString:
				if v.ValueString == "RUINE_KONTOR_1" {
					ret.Value = &Variable_ValueInt{
						ValueInt: 424242,
					}
					return &ret, nil
				}
			}
			if matches[1] == "+" {
				switch v := oldValue.Value.(type) {
				case *Variable_ValueInt:
					i, err := strconv.Atoi(matches[2])
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueInt{
						ValueInt: v.ValueInt + int32(i),
					}
					return &ret, nil
				case *Variable_ValueFloat:
					i, err := strconv.ParseFloat(matches[2], 64)
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueFloat{
						ValueFloat: v.ValueFloat + float32(i),
					}
					return &ret, nil
				case *Variable_ValueString:
					i, err := strconv.Atoi(matches[2])
					if err != nil {
						return &ret, err
					}
					old, err := strconv.Atoi(v.ValueString)
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueInt{
						ValueInt: int32(old) + int32(i),
					}
					return &ret, nil
				}
			}
			if matches[1] == "-" {
				switch v := oldValue.Value.(type) {
				case *Variable_ValueInt:
					i, err := strconv.Atoi(matches[2])
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueInt{
						ValueInt: v.ValueInt - int32(i),
					}
					return &ret, nil
				case *Variable_ValueFloat:
					i, err := strconv.ParseFloat(matches[2], 64)
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueFloat{
						ValueFloat: v.ValueFloat - float32(i),
					}
					return &ret, nil
				case *Variable_ValueString:
					i, err := strconv.Atoi(matches[2])
					if err != nil {
						return &ret, err
					}
					old, err := strconv.Atoi(v.ValueString)
					if err != nil {
						return &ret, err
					}
					ret.Value = &Variable_ValueInt{
						ValueInt: int32(old) - int32(i),
					}
					return &ret, nil
				}
			}
		}
	}
	return nil, nil
}
