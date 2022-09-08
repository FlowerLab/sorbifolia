package httprouter

type Params []Param

type Param struct {
	Key, Val string
}

func (p *Params) clean()              { *p = (*p)[:0] }
func (p *Params) Add(key, val string) { *p = append(*p, Param{key, val}) }
func (p *Params) Get(key string) (string, bool) {
	for _, v := range *p {
		if v.Key == key {
			return v.Val, true
		}
	}
	return "", false
}
