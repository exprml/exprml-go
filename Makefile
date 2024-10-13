
.PHONY: protobuf
protobuf: ## make protobuf BRANCH=main
	rm -rf protobuf
	git clone -b ${BRANCH} https://github.com/exprml/exprml-api.git protobuf
	cp buf.gen.yaml protobuf/buf.gen.yaml
	rm -rf protobuf/.git
	cd protobuf && buf generate

.PHONY: schema
schema: ## make schema BRANCH=main
	curl -o schema.json https://raw.githubusercontent.com/exprml/exprml-language/refs/heads/${BRANCH}/schema.json

.PHONY: testdata
testdata: ## make testdata BRANCH=main
	rm -rf testdata
	rm -rf exprml-testsuite
	git clone -b ${BRANCH} https://github.com/exprml/exprml-testsuite.git
	mv exprml-testsuite/testdata testdata
	rm -rf exprml-testsuite
