package cmdx

func EnvVarExamplesHelpMessage(name string) string {
	return `This command exposes a variety of controls via environment variables. Here are some examples on how to
configure environment variables:

Linux / macOS:
	$ export FOO=bar
	$ export BAZ=bar
	$ ` + name + ` ...

	$ FOO=bar BAZ=bar ` + name + ` ...

Docker:
	$ docker run -e FOO=bar -e BAZ=bar ...

Windows (cmd):
	> set FOO=bar
	> set BAZ=bar
	> ` + name + ` ...

Windows (powershell):
	> $env:FOO = "bar"
	> $env:BAZ = "bar"
	> ` + name + `
`
}
