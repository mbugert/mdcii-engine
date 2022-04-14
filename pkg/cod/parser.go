package cod

import (
	fmt "fmt"
	"strings"
)

func (c *Cod) Parse() error {
	for _, line := range c.Lines {
		spaces := countFrontSpaces(line)
		fmt.Println(spaces)
		line = strings.ReplaceAll(line, " ", "")
		if strings.Contains(line, "Nahrung:") || strings.Contains(line, "Soldat:") || strings.Contains(line, "Turm:") {
			continue
		}

		// Check relative array assignements before relative assignements
		// relative array assignment, examples:
		// example: '@Pos:       +0, +42'
		ok, err := c.relativeArrayAssignment(line)
		if ok {
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		// relative assignment, examples:
		// example: '@Pos: +42'
		ok, err = c.relativeAssignment(line)
		if ok {
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		// constant assignment, examples:
		// VARIABLEA = 5000
		// VARIABLEB = Nummer
		// Nummer = 1000
		// FOO = BAR+100
		ok, err = c.constantAssignment(line)
		if ok {
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}
	}
	return nil
}
