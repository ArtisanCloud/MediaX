# 定义文件夹和输出文件
PKG_DIR := pkg/client
DOCS_DIR := docs

# 生成文档的目标
gen_all_docs: generate_docs

# 遍历 pkg/client 目录，生成每个包的文档
generate_docs:
	 go run scripts/genDocs.go -target="pkg/client" -docs="../MediaXDoc/docs/mediax"
