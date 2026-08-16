.DEFAULT_GOAL := build

BUILD_DIR := build
EXPERIMENTS := $(notdir $(patsubst %/Makefile,%,$(wildcard */Makefile)))

.PHONY: build bootstrap clean help list $(EXPERIMENTS)

# For this to work, every directory must have its own Makefile.
build: $(EXPERIMENTS)

# Keep the workspace in sync before compiling any experiment.
bootstrap:
	go work use -r .

# Build one experiment with: make <experiment-name>
$(EXPERIMENTS): | bootstrap
$(EXPERIMENTS):
	$(MAKE) -C "$@" build BUILD_DIR="$(abspath $(BUILD_DIR))"

clean:
	$(RM) -r "$(BUILD_DIR)"

# List should print every valid directory.
list:
	@for experiment in $(EXPERIMENTS); do printf '%s\n' "$$experiment"; done
