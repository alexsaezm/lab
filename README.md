# lab

This repository contains little experiments that popped out of the blue during
the coding session of another project. They might seem dumb or out of context,
but they probably made sense at that time.

To keep the repository clean, each directory is an experiment, and it should
contain a `Makefile` with a `build` target. The `Makefile` at the top-level will
run that target with `BUILD_DIR` set to the `build/` directory so it's easy to
ignore or to clean.

For example, the Go experiment in `hello_world_go/Makefile` contains:

```Makefile
TARGET ?= $(notdir $(CURDIR))

.PHONY: build

build:
	mkdir -p "$(BUILD_DIR)"
	go build -o "$(BUILD_DIR)/$(TARGET)" .
```

While this approach is pretty nice for building, it makes building from within a
directory a cumbersome experience:

```bash
make BUILD_DIR=../build
```
