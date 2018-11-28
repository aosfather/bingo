package bingo

type applicationConfig interface {
	Load(path string)
	GetProperty(key string) string
	GetPropertyAsInt(key string) int
	GetPropertyAsBool(key string) bool
	GetPropertyAsList(key string) []string
}
