TARGET ?= $(notdir $(CURDIR))

.PHONY: build

build:
	mkdir -p "$(BUILD_DIR)"
	gcc -mavx512f test.c -o "$(BUILD_DIR)/$(TARGET)"
