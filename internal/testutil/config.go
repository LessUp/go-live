package testutil

import "live-webrtc-go/internal/config"

func TestConfig() *config.Config {
	return &config.Config{
		HTTPAddr:          ":8080",
		AllowedOrigin:     "*",
		AuthToken:         "",
		STUN:              []string{"stun:stun.l.google.com:19302"},
		TURN:              []string{},
		TLSCertFile:       "",
		TLSKeyFile:        "",
		RecordEnabled:     false,
		RecordDir:         "records",
		MaxSubsPerRoom:    0,
		RoomTokens:        map[string]string{},
		TURNUsername:      "",
		TURNPassword:      "",
		UploadEnabled:     false,
		DeleteAfterUpload: false,
		S3Endpoint:        "",
		S3Region:          "",
		S3Bucket:          "",
		S3AccessKey:       "",
		S3SecretKey:       "",
		S3UseSSL:          true,
		S3PathStyle:       false,
		S3Prefix:          "",
		AdminToken:        "",
		RateLimitRPS:      0,
		RateLimitBurst:    0,
		JWTSecret:         "",
		JWTAudience:       "",
		PprofEnabled:      false,
		OTELServiceName:   "live-webrtc-go",
	}
}
