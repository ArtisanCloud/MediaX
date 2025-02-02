build.plugin.all.bundle: build.plugin.bundle.wechat build.plugin.bundle.douYin build.plugin.bundle.redBook

WORKDIR := plugins/open

.PHONY: build.plugin.bundle.wechat
build.plugin.bundle.wechat:
	@echo "正在构建Wechat Bundle插件组..."
	cd $(WORKDIR)/MediaXPlugin-Wechat && go mod tidy && $(MAKE) plugin.all.build
	@echo "Wechat Bundle插件组构建完成"

.PHONY: build.plugin.bundle.douYin
build.plugin.bundle.douYin:
	@echo "正在构建DouYin插件组..."
	cd $(WORKDIR)/MediaXPlugin-DouYin && go mod tidy && $(MAKE) plugin.all.build
	@echo "DouYin插件组构建完成"

.PHONY: build.plugin.bundle.redBook
build.plugin.bundle.redBook:
	@echo "正在构建redBook插件组..."
	cd $(WORKDIR)/MediaXPlugin-RedBook && go mod tidy && $(MAKE) plugin.all.build
	@echo "redBook插件组构建完成"