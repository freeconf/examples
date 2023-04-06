
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

doc-images:
	rsync -av  --include '*/' --include '*.png' --exclude '*'  ./ $(SITE_EXAMPLE_DIR)

test:
	go test ./...
