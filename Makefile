edit:
	tinygo-edit --target challenger-rp2040 --editor code
deploy:
	tinygo flash -target challenger-rp2040 pico/main.go
runlocal:
	nix-shell local/shell.nix --command "go run local/main.go"
