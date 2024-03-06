module github.com/ellypaws/inkbunny-sd/cmd

go 1.22

replace github.com/ellypaws/inkbunny-sd => ../

require (
	github.com/ellypaws/inkbunny-sd v0.0.0
	github.com/ellypaws/inkbunny/api v0.0.0-20240306094519-f8fede62380c
	golang.org/x/term v0.18.0
)

require golang.org/x/sys v0.18.0 // indirect
