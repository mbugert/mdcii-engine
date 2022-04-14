package cod

import (
	"bytes"
	fmt "fmt"
	"strings"
)

func (c *Cod) createOrReuseVariable(name string) (*Variable, error) {
	if c.Intern.currentObject != nil {
		optionalVar := c.getVariableFromObject(c.Intern.currentObject, name)
		if optionalVar != nil {
			return optionalVar, nil
		}
		variable := &Variable{}
		variables := &Variables{
			Variable: []*Variable{variable},
		}
		c.Intern.currentObject.Variables = variables

		return variable, nil
	}
	return nil, fmt.Errorf("no current object")
}

func (c *Cod) getVariableFromConstants(key string) Variable {
	for i := 0; i < len(c.Intern.constants.Variable); i++ {
		if c.Intern.constants.Variable[i].Name == key {
			return *c.Intern.constants.Variable[i]
		}
	}
	return Variable{}
}

func (c *Cod) getVariableFromObject(obj *Object, name string) *Variable {
	if obj.Variables != nil {
		for i := 0; i < len(obj.Variables.Variable); i++ {
			if obj.Variables.Variable[i].Name == name {
				return obj.Variables.Variable[i]
			}
		}
	}
	return nil
}

func calculateOperation(oldValue int, operation string, op int) int {
	currentValue := oldValue
	if operation == "+" {
		currentValue += op
	} else if operation == "-" {
		currentValue -= op
	} else {
		currentValue = op
	}
	return currentValue
}

func (c *Cod) existsInCurrentObject(variableName string) int {
	if c.Intern.currentObjectIndex != -1 {
		for index, v := range c.Intern.currentObject.Variables.Variable {
			if v.Name == variableName {
				return index
			}
		}
	}
	return -1
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
