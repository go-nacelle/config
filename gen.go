package config

//go:generate go-mockgen -f github.com/go-nacelle/config -i Config -o config_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Logger -o logger_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Sourcer -o sourcer_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i FileSystem -o fs_mock_test.go
