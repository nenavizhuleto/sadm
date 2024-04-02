package sadm

import "fmt"

func unknownCommandErr(cmd string) error {
	return fmt.Errorf("unknown command: %s", cmd)
}
