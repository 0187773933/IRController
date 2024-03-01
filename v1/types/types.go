package types

type Key struct {
	Name string `yaml:"name"`
	Code string `yaml:"code"`
	KeyPath string `yaml:"key_path"`
}

type Remote struct {
	Name string `yaml:"name"`
	Keys map[string]Key `yaml:"keys"`
}

type ConfigFile struct {
	KeySaveFileBasePath string `yaml:"key_save_file_base_path"`
	DefaultRemote string `yaml:"default_remote"`
	Remotes map[string]Remote `yaml:"remotes"`
}