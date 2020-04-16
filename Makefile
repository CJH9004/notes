.PHONY: all
all: 
	gitbook build . docs && touch docs/.nojekyll && rm docs/Makefile && git add . && git commit -m "edited" && git push origin master