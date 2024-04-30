module github.com/ellypaws/inkbunny-sd/cmd

go 1.22

replace github.com/ellypaws/inkbunny-sd => ../

require (
	github.com/ellypaws/inkbunny-sd v0.0.0-00010101000000-000000000000
	github.com/ellypaws/inkbunny/api v0.0.0-20240411110242-d491ced97f23
	github.com/go-errors/errors v1.5.1
	github.com/stretchr/testify v1.9.0
	golang.org/x/term v0.18.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
