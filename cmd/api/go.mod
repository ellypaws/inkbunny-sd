module github.com/ellypaws/inkbunny-sd/cmd/api

go 1.22

replace github.com/ellypaws/inkbunny-sd => ../../

replace github.com/ellypaws/inkbunny-sd/cmd => ../

require (
	github.com/ellypaws/inkbunny-sd v0.0.0-00010101000000-000000000000
	github.com/ellypaws/inkbunny-sd/cmd v0.0.0-00010101000000-000000000000
	github.com/ellypaws/inkbunny/api v0.0.0-20240411110242-d491ced97f23
	github.com/go-errors/errors v1.5.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
