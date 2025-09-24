SERVICES := auth store
GO_OUT := go
TS_OUT := ts
SWAGGER_CLI := swagger-codegen-cli-2.4.46.jar

gen: clean gen-go gen-ts

clean:
	@rm -rf $(GO_OUT)/ $(TS_OUT)/

gen-go:
	@for s in $(SERVICES); do \
		echo "Generating Go client for $$s..."; \
		java -jar $(SWAGGER_CLI) generate -i ../$$s/docs/swagger.json -l go -o $(GO_OUT)/$$s; \
	done

gen-ts:
	@for s in $(SERVICES); do \
		echo "Generating TypeScript client for $$s..."; \
		java -jar $(SWAGGER_CLI) generate -i ../$$s/docs/swagger.json -l typescript-node -o $(TS_OUT)/$$s; \
	done
