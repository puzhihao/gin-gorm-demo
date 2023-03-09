package conf

import "github.com/BurntSushi/toml"

func LoadConfigToml(path string) error {
	//new一个配置对象
	cfg := NewDefaultConfig()
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return err
	}
	SetGlobalConfig(cfg)
	return nil
}
