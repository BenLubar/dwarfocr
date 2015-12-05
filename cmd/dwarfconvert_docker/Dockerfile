FROM golang

RUN go get github.com/BenLubar/dwarfocr/cmd/dwarfconvert_docker
WORKDIR /go/src/github.com/BenLubar/dwarfocr/cmd/dwarfconvert_docker

EXPOSE 80

CMD ["dwarfconvert_docker"]
