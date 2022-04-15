package password

// Generator is an interface that implements the Generate function. This
// is useful for testing where you can pass this interface instead of a real
// password generator to mock responses for predictability.
type Generator interface {
	MustGenerate(password string) string
	Generate(password string) (string, error)
	Compare(hashedPassword, password string) bool
}
