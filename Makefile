BUILDDIR = build
THEMENAME = default
THEMEDIR = theme/$(THEMENAME)
JSNAME = fotoDen
JSDIR = js
JSMIN = terser
TOOLPKG = github.com/vulppine/fotoDen/tool
JSSUM = md5sum $(JSDIR)/$(JSNAME).js | cut -d " " -f 1
JSMINSUM = md5sum $(BUILDDIR)/$(JSNAME).min.js | cut -d " " -f 1
.PHONY: all minjs tool theme

$(shell mkdir $(BUILDDIR))

all: minjs theme tool

minjs:
	@echo "Minifying fotoDen..."
	$(JSMIN) -c -m -o $(BUILDDIR)/$(JSNAME).min.js $(JSDIR)/$(JSNAME).js

theme:
	@echo "Packaging theme..."
	$(shell mkdir $(BUILDDIR)/theme/)
	$(shell mkdir $(BUILDDIR)/theme/html)
	$(shell mkdir $(BUILDDIR)/theme/js)
	$(foreach js,$(notdir $(wildcard $(THEMEDIR)/js/*)),$(shell $(JSMIN) $(THEMEDIR)/js/$(js) -c -m -o $(BUILDDIR)/theme/js/$(js)))
	cp -r $(THEMEDIR)/html/* $(BUILDDIR)/theme/html/
	cp $(THEMEDIR)/theme.json $(BUILDDIR)/theme/
	cd $(BUILDDIR)/theme;\
	zip ../$(THEMENAME)_theme.zip -r *
	@echo "Cleaning up..."
	rm -r build/theme/

tool:
	@echo "Making tool..."
	if [ -e $(BUILDDIR)/$(JSNAME).min.js ]; then\
		echo "Building with minified JS...";\
		go build -o $(BUILDDIR)/fotoDen\
		-ldflags "-X '$(TOOLPKG).JSSum=$(shell $(JSSUM))' -X '$(TOOLPKG).JSMinSum=$(shell $(JSMINSUM))'";\
	else\
		echo "Building without minified JS...";\
		go build -o $(BUILDDIR)/fotoDen\
		-ldflags "-X '$(TOOLPKG).JSSum=$(shell $(JSSUM))'";\
	fi;
