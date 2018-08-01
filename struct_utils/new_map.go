package struct_utils

func NewMapRef() map[string]ReflectFieldMap{
	return make(map[string]ReflectFieldMap)
}

func NewMapMapString() map[string]map[string]string{
	return make(map[string]map[string]string)
}
