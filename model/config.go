package model

type Config struct {
	Mqtt Mqtt `yaml:"mqtt"`
}

type Mqtt struct {
	Server      string `yaml:"server"`
	Port        int64  `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	ClientId    string `yaml:"client_id"`
	TopicPrefix string `yaml:"topic_prefix"`
}
