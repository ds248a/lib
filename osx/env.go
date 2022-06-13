package osx

import (
	"os"
	"path/filepath"
	"strings"
)

// UnsetEnv удаляет из окружения заданый параметр.
func UnsetEnv(key string) error {
	envs := os.Environ()
	os.Clearenv()
	for _, e := range envs {
		strs := strings.SplitN(e, "=", 2)
		if strs[0] == key {
			continue
		}
		if err := os.Setenv(strs[0], strs[1]); err != nil {
			return err
		}
	}
	return nil
}

// GetEnv ищет переменную окружения ко заданному ключу.
// Если переменная содержит пустое знацение, то возвращает 'default'.
func GetEnv(key string, dfault string, combineWith ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = dfault
	}

	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
}

func HostProc(combineWith ...string) string {
	return GetEnv("HOST_PROC", "/proc", combineWith...)
}

func HostSys(combineWith ...string) string {
	return GetEnv("HOST_SYS", "/sys", combineWith...)
}

func HostEtc(combineWith ...string) string {
	return GetEnv("HOST_ETC", "/etc", combineWith...)
}
