
.SECONDEXPANSION:
SITE_DIR = ../freeconf-docs
SITE_EXAMPLE_DIR = $(SITE_DIR)/content/en/docs/Examples

SITE_DOCS_SRC = $(wildcard */*.gotmpl)
SITE_DOCS = $(foreach F,$(SITE_DOCS_SRC),$(SITE_EXAMPLE_DIR)/$(dir $(F))$(basename $(notdir $(F))))

# Generate the freeconf.org site example docs from templates from here
docs: $(SITE_DOCS) doc-images

$(SITE_DOCS) : $(SITE_EXAMPLE_DIR)/% : %.gotmpl ./site-docs/main.go $$(shell find $$(dir $$*) -path '*/venv' -prune -type f)
	test -d $(dir $@) || mkdir -p $(dir $@)
	go run ./site-docs/main.go $< > $@

# try to restrict depth otherwise dirs like ansible's venv dir creates a massive list of empty
# dirs that confuses the hell out of hugo
doc-images:
	rsync -av  --include '*/' --include '*.png' --exclude '*' --exclude '/*/*/*/'  ./ $(SITE_EXAMPLE_DIR)

test: test-go test-py

test-go:
	go test ./...

test-py:
	cd node-basic/python; \
		./test_manage.py
	cd node-extend/python; \
		./test_manage.py
