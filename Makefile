build:
	go build ./cmd/wykobot

run: build
	./wykobot

build-linux:
	GOOS=linux go build ./cmd/wykobot

deploy: build-linux
	scp wykobot "${DEPLOY_HOST}:/usr/bin"
	ssh "${DEPLOY_HOST}" "mkdir -p /etc/wykobot/"
	scp config.yaml "${DEPLOY_HOST}:/etc/wykobot/config.yaml"
	scp systemd/wykobot.service "${DEPLOY_HOST}:/etc/systemd/system/"
	scp systemd/wykobot.timer "${DEPLOY_HOST}:/etc/systemd/system/"
	ssh "${DEPLOY_HOST}" "systemctl daemon-reload"
	ssh "${DEPLOY_HOST}" "systemctl start wykobot"