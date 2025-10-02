SERVICES := auth store fstorage
GO_OUT := go
TS_OUT := ts
SWAGGER_CLI := swagger-codegen-cli-2.4.46.jar

gen: clean gen-go gen-ts

clean:
	@for s in $(SERVICES); do \
		if [ "$$s" = "fstorage" ]; then \
			echo "Skipping fstorage - manually modified client"; \
		else \
			echo "Cleaning Go client for $$s..."; \
			rm -rf $(GO_OUT)/$$s; \
		fi; \
	done

gen-go:
	@for s in $(SERVICES); do \
		echo "Generating Go client for $$s..."; \
		if [ "$$s" = "fstorage" ]; then \
			echo "Skipping fstorage - manually modified client"; \
		else \
			java -jar $(SWAGGER_CLI) generate -i ../$$s/docs/swagger.json -l go -o $(GO_OUT)/$$s; \
		fi; \
	done

gen-ts:
	@for s in $(SERVICES); do \
		echo "Generating TypeScript client for $$s..."; \
		java -jar $(SWAGGER_CLI) generate -i ../$$s/docs/swagger.json -l typescript-node -o $(TS_OUT)/$$s; \
	done
