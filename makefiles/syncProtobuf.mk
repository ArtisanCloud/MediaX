# 源 proto 文件目录
SRC_DIR := manifest/protobuf/proto/

# 目标仓库路径
TARGET_REPO := ../MediaXProtobuf
TARGET_DIR := $(TARGET_REPO)/proto

# 创建目标目录并复制 proto 文件
copy_protos:
	@echo "创建目标目录: $(TARGET_DIR)"
	mkdir -p $(TARGET_DIR)
	@echo "复制 proto 文件到 MediaXProtobuf..."
	cp -rf $(SRC_DIR)/* $(TARGET_DIR)
	@echo "复制完成 ✅"

