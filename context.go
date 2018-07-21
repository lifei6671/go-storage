package storage

//上下文信息
type Context map[string]interface{}

func NewContext() Context {
	return make(map[string]interface{})
}
func (m Context) Has(key string) bool  {
	_,ok := m[key];
	return ok
}

func (m Context) Set(key string ,value interface{}) Context {
	m[key] = value
	return m
}

func (m Context) Get(key string) (interface{} ,bool) {
	v,ok := m[key];

	return v,ok
}