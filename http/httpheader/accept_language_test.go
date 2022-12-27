package httpheader

import (
	"fmt"
	"testing"
)

// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5

func TestAcceptLanguage_Each(t *testing.T) {
	al := AcceptLanguage("fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5")
	al.Each(func(v QualityValue) bool {
		fmt.Println(string(v.Value), v.Priority)
		return true
	})
}
