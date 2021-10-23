PROJECT_NAME := melt
IMAGE_NAME := melt_env
CONTAINER_NAME := melt_dev
IMAGE_ID := $$(docker image ls --filter=reference=$(IMAGE_NAME) --quiet)
CONTAINER_ID := $$(docker container ls --filter=name=$(CONTAINER_NAME) --quiet --all)


build: ## Dockerfileをビルド
	docker build --tag $(IMAGE_NAME) .docker/
rmi: ## イメージを破棄
	docker image rm -f $(IMAGE_NAME)
	docker image prune -f

run: ## コンテナを起動
	@if [ ! -n "$(IMAGE_ID)" ]; then make build; fi
	@if [ ! -n "$(CONTAINER_ID)" ]; then docker run -itd -v $(PWD):/work/melt --name $(CONTAINER_NAME) $(IMAGE_NAME) bash; fi
stop: ## コンテナを停止
	docker container stop $(CONTAINER_NAME)
start: ## コンテナを起動
	docker container start $(CONTAINER_NAME)
delete: ## コンテナを破棄
	@make stop
	docker container rm $(CONTAINER_NAME)
restart: ## コンテナを再起動
	@make delete
	@make run

destoroy: ## 環境を破棄
	@make delete
	@make rmi
remake: ## 環境の再作成
	@make destoroy
	@make run

attach: ## コンテナに接続(Ctrl-p,qでコンテナから抜ける)
	@if [ -n "$(CONTAINER_ID)" ]; then docker container start $(CONTAINER_ID); docker container exec -it $(CONTAINER_ID) bash; else echo "Error: No such container: $(CONTAINER_NAME)"; fi

init: ## 初期設定
	docker container exec -t $(CONTAINER_NAME) go mod init "github.com/schnell3526/"$(PROJECT_NAME)
	docker container exec -t $(CONTAINER_NAME) go get -d golang.org/x/tools/gopls
	docker container exec -t $(CONTAINER_NAME) go get -d github.com/ramya-rao-a/go-outline
	docker container exec -t $(CONTAINER_NAME) go get -d golang.org/x/tools/cmd/goimports


ls: ## コンテナの一覧
	docker container ls -a
ils: ## イメージ一覧
	docker image ls -a

echo:
	@if [ -n "$(IMAGE_ID)" ]; then echo "melt exists"; else echo "melt not exists"; fi

# help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'