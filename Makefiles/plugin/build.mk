# 默认插件路径
module ?= $(wildcard plugins/open/*)

# 构建所有插件
.PHONY: plugin.build
plugin.build:
	@for dir in $(module); do \
		if [ -d $$dir ]; then \
			echo "正在构建插件 $$dir..."; \
			go build -o $$dir/plugin/plugin.so -buildmode=plugin $$dir/plugin/plugin.go; \
			echo "插件 $$dir 构建完成"; \
		fi \
	done