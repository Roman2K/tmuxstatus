all: main.cr
	crystal build --release $?

spec: *_spec.cr
	crystal spec $?
