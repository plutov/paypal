oapi:
	git submodule update --init --remote
	for file in openapi-specs/openapi/*.json; do \
		name=$$(basename $$file .json); \
		go tool oapi-codegen -o mockserver/$$name/gen.go --package srv --generate std-http-server,models,strict-server $$file; \
	done

