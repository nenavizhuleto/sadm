package sadm

func PrintHelp(c *Connection, commands []Command) error {
	c.Printf("available commands:\n")
	for _, cmd := range commands {
		if err := c.Printf("\t%s", cmd.String()); err != nil {
			return err
		}
	}
	return nil
}
