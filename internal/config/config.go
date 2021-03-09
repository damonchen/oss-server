package config

// AliOSS aliyun oss
type AliOSS struct {
	APIKey     string `yaml:"apiKey"`
	APISecret  string `yaml:"apiSecret"`
	BucketName string `yaml:"bucket"`
	Region     string `yaml:"region"`
}

// TencentOSS tencent oss
type TencentOSS struct {
	SecretID   string `yaml:"apiId"`
	SecretKey  string `yaml:"apiKey"`
	BucketName string `yaml:"bucket"`
	Region     string `yaml:"region"`
}

// Auth auth for request
type ProxyAuth struct {
	// 认证服务器地址，http或https协议
	AuthPath string `yaml:"authPath"`
}

// BasicAuth basic auth
type BasicAuth struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

// Auth auth
type Auth struct {
	// Type, value may basic or proxy
	Type      string    `yaml:"type,omitempty"`
	BasicAuth BasicAuth `yaml:"auth,omitempty"`
	ProxyAuth ProxyAuth `yaml:"proxy,omitempty"`
}

// Configuration configuration
type Configuration struct {
	Port       string     `yaml:"port"`
	Providers  []string   `yaml:"providers"`
	AliOSS     AliOSS     `yaml:"aliyun,omitempty"`
	TencentOSS TencentOSS `yaml:"tencent,omitempty"`
	Auth       Auth       `yaml:"auth,omitempty"`
}
