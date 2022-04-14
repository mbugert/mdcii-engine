package cod

import (
	fmt "fmt"
	"strconv"
	"strings"
)

// relative array assignment, examples:
// example: '@Pos:       +0, +42'
func (c *Cod) relativeArrayAssignment(line string) (bool, error) {
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
						return true, err
					}
					if number[0] == "-" {
						offset *= -1
					}
					offsets = append(offsets, int(offset))

				}
			}
		}

		index := c.existsInCurrentObject(name)
		currentArrayValues := []int{}
		for i := 0; i < len(offsets); i++ {
			var currentValue int
			if index != -1 {
				if _, ok := c.Intern.variableNumbersArray[c.Intern.currentObject.Name]; ok {
					currentValue = c.Intern.variableNumbersArray[c.Intern.currentObject.Name][i]
					currentValue = calculateOperation(currentValue, "+", offsets[i])
					c.Intern.currentObject.Variables.Variable[i].Value = &Variable_ValueInt{
						ValueInt: int32(currentValue),
					}
				} else {
					return true, fmt.Errorf("no current object")
				}
			} else {
				if _, ok := c.Intern.variableNumbersArray[name]; ok {
					currentValue = calculateOperation(c.Intern.variableNumbersArray[name][i], "+", offsets[i])
					variable, err := c.createOrReuseVariable(name)
					if err != nil {
						return true, err
					}
					if variable != nil {
						variable.Name = name
						variable.Value = &Variable_ValueInt{
							ValueInt: int32(currentValue),
						}
					}
				} else {
					return true, fmt.Errorf("no current object")
				}
			}
			currentArrayValues = append(currentArrayValues, currentValue)
		}
		c.Intern.variableNumbersArray[name] = currentArrayValues
	} else {
		return false, fmt.Errorf("line no relative array assignment")
	}
	return true, nil
}

// relative assignment, examples:
// example: '@Pos: +42'
func (c *Cod) relativeAssignment(line string) (bool, error) {
	if matches := regexSearch(`(@)(\w+)\s*:\s*([+|-]?)(\d+)`, line); len(matches) > 0 {
		name := matches[1]
		if name == NUMBER_INCREMENT {
			return false, fmt.Errorf("line no relative assignment")
		}
		index := c.existsInCurrentObject(name)
		operand, err := strconv.Atoi(matches[3])
		if err != nil {
			return true, err
		}
		currentValue := 0
		if index != -1 {
			currentValue = c.Intern.variableNumbers[c.Intern.currentObject.Variables.Variable[index].Name]
			currentValue = calculateOperation(currentValue, matches[2], operand)
			c.Intern.currentObject.Variables.Variable[index].Value = &Variable_ValueInt{
				ValueInt: int32(currentValue),
			}
		} else {
			currentValue = c.Intern.variableNumbers[matches[1]]
			currentValue = calculateOperation(currentValue, matches[2], operand)
			variable, err := c.createOrReuseVariable(name)
			if err != nil {
				return true, err
			}
			if variable != nil {
				variable.Name = name
				variable.Value = &Variable_ValueInt{
					ValueInt: int32(currentValue),
				}
			}
		}
		c.Intern.variableNumbers[name] = currentValue
	} else {
		return false, fmt.Errorf("line no relative assignment")
	}
	return true, nil
}
