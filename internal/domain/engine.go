package domain

type Engine string

const (
	EngineSQLServer Engine = "sqlserver"
	EnginePostgres  Engine = "postgres"
)

func ParseEngine(s string) (Engine, bool) {
	switch Engine(s) {
	case EngineSQLServer, EnginePostgres:
		return Engine(s), true
	default:
		return "", false
	}
}