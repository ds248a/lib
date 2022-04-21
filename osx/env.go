package osx

import (
	"os"
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
