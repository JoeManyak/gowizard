package gen

import "os"

type InterfaceMethod struct {
	Name string
	// Args should be a pairs of type and name
	Args    []string
	Returns []string
}

func AddInterface(file *os.File, name string, methods []InterfaceMethod) error {
	if _, err := file.WriteString("type " + name + " interface {\n"); err != nil {
		return err
	}

	for i := range methods {
		if _, err := file.WriteString(methods[i].Name + "("); err != nil {
			return err
		}

		for j := 0; j < len(methods[i].Args)-1; j++ {
			comma := ""
			if j < len(methods[i].Args)-2 {
				comma = ", "
			}

			if _, err := file.WriteString(methods[i].Args[j] + " " + methods[i].Args[j+1] + comma); err != nil {
				return err
			}
		}

		if _, err := file.WriteString(")"); err != nil {
			return err
		}

		switch len(methods[i].Returns) {
		case 0:
			if _, err := file.WriteString(" {\n"); err != nil {
				return err
			}
			break
		case 1:
			if _, err := file.WriteString(" " + methods[i].Returns[0] + " {\n"); err != nil {
				return err
			}
			break
		default:
			if _, err := file.WriteString(" ("); err != nil {
				return err
			}

			for j := range methods[i].Returns {
				comma := ""
				if j < len(methods[i].Returns)-1 {
					comma = ", "
				}
				if _, err := file.WriteString(methods[i].Returns[j] + comma); err != nil {
					return err
				}
			}
		}

		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}
