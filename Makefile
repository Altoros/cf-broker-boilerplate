gopackage = github.com/Altoros/cf-broker-boilerplate
gopackage_path = /go/src/$(gopackage)



build-linux:
		docker build -t cf-broker-boilerplate/build -f docker/build.Dockerfile .
		docker run -v $(PWD):$(gopackage_path) -t cf-broker-boilerplate/build go build -o $(gopackage_path)/out/broker $(gopackage)

build-osx:
		docker build -t cf-broker-boilerplate/build -f docker/build.Dockerfile .
		docker run -v $(PWD):$(gopackage_path) -e "GOOS=darwin" -t cf-broker-boilerplate/build go build -o $(gopackage_path)/out/broker $(gopackage)
