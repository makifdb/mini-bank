package config

type Config struct {
	Port         int    `env:"PORT,default=8080"`
	JWTSecret    string `env:"JWT_SECRET,default=secret"`
	DBURL        string `env:"DB_URL,default=postgres://postgres:password@localhost:5432/postgres?sslmode=disable"`
	RedisAddr    string `env:"REDIS_ADDR,default=localhost:6379"`
	SmtpServer   string `env:"SMTP_SERVER,default=smtp.gmail.com"`
	SmtpPort     int    `env:"SMTP_PORT,default=587"`
	SmtpUser     string `env:"SMTP_USER,default=user"`
	SmtpPassword string `env:"SMTP_PASSWORD,default=password"`
	Environment  string `env:"ENVIRONMENT,default=development"`
}
