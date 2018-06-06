package mongodb

var (
	Pool SessionManager
)

func init() {
	Pool = new(pool)
}
