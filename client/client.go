package client

type Config struct {
	BaseUrl     string
	ErrorLevel  string
	ServiceName string
	SessionName string
}

func (c *Config) Error() {

}

func NewClient(url, errLevel, svcName, sessName string) *Config {
	return &Config{
		BaseUrl:     url,
		ErrorLevel:  errLevel,
		ServiceName: svcName,
		SessionName: sessName,
	}
}
