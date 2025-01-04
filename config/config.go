package config

import (
	"errors"
	"github.com/goccy/go-yaml"
	"io/fs"
	"os"
	"path/filepath"
)

type Config struct {
	Interpreter InterpreterConfig
	Debugger    DebuggerConfig
}

func GetConfig(fileOverride string, overrides map[string]string) (*Config, error) {
	var path string
	if fileOverride != "" {
		path = fileOverride
	} else {
		p, err := getPropFilePath()
		if err != nil {
			return nil, err
		}
		path = p
	}
	config, err := getConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	err = overridePropsFromEnv(config)
	if err != nil {
		return nil, err
	}
	err = overridePropsFromMap(config, overrides)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func DefaultConfig() Config {
	return Config{
		Interpreter: InterpreterConfig{
			DivideByZeroBehaviour:       Div0PromptForInput,
			ModulusByZeroBehaviour:      Div0PromptForInput,
			PutOutOfBoundsBehaviour:     OobNoOp,
			GetOutOfBoundsBehaviour:     OobZero,
			EnforceTorusSizeRestriction: false,
			TorusSizeRestrictionWidth:   80,
			TorusSizeRestrictionHeight:  25,
		},
		Debugger: DebuggerConfig{
			ShowTorus:            true,
			ShowTorusCoordinates: true,
			ShowStack:            true,
			EnableColors:         true,
		},
	}
}

func getConfigFromFile(path string) (*Config, error) {
	props := DefaultConfig()
	file, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) { // file not found; no config exists just use defaults
			return &props, nil
		}
		return nil, err
	}
	err = yaml.Unmarshal(file, &props)
	if err != nil {
		return nil, err
	}
	return &props, nil
}

func getPropFilePath() (string, error) {
	fromEnv := os.Getenv("KGF_CONFIG_PATH")
	if fromEnv != "" {
		return fromEnv, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kgf", "config.yaml"), nil
}

func oobMapper(s string) (OutOfBoundsBehaviour, error) {
	behaviour := outOfBoundsBehaviours[s]
	if behaviour == "" {
		return "", errors.New("Unknown out of bounds behaviour " + s)
	}
	return behaviour, nil
}

func div0Mapper(s string) (DivideByZeroBehaviour, error) {
	behaviour := divideByZeroBehaviours[s]
	if behaviour == "" {
		return "", errors.New("Unknown division by zero behaviour " + s)
	}
	return behaviour, nil
}
