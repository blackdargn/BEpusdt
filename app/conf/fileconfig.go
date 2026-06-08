package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	pkgaes "github.com/v03413/bepusdt/pkg/aes"
)

type FileConfig struct {
	Listen        string `yaml:"listen"`
	Log           string `yaml:"log"`
	SQLite        string `yaml:"sqlite"`
	MySQLDSN      string `yaml:"mysql_dsn"`
	PostgreSQLDSN string `yaml:"postgresql_dsn"`
}

var fileConfig *FileConfig

func GetFileConfig() *FileConfig {
	return fileConfig
}

func LoadFileConfig(path string) error {
	candidates := configCandidates(path)
	for _, candidate := range candidates {
		cfg, err := loadConfigFile(candidate)
		if err != nil {
			return err
		}
		if cfg != nil {
			fileConfig = cfg

			return nil
		}
	}

	return nil
}

func configCandidates(path string) []string {
	if path != "" {
		return []string{path}
	}

	return DefaultConfigCandidates()
}

func loadConfigFile(path string) (*FileConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}

	content := strings.TrimSpace(string(raw))
	if content == "" {
		return nil, fmt.Errorf("配置文件 %s 为空", path)
	}

	if strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".enc") {
		content, err = pkgaes.Decrypt(content, pkgaes.ConfigCipherKey)
		if err != nil {
			return nil, fmt.Errorf("解密配置文件 %s 失败: %w", path, err)
		}
	}

	var cfg FileConfig
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
	}

	return &cfg, nil
}
