.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = qgnet

# 初始化
init:
	go mod init github.com/suisrc/${APP}

# 修正依赖
tidy:
	go mod tidy
