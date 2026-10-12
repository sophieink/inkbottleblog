.PHONY: dev build deploy clean

# 本地预览：先构建，再起 server
dev:
	go run ./cmd/generator
	go run ./cmd/server

# 只构建
build:
	go run ./cmd/generator

# 一键发布
deploy:
	git add .
	git commit -m "publish: $$(date +%Y-%m-%d)" || true
	git push

# 清理产物
clean:
	rm -rf docs/
