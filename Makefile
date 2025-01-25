# Makefile: 主要文件，包含其他文件
.PHONY: help submodule build plugins

# 目标是显示帮助信息
help:
	@echo "Usage:"
	@echo "  make submodule    # 拉取子模块"
	@echo "  make build        # 构建插件"
	@echo "  make plugins      # 拉取并构建插件"

# 引入插件相关的子文件
include ./Makefiles/plugin/submodule.mk
include ./Makefiles/plugin/build.mk
