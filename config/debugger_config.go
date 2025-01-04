package config

type DebuggerConfig struct {
	ShowTorus            bool `yaml:"show-torus"`
	ShowTorusCoordinates bool `yaml:"show-torus-coordinates"`
	ShowStack            bool `yaml:"show-stack"`
	EnableColors         bool `yaml:"enable-colors"`
}
