package cod

import (
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
)

type CodIntern struct {
	variableNumbers      map[string]int
	variableNumbersArray map[string][]int
	constants            Variables
	currentObjectIndex   int
	currentObject        *Object
	objectStack          *Stack
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
			constants: Variables{
				Variable: []*Variable{},
			},
			currentObjectIndex: -1,
			objectStack:        NewStack(),
			currentObject:      nil,
		},
	}, nil
}

func (cod *Cod) DumpConstants() error {
	jsonBytes, err := json.MarshalIndent(cod.Intern.constants.Variable, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println("Constants:", string(jsonBytes))
	return nil
}

func (cod *Cod) DumpVariables() error {
	fmt.Println("Variables:", cod.Intern.variableNumbers)
	return nil
}

func (cod *Cod) DumpArrayVariables() error {
	fmt.Println("Array Variables:", cod.Intern.variableNumbersArray)
	return nil
}
