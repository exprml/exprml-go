
.PHONY: protobuf
protobuf:
	rm -rf protobuf
	git clone https://github.com/exprml/exprml-api.git protobuf
	cp buf.gen.yaml protobuf/buf.gen.yaml
	rm -rf protobuf/.git
	cd protobuf && buf generate

.PHONY: schema
schema:
	curl -o schema.json https://raw.githubusercontent.com/exprml/exprml-language/refs/heads/main/schema.json

.PHONY: testdata
testdata:
	rm -rf testdata
	rm -rf exprml-testsuite
	git clone https://github.com/exprml/exprml-testsuite.git
	mv exprml-testsuite/testdata testdata
	rm -rf exprml-testsuite
