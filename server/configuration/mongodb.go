package configuration

//MongoDB configurations
type MongoDB struct {
	URI         string `cfgDefault:"mongodb://localhost:27017"`
	Database    string `cfgDefault:"test"`
	MaxPoolSize int    `cfgDefault:"2048"`
}
