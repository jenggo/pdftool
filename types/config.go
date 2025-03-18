package types

import "github.com/gofiber/storage/minio"

const (
	AppName          string = "pdfTool"
	AppVersion       string = "0.0.1-encryptOnly"
	MistralOcrApiUrl string = "https://api.mistral.ai/v1/ocr"
)

var Config struct {
	App struct {
		Listen     string `yaml:"listen" env:"LISTEN" env-default:":2804"`
		PPROF      string `yaml:"pprof" env:"PPROF"`
		LogLevel   int8   `yaml:"log_level" env:"LOG_LEVEL" env-default:"2"` // 0: debug, 1: info, 2: warning, 3: error, 4: fatal, 5: panic
		Cloudflare bool   `yaml:"cloudflare" env:"CLOUDFLARE" env-default:"true"`
		Sentry     string `yaml:"sentry" env:"SENTRY"`
		BaseURL    string `yaml:"base_url" env:"BASE_URL" env-default:"http://localhost:2804"`
		Auth       struct {
			User string `yaml:"user" env:"AUTH_USER" env-default:"lorem"`
			Pass string `yaml:"pass" env:"AUTH_PASS" env-default:"ipsumDOLORSITamet"`
		} `yaml:"auth"`
	} `yaml:"app"`

	Keys struct {
		API     string `yaml:"api_key" env:"API_KEY"`
		Mistral string `yaml:"mistral" env:"MISTRAL"`
	} `yaml:"keys"`

	Swagger struct {
		Enable bool   `yaml:"enable" env:"SWAGGER_ENABLE" env-default:"false"`
		Path   string `yaml:"path" env:"SWAGGER_PATH" env-default:"/use"`
	} `yaml:"swagger"`

	S3 struct {
		Enable   bool   `yaml:"enable" env:"S3_ENABLE" env-default:"false"`
		Endpoint string `yaml:"endpoint" env:"S3_ENDPOINT"`
		Bucket   string `yaml:"bucket" env:"S3_BUCKET"`
		Key      struct {
			Access string `yaml:"access" env:"S3_ACCESS"`
			Secret string `yaml:"secret" env:"S3_SECRET"`
		} `yaml:"key"`
		Storage *minio.Storage
	} `yaml:"s3"`
}
