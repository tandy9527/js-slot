package yamlutil

type Config[T ConfigGetter] struct {
	Cfg T
}

func (c Config[T]) GetString(path ...string) string { return GetNested[string](c.Cfg, path...) }
func (c Config[T]) GetInt(path ...string) int       { return GetNested[int](c.Cfg, path...) }
