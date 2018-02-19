#build our docker image with name newstreeter/midiparser
#docker build -t newstreeter/midiparser .

#run our docker container afterwards remove himself
#docker run --rm -it newstreeter/midiparser:latest

#Docker Remove All <none> images (only run in bash terminal)
#docker rmi $(docker images -f "dangling=true" -q)

################################################################



#name of base image
FROM golang:alpine

#create a folder where our program will be located
RUN mkdir -p /go/src/bitbucket.org/NewStreeter/MIDIParser

#set a working directory with a created folder
WORKDIR /go/src/bitbucket.org/NewStreeter/MIDIParser

#Copy all files from source to the Docker's path in the image's filesystem
COPY . /go/src/bitbucket.org/NewStreeter/MIDIParser

#run all tests include subpackages with coverage and list root files
CMD go test -v -cover ./... && ls -la