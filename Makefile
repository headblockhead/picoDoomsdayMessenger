edit:
	tinygo-edit --target pico --editor code
deploy:
	tinygo flash -target=pico pico/main.go
runlocal:
	nix-shell local/shell.nix --command "go run local/main.go"
