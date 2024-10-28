package configloader

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(path string) (Config, error) {
	if file_info, stat_err := os.Stat(path); stat_err == nil {
		if file_info.IsDir() {
			return Config{}, fmt.Errorf("path %v is a directory, not a file", path)
		}
		if file, err := os.ReadFile(path); err == nil {
			var loaded_config Config

			if err := json.Unmarshal(file, &loaded_config); err == nil {
				return loaded_config, nil
			} else {
				return Config{}, err
			}

		} else {
			return Config{}, err
		}

	} else {
		return Config{}, fmt.Errorf("configuration file %v does not exist", path)
	}
}
