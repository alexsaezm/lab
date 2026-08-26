TARGET ?= $(notdir $(CURDIR))

.PHONY: build

build:
	mkdir -p "$(BUILD_DIR)"
	go build -o "$(BUILD_DIR)/$(TARGET)" .
	go build -gcflags="all=-N -l" -o "$(BUILD_DIR)/$(TARGET)_no_opt" .
