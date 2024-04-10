.PHONY:docker
docker:
	@rm webook ||true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -buildvcs=false -o webook .

	@docker rmi -f sushuangbao88/webook:v0.0.1
	@docker build -t sushuangbao88/webook:v0.0.1 .