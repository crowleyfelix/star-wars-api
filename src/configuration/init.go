package configuration

var (
	config *Configuration
)

func init() {
	load()
}
