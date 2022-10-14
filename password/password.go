package password

// Generator is an interface that implements the Generate function. This
// is useful for testing where you can pass this interface instead of a real
// password generator to mock responses for predictability.
type Generator interface {
	MustGenerate(password string) string
	Generate(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

// defaultGenerator may change the underlying algorithm in the future
var defaultGenerator = New()

func MustGenerate(password string) string      { return defaultGenerator.MustGenerate(password) }
func Generate(password string) (string, error) { return defaultGenerator.Generate(password) }
func Compare(hashedPassword, password string) bool {
	return defaultGenerator.Compare(hashedPassword, password)
}
