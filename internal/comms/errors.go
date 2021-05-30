package comms

import "errors"

var ErrNoHTTPConfig = errors.New("unable to create HTTP server with empty config")

var ErrNoHTTPListenAddress = errors.New("unable to create HTTP server with empty ListenAddress")

var ErrNilVaultConfig = errors.New("nil vault config")

var ErrMissingVaultConfig = errors.New("missing vault config")
