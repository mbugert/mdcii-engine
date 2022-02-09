package cod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestNewPng(t *testing.T) {
// 	// assert := assert.New(t)
// 	data := []byte{}
// 	cod := &Cod{
// 		Content: data,
// 		Decode:  false,
// 	}
// 	cod.Parse()
// }

func TestReadLines1(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\r\nfollowed by a new line\r\nfin\r\n"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test", string(bytes[0]))
	assert.Equal("followed by a new line", string(bytes[1]))
	assert.Equal("fin", string(bytes[2]))
}

func TestReadLines2(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\r\nfollowed by a new line\r\nfin"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test", string(bytes[0]))
	assert.Equal("followed by a new line", string(bytes[1]))
}

func TestReadLines3(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\rfollowed by a new line\r\nfin"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test\rfollowed by a new line", string(bytes[0]))
}

func TestReadLines4(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\rfollowed by a new line\r\nfin\r\n"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test\rfollowed by a new line", string(bytes[0]))
	assert.Equal("fin", string(bytes[1]))
}

func TestReadLines5(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\r\nfollowed by a new line\r\n;----comment\r\nfin\r\n"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test", string(bytes[0]))
	assert.Equal("followed by a new line", string(bytes[1]))
	assert.Equal("fin", string(bytes[2]))
}

func TestReadLines6(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test\r\nfollowed by a new line\r\n    ;----comment\r\nfin\r\n"
	bytes, err := readLines([]byte(data))
	assert.Nil(err)

	assert.Equal("this is a test", string(bytes[0]))
	assert.Equal("followed by a new line", string(bytes[1]))
	assert.Equal("fin", string(bytes[2]))
}

func TestCountFrontSpaces(t *testing.T) {
	assert := assert.New(t)
	data := "this is a test"
	assert.Equal(0, countFrontSpaces(data))
	data = " this is a test"
	assert.Equal(1, countFrontSpaces(data))
	data = "  this is a test"
	assert.Equal(2, countFrontSpaces(data))
	data = "   this is a test"
	assert.Equal(3, countFrontSpaces(data))
	data = "    this is a test"
	assert.Equal(4, countFrontSpaces(data))
}

func TestGetValueConstantNotExists(t *testing.T) {
	assert := assert.New(t)
	cod := Cod{
		Lines:   []string{},
		Objects: Objects{},
		Intern: CodIntern{
			variableNumbers:      map[string]int{},
			variableNumbersArray: map[string][]int{},
			constants:            Variables{},
		},
	}
	v, err := cod.getValue("FOO", "BAR+5", true)
	assert.Nil(err)
	assert.Equal(&Variable_ValueInt{ValueInt: 5}, v.Value)

}

func TestGetValueConstantExists(t *testing.T) {
	assert := assert.New(t)
	cod := Cod{
		Lines:   []string{},
		Objects: Objects{},
		Intern: CodIntern{
			variableNumbers:      map[string]int{},
			variableNumbersArray: map[string][]int{},
			constants:            Variables{},
		},
	}

	cod.Intern.constants.Variable = append(cod.Intern.constants.Variable,
		&Variable{
			Name:  "BAR",
			Value: &Variable_ValueInt{ValueInt: 101},
		},
	)
	cod.Intern.constants.Variable = append(cod.Intern.constants.Variable,
		&Variable{
			Name:  "BAZ",
			Value: &Variable_ValueInt{ValueInt: 100},
		},
	)
	v, err := cod.getValue("FOO", "BAR+5", true)
	assert.Nil(err)
	assert.Equal(&Variable_ValueInt{ValueInt: 106}, v.Value)
	v, err = cod.getValue("FOO", "BAZ-5", true)
	assert.Nil(err)
	assert.Equal(&Variable_ValueInt{ValueInt: 95}, v.Value)
}
