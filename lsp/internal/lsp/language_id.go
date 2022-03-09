package lsp

// LanguageID represents the coding language
// of a file
type LanguageID string

const (
	Terraform LanguageID = "demo"
	Tfvars    LanguageID = "demo-vars"
)

func (l LanguageID) String() string {
	return string(l)
}
