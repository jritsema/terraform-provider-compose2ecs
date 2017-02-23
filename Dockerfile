FROM golang:1.7

RUN apt-get update
RUN apt-get install unzip

RUN mkdir -p /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

RUN go get github.com/tools/godep
RUN go get github.com/mitchellh/gox
RUN go get github.com/inconshreveable/mousetrap

ENV GHR_VERSION v0.5.2
RUN wget https://github.com/tcnksm/ghr/releases/download/${GHR_VERSION}/ghr_${GHR_VERSION}_linux_amd64.zip && unzip ghr_${GHR_VERSION}_linux_amd64.zip -d /usr/local/bin

# build
RUN go vet -x $(go list ./... | grep -v /vendor/)
RUN godep go test -v $(go list ./... | grep -v /vendor/)
RUN gox -osarch="darwin/amd64" -osarch="linux/amd64" -osarch="windows/amd64" -output "dist/ncd_{{.OS}}_{{.Arch}}"

# push binaries to github
RUN ghr -t $GITHUB_TOKEN -u $CIRCLE_PROJECT_USERNAME -r $CIRCLE_PROJECT_REPONAME --replace `git describe --tags` dist/'