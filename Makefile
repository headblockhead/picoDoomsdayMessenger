edit:
	tinygo-edit --target pico --editor code
deploy:
	tinygo flash -target=pico pico/main.go
runlocal:
	go run local/main.go
