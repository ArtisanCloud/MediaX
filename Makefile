ROOT_DIR    = $(shell pwd)

include $(ROOT_DIR)/makefiles/genDoc.mk
include $(ROOT_DIR)/makefiles/syncProtobuf.mk

.PHONY: sessiontoken
sessiontoken:
	@echo "[MediaX] 启动 SessionToken 服务 (默认监听 :7070)"
	GO111MODULE=on go run ./cmd/sessiontoken

.PHONY: sessiontoken-bootstrap
sessiontoken-bootstrap:
	@echo "[MediaX] 初始化 SessionToken 配置"
	./scripts/sessiontoken-bootstrap.sh
