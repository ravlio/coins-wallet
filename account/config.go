package account

type Config struct {
	PostgresURL string `yaml:"postgresUrl"`
	Kafka       struct {
		GroupID         string   `yaml:"groupId"`
		Topic           string   `yaml:"topic"`
		DeadLetterTopic string   `yaml:"deadLetterTopic"`
		Brokers         []string `yaml:"brokers"`
	}
	MetricsPort int `yaml:"metricsPort"`
}
