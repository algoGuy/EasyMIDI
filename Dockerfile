################################################################

#name of base image
FROM golang:alpine

#create a folder where our program will be located
RUN mkdir -p /go/src/github.com/algoGuy/EasyMIDI

#set a working directory with a created folder
WORKDIR /go/src/github.com/algoGuy/EasyMIDI

#Copy all files from source to the Docker's path in the image's filesystem
COPY . /go/src/github.com/algoGuy/EasyMIDI

#run all tests include subpackages with coverage and list root files
CMD go test -v -cover ./... && ls -la