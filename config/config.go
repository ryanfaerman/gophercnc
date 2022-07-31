package config

type Config struct {
	Version int

	Tools []struct {
		Name    string
		Path    string
		Default bool
	}

	Machines []struct {
		Name    string
		Path    string
		Default bool
	}

	dirty bool
}

/*

tool:
  - name: mpcnc
	path: some/path/tools.tools
	default: true

  - name: banana
	path: something

machines:
  - name: MPCNC
	path: some/path/machine.machine

/**/
