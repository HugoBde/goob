TEMPLATES = $(wildcard pkg/*.templ)
COMPILED_TEMPLATES = $(patsubst pkg/%.templ,pkg/%_templ.go, $(TEMPLATES))

all: goob public/index.css

%_templ.go: %.templ
	templ generate -f $^

public/index.css: public/input.css pkg/*.templ 
	tailwind -i $< -o $@

shitter: 
	echo $(TEMPLATES)
	echo $(COMPILED_TEMPLATES)

goob: cmd/goob.go pkg/*.go $(COMPILED_TEMPLATES)
	go build $<

