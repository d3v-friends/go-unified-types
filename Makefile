SHELL:=/bin/zsh

upmod:
	go get -u ./...
	go mod tidy
tag:
	sh script/tag.sh $(shell cat ./version)
gen:
	sh script/gen.sh utypes