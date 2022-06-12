package cpu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ds248a/lib/osx"
)

type InfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
}

// CPUInfo on linux will return 1 item per physical thread.
//
// CPUs have three levels of counting: sockets, cores, threads.
// Cores with HyperThreading count as having 2 threads per core.
// Sockets often come with many physical CPU cores.
// For example a single socket board with two cores each with HT will
// return 4 CPUInfoStat structs on Linux and the "Cores" field set to 1.
func Info() ([]InfoStat, error) {
	filename := osx.HostProc("cpuinfo")
	lines, _ := osx.ReadLines(filename)

	var ret []InfoStat
	c := InfoStat{CPU: -1, Cores: 1}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "processor":
			if c.CPU >= 0 {
				err := finishCPUInfo(&c)
				if err != nil {
					return ret, err
				}
				ret = append(ret, c)
			}
			c = InfoStat{Cores: 1}
			t, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, err
			}
			c.CPU = int32(t)

		case "vendorId", "vendor_id":
			c.VendorID = value

		case "cpu family":
			c.Family = value

		case "model":
			c.Model = value

		case "model name", "cpu":
			c.ModelName = value
			if strings.Contains(value, "POWER8") ||
				strings.Contains(value, "POWER7") {
				c.Model = strings.Split(value, " ")[0]
				c.Family = "POWER"
				c.VendorID = "IBM"
			}

		case "stepping", "revision":
			val := value
			if key == "revision" {
				val = strings.Split(value, ".")[0]
			}
			t, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return ret, err
			}
			c.Stepping = int32(t)

		case "cpu MHz", "clock":
			// treat this as the fallback value, thus we ignore error
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				c.Mhz = t
			}

		case "cache size":
			t, err := strconv.ParseInt(strings.Replace(value, " KB", "", 1), 10, 64)
			if err != nil {
				return ret, err
			}
			c.CacheSize = int32(t)

		case "physical id":
			c.PhysicalID = value

		case "core id":
			c.CoreID = value

		case "flags", "Features":
			c.Flags = strings.FieldsFunc(value, func(r rune) bool {
				return r == ',' || r == ' '
			})
		}
	}

	if c.CPU >= 0 {
		err := finishCPUInfo(&c)
		if err != nil {
			return ret, err
		}
		ret = append(ret, c)
	}

	return ret, nil
}

// NumCPU returns the count of CPUs in the CPU affinity mask of the pid 1 process.
func NumCPU() (int, error) {
	cpuInfos, err := Info()
	if err != nil {
		return 0, err
	}
	var count int32
	for _, inf := range cpuInfos {
		count += inf.Cores
	}
	return int(count), nil
}

func sysCPUPath(cpu int32, relPath string) string {
	return osx.HostSys(fmt.Sprintf("devices/system/cpu/cpu%d", cpu), relPath)
}

func finishCPUInfo(c *InfoStat) error {
	var lines []string
	var err error
	var value float64

	if len(c.CoreID) == 0 {
		lines, err = osx.ReadLines(sysCPUPath(c.CPU, "topology/core_id"))
		if err == nil {
			c.CoreID = lines[0]
		}
	}

	// override the value of c.Mhz with cpufreq/cpuinfo_max_freq regardless
	// of the value from /proc/cpuinfo because we want to report the maximum
	// clock-speed of the CPU for c.Mhz, matching the behaviour of Windows
	lines, err = osx.ReadLines(sysCPUPath(c.CPU, "cpufreq/cpuinfo_max_freq"))
	// if we encounter errors below such as there are no cpuinfo_max_freq file,
	// we just ignore. so let Mhz is 0.
	if err != nil {
		return nil
	}
	value, err = strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil
	}
	c.Mhz = value / 1000.0 // value is in kHz
	if c.Mhz > 9999 {
		c.Mhz = c.Mhz / 1000.0 // value in Hz
	}
	return nil
}
