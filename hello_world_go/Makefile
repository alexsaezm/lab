TARGET ?= $(notdir $(CURDIR))

.PHONY: build

build:
	mkdir -p "$(BUILD_DIR)"
	go build -o "$(BUILD_DIR)/$(TARGET)" .
