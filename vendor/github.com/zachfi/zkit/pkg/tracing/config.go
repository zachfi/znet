package tracing

import "flag"

type Config struct {
	OtelEndpoint string `yaml:"otel_endpoint"`
	OrgID        string `yaml:"org_id"`
}

func (c *Config) RegisterFlagsAndApplyDefaults(prefix string, f *flag.FlagSet) {
	f.StringVar(&c.OtelEndpoint, "otel.endpoint", "", "otel endpoint, eg: tempo:4317")
	f.StringVar(&c.OrgID, "org.id", "", "org ID to use when sending traces")
}
