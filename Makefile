all: center template flow

clean: 
	rm -rf ./output

center:
	mkdir -p output
	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/merge-img_arm ./bin/center/main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/merge-img_amd ./bin/center/main.go
	makefat ./output/merge-img ./output/merge-img_*
	rm -rf ./output/merge-img_*

template:
	mkdir -p output
	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/auto-outline_arm ./bin/template/clipboard/main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/auto-outline_amd ./bin/template/clipboard/main.go
	makefat ./output/auto-outline ./output/auto-outline_*
	rm -rf ./output/auto-outline_*

	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/img-outline_arm ./bin/template/file/main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/img-outline_amd ./bin/template/file/main.go
	makefat ./output/img-outline ./output/img-outline_*
	rm -rf ./output/img-outline_*

	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/screencapture-outline_arm ./bin/template/screencapture/main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/screencapture-outline_amd ./bin/template/screencapture/main.go
	makefat ./output/screencapture-outline ./output/screencapture-outline_*
	rm -rf ./output/screencapture-outline_*

flow:
	mkdir -p output/flow
	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o output/flow/img-outline-flow_arm ./bin/flow/main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o output/flow/img-outline-flow_amd ./bin/flow/main.go
	makefat ./output/flow/img-outline-flow ./output/flow/img-outline-flow_*
	rm -rf ./output/flow/img-outline-flow_*

	cp ./bin/flow/block_snipaste.sh ./output/flow/block_snipaste.sh
	cp -r ./bin/template/imgs ./output/flow/imgs
