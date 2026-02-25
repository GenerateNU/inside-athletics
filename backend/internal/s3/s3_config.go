package s3

import (
	"os"
	"strconv"
	"time"
)

// Env keys for S3 config.
const (
	EnvBucket    = "S3_BUCKET"
	EnvRegion    = "AWS_REGION"
	EnvExpirySec = "S3_PRESIGNED_EXPIRY_SECONDS"
)

// Returns config from env; ok is false if S3_BUCKET or AWS_REGION is missing.
func LoadConfigFromEnv() (cfg Config, ok bool) {
	cfg.Bucket = os.Getenv(EnvBucket)
	cfg.Region = os.Getenv(EnvRegion)
	if cfg.Bucket == "" || cfg.Region == "" {
		return cfg, false
	}
	if s := os.Getenv(EnvExpirySec); s != "" {
		if sec, err := strconv.Atoi(s); err == nil && sec > 0 {
			cfg.PresignedURLExpiry = time.Duration(sec) * time.Second
		}
	}
	if cfg.PresignedURLExpiry == 0 {
		cfg.PresignedURLExpiry = DefaultPresignedURLExpiry
	}
	return cfg, true
}
