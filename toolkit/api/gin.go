package api

// GinTemplate 模板
func GinTemplate() string {
	outString := `
package {{.Pkg}}
`
	return outString
}
