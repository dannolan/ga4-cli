package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	LegacyDirName = ".ga4-cli"
	ConfigDirName = ".config/ga4-cli"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	PropertyID   string `json:"property_id,omitempty"`
}

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, ConfigDirName)
	}
	return &Store{dir: dir}, nil
}

func LegacyTokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, LegacyDirName, "token.json"), nil
}

func (s *Store) Dir() string        { return s.dir }
func (s *Store) envPath() string    { return filepath.Join(s.dir, "env") }
func (s *Store) configPath() string { return filepath.Join(s.dir, "config.json") }
func (s *Store) tokenPath() string  { return filepath.Join(s.dir, "token.json") }

func (s *Store) TokenPath() string { return s.tokenPath() }

func (s *Store) Load() (Config, error) {
	cfg := Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		PropertyID:   firstNonEmpty(os.Getenv("GA4_PROPERTY_ID"), os.Getenv("PROPERTY_ID")),
	}

	if fileCfg, err := s.loadJSON(); err == nil {
		cfg.ClientID = firstNonEmpty(cfg.ClientID, fileCfg.ClientID)
		cfg.ClientSecret = firstNonEmpty(cfg.ClientSecret, fileCfg.ClientSecret)
		cfg.PropertyID = firstNonEmpty(cfg.PropertyID, fileCfg.PropertyID)
	} else if !errors.Is(err, os.ErrNotExist) {
		return cfg, err
	}

	if envCfg, err := loadEnvFile(s.envPath()); err == nil {
		cfg.ClientID = firstNonEmpty(cfg.ClientID, envCfg.ClientID)
		cfg.ClientSecret = firstNonEmpty(cfg.ClientSecret, envCfg.ClientSecret)
		cfg.PropertyID = firstNonEmpty(cfg.PropertyID, envCfg.PropertyID)
	} else if !errors.Is(err, os.ErrNotExist) {
		return cfg, err
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return cfg, fmt.Errorf("missing CLIENT_ID or CLIENT_SECRET; set env vars or add %s", s.envPath())
	}
	return cfg, nil
}

func (s *Store) Save(cfg Config) error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.configPath(), append(data, '\n'), 0o600)
}

func (s *Store) loadJSON() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(s.configPath())
	if err != nil {
		return cfg, err
	}
	return cfg, json.Unmarshal(data, &cfg)
}

func loadEnvFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	if err := scanner.Err(); err != nil {
		return Config{}, err
	}
	return Config{
		ClientID:     values["CLIENT_ID"],
		ClientSecret: values["CLIENT_SECRET"],
		PropertyID:   firstNonEmpty(values["GA4_PROPERTY_ID"], values["PROPERTY_ID"]),
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
