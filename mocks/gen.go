package mocks

//go:generate go-mockgen -f github.com/go-nacelle/config -i Config -o config.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Logger -o logger.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Sourcer -o sourcer.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i TagModifier -o tag_modifier.go
