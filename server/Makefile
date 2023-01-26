
build: clean my-collection ui

ui:
	cd ../frontend && yarn build && cd ../server && cp -r ../frontend/build output/ui

my-collection:
	go build -o output/my-collection

clean:
	rm -rf output && mkdir output
