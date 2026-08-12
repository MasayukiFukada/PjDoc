package get_config

// AppConfig はクライアント（フロントエンド）に提供するアプリケーション設定です。
type AppConfig struct {
	PlantUMLServer string `json:"plantumlServer"`
}

// ConfigService は設定保持と提供を行うサービスです。
type ConfigService struct {
	config AppConfig
}

// NewConfigService は新しい ConfigService を生成します。
func NewConfigService(plantUMLServer string) *ConfigService {
	if plantUMLServer == "" {
		plantUMLServer = "https://kroki.io"
	}
	return &ConfigService{
		config: AppConfig{
			PlantUMLServer: plantUMLServer,
		},
	}
}

// GetConfig は現在の設定を返します。
func (s *ConfigService) GetConfig() AppConfig {
	return s.config
}
