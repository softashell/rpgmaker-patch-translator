package engine

type EngineType int

const (
	None EngineType = iota
	RPGMVX
	Wolf
)

var engine EngineType

func Set(e EngineType) {
	engine = e
}

func Get() EngineType {
	return engine
}

func Is(e EngineType) bool {
	return engine == e
}
