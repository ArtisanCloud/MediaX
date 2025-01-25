# 添加子模块
.PHONY: submodule.add
submodule.add:
	@git submodule add $(repo) $(path)
	@echo "子模块已添加: $(repo) 到 $(path)"

# 更新并拉取指定子模块最新代码
.PHONY: submodule.update
submodule.update:
	@git submodule update --remote --recursive $(module)
	@echo "子模块 $(module) 更新完成"

# 切换子模块版本
.PHONY: submodule.checkout
submodule.checkout:
	@cd $(module) && git checkout $(commit)
	@echo "子模块 $(module) 已切换到版本 $(commit)"