# 源 proto 文件目录
SRC_DIR := manifest/protobuf/proto/

# 目标仓库路径
TARGET_REPO := ../MediaXProtobuf
TARGET_DIR := $(TARGET_REPO)/proto

# 创建目标目录并复制 proto 文件
copy_protos:
	go run scripts/syncProtobuf/main.go -src=$(SRC_DIR) -dst=$(TARGET_DIR)

