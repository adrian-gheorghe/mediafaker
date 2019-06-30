base:
	go run $$(ls -1 *.go | grep -v _test.go)
local:
	go run $$(ls -1 *.go | grep -v _test.go) local --source $(SOURCEPATH) --destination $(DESTPATH) --extcopy=".htaccess" --extcopy=".gitignore" --jsonlog=$(JSONLOG)
url:
	go run $$(ls -1 *.go | grep -v _test.go) url --extcopy=".htaccess" --extcopy=".gitignore" --destination "$(DESTPATH)" --source "$(URL)" --jsonlog=$(JSONLOG)

ssh:
	go run $$(ls -1 *.go | grep -v _test.go) ssh --extcopy=".htaccess" --extcopy=".gitignore" --destination $(DESTPATH) --source "$(SOURCE)" --ssh-host "$(SSHHOST)" --ssh-user "$(SSHUSER)" --ssh-key "$(SSHKEY)" --jsonlog=$(JSONLOG)

ver:
	go run $$(ls -1 *.go | grep -v _test.go) --version

test:
	go test -v ./...