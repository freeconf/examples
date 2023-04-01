
.SECONDEXPANSION:
SITE_DIR = ../freeconf-docs
SITE_EXAMPLE_DIR = $(SITE_DIR)/content/en/docs/Examples

SITE_DOCS_SRC = $(wildcard */index.md.gotmpl)
SITE_DOCS = $(foreach F,$(SITE_DOCS_SRC),$(SITE_EXAMPLE_DIR)/$(dir $(F))index.md)

GO_DIRS = \
	node-list

# Generate the freeconf.org site example docs from templates from here
docs: $(SITE_DOCS)

$(SITE_DOCS) : $(SITE_EXAMPLE_DIR)/%/index.md : %/index.md.gotmpl ./site-docs/main.go $(shell find $$* -type f)
	test -d $(dir $@) || mkdir -p $(dir $@)
	go run ./site-docs/main.go $< > $@

test:
	$(foreach D,$(GO_DIRS), cd $(D); go test . ;)