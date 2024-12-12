module github.com/muesli/termenv/examples/ssh

go 1.18

require (
	github.com/charmbracelet/ssh v0.0.0-20221117183211-483d43d97103
	github.com/charmbracelet/wish v1.0.0
	github.com/creack/pty v1.1.18
	github.com/muesli/termenv v0.13.0
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/caarlos0/sshmarshal v0.1.0 // indirect
	github.com/charmbracelet/keygen v0.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/muesli/termenv => ../../
