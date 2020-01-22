runjqbot: $(shell find . -name "*.go")
	go build -ldflags="-s -w" -o ./runjqbot

deploy: runjqbot
	ssh root@nusakan-58 'systemctl stop runjqbot'
	scp runjqbot nusakan-58:runjqbot/runjqbot
	ssh root@nusakan-58 'systemctl start runjqbot'
