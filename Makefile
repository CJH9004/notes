build: 
	gitbook build . docs && touch docs/.nojekyll && rm docs/Makefile