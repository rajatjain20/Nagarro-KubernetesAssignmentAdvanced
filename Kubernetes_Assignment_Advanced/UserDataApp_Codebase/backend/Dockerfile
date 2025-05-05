# Use an official Ubuntu image as the base image
FROM ubuntu:latest

# update the package list and install necessary packages
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    build-essential

# install Go
RUN wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz && \
    rm go1.23.3.linux-amd64.tar.gz

# set Go environment variables
# and set the TimeZone env variable, so that the logs will have IST time
ENV PATH="/usr/local/go/bin:${PATH}" \
    GOPATH="/go" \
    TZ="Asia/Kolkata"

# Create a group and user 
RUN groupadd -r backendusergroup && useradd -r -g backendusergroup backenduser

# Create a directory for the app 
RUN mkdir -p /go/src/app 

# Create and set permissions for the build cache directory 
RUN mkdir -p /home/backenduser/.cache/go-build && \ 
   chown -R backenduser:backendusergroup /home/backenduser/.cache

# Create necessary directories and set permissions
RUN mkdir -p /go/pkg && \ 
   chown -R backenduser:backendusergroup /go

# set work directory inside container
WORKDIR /go/src/app/

# copy everything from current directory from host to workimng directory inside container
# this will ignore few files specified in .dockerignore file
COPY ./src .

# Change the ownership of the application files to the new user 
RUN chown -R backenduser:backendusergroup /go/src/app 

# Switch to the new user 
USER backenduser

# build the go code inside container
# this will generate executale "backend"
RUN go build -v -o ./ .

# Expose the port the application runs on
EXPOSE 3030

# command to execute the application
CMD [ "./backend" ]

