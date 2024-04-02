package perfmon

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/nenavizhuleto/sadm"
)

func cpuInfo(c *sadm.Connection) error {
	info, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		return err
	}

	p := info.Processors[0]

	c.Println("model name:\t", p.ModelName)
	c.Println("vendor id:\t", p.VendorId)
	c.Println("num cores:\t", p.Cores)

	return nil
}

func stat(c *sadm.Connection) error {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return err
	}

	c.Printf("procs:\t%d\n", stat.Processes)
	c.Printf("boot:\t%s\n", stat.BootTime)

	c.Println()

	s := stat.CPUStatAll
	c.Println("WARN:\tcpu stats are inaccurate")
	c.Printf("user:\t%.02f\n", float64(s.User/1000))
	c.Printf("system:\t%.02f\n", float64(s.System/1000))

	return nil
}
