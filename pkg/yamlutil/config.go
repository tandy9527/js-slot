package yamlutil

type ConfigGetter interface {
	ExtraMap() map[string]any
}
