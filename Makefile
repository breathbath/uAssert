default: build

go/ap.pb:
	cd projects/accessProxy/protos && \
	protoc -I . \
	--include_imports --include_source_info \
	--descriptor_set_out=accessProxy.pb \
	--go_out=plugins=grpc:. \
	access_proxy/*.proto

# Protobuf compiler helper functions
protoc_check:
ifeq ("", "$(shell which protoc)")
	@echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	@echo "It looks like you don't have a version of protocol buffer tools."
	@echo "To install the protocol buffer toolchain on Linux, you can run:"
	@echo "    make install-protoc"
	@echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	exit 1
endif
ifneq ("libprotoc 3.7.1", "$(shell protoc --version)")
	@echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	@echo "You have the wrong version of protoc installed"
	@echo "Please install version 3.7.0"
	@echo "See https://github.com/golang/protobuf"
	@echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	exit 1
endif

ensure-simlink:
	[ -L projects/AccessProxy/protos/voltha_protos ] || (ln -sf /Users/breathbath/go/pkg/mod/github.com/opencord/voltha-protos@v0.0.0-20190711063307-f98ca1386c16/protos/voltha_protos projects/AccessProxy/protos/voltha_protos)
	[ -L projects/AccessProxy/protos/google ] || (ln -sf /Users/breathbath/go/pkg/mod/github.com/opencord/voltha-protos@v0.0.0-20190711063307-f98ca1386c16/protos/google projects/AccessProxy/protos/google)

build: protoc_check ensure-simlink go/ap.pb