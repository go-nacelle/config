package config

//go:generate go-mockgen -f github.com/go-nacelle/config -i Logger -o mock_logger_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Sourcer -o mock_sourcer_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i FileSystem -o mock_fs_test.go
