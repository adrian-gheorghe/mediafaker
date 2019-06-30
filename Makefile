local:
	go run $$(ls -1 *.go | grep -v _test.go) local --source $(SOURCEPATH) --destination $(DESTPATH) --extcopy=".htaccess" --extcopy=".gitignore"
tree-url:
	go run $$(ls -1 *.go | grep -v _test.go) url --extcopy=".htaccess" --extcopy=".gitignore" --destination "$(DESTPATH)" --source "$(URL)"