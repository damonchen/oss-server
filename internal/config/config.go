package config

// Aliyun aliyun oss
type Aliyun struct {
	Name   string `yaml:"name"`
	ApiID  string `yaml:"apiId"`
	ApiKey string `yaml:"apiKey"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

// Tencent tencent oss
type Tencent struct {
	Name   string `yaml:"name"`
	ApiID  string `yaml:"apiId"`
	ApiKey string `yaml:"apiKey"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
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
	Host      string    `yaml:"host,omitempty"`
	Port      string    `yaml:"port,omitempty"`
	Providers []string  `yaml:"providers"`
	Aliyun    []Aliyun  `yaml:"aliyun,omitempty"`
	Tencent   []Tencent `yaml:"tencent,omitempty"`
	Auth      Auth      `yaml:"auth,omitempty"`
}
