package es

type ElasticSearch struct {
	URL string
	Port           int
	ApiKey         string
	Username       string
	Password       string
	CloudHost      string
	MaxConnections int32
}
