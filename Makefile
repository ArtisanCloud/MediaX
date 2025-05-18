ROOT_DIR    = $(shell pwd)

include $(ROOT_DIR)/makefiles/genDoc.mk
include $(ROOT_DIR)/makefiles/syncProtobuf.mk