MOTMOT=Motmot -load-cache
INPUT=Notes.txt.mot
PROCESSING=Process.mot

MARKDOWN=../README.md

.PHONY: run markdown
run: markdown
markdown: $(MARKDOWN)

$(MARKDOWN): $(INPUT) $(PROCESSING)
	$(MOTMOT) -file $(PROCESSING) -exp 'Notes.Processing.run ()' -prefix Markdown
