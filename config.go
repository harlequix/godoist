package godoist

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	koanf "github.com/knadh/koanf/v2"

	ktoml "github.com/knadh/koanf/parsers/toml"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Token      string `koanf:"token"`
	ApiURL     string `koanf:"api_url"`
	Timeout    int    `koanf:"timeout"`
	Debug      bool   `koanf:"debug"`
	UseSyncAPI bool   `koanf:"use_sync_api"`
}

func (config Config) Merge(other *Config) {
	k := koanf.New(".")
	k.Load(structs.Provider(&config, "koanf"), nil)
	k.Load(structs.Provider(other, "koanf"), nil)
	k.Unmarshal("", &config)
}

func defaultConfig() *Config {
	return &Config{
		Token:      "",
		ApiURL:     "https://api.todoist.com/api/v1",
		Timeout:    30,
		Debug:      false,
		UseSyncAPI: false,
	}
}

func BuildConfig(files []string, envPrefix string, external interface{}) (*Config, error) {
	k := koanf.New(".")
	out := defaultConfig()
	k.Load(structs.Provider(out, "koanf"), nil)
	// Files (lowest → highest precedence in the order provided)
	for _, f := range files {
		if f == "" {
			continue
		}
		if _, err := os.Stat(f); err != nil {
			// Skip missing files quietly; change to return err if you want strict behavior.
			continue
		}
		parser, perr := parserFor(f)
		if perr != nil {
			return nil, perr
		}
		if err := k.Load(file.Provider(f), parser); err != nil {
			return nil, fmt.Errorf("load file %q: %w", f, err)
		}
	}

	// Env (highest precedence within base)
	xfm := func(s string) string {
		s = strings.TrimPrefix(s, envPrefix)
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "__", ".") // double underscore => dot
		return s
	}
	if err := k.Load(env.Provider(envPrefix, ".", xfm), nil); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}
	k.Print()
	// External (highest precedence overall)
	err := k.Load(structs.Provider(external, "koanf"), nil, koanf.WithMergeFunc(func(src, dest map[string]interface{}) error {
		for k, v := range src {
			if str, ok := v.(string); ok && str == "" {
				// Skip empty string, keep original
				continue
			}
			dest[k] = v
		}
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := k.Unmarshal("", &out); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	k.Print()
	return out, nil
}

func parserFor(path string) (koanf.Parser, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return kyaml.Parser(), nil
	case ".toml":
		return ktoml.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported config format for %q", path)
	}
}
