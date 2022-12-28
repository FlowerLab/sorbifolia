package httpheader

type AcceptLanguage []byte

func (v AcceptLanguage) Each(fn EachQualityValue) { eachQualityValue(v, fn) }
