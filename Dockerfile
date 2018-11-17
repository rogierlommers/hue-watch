FROM ubuntu
LABEL description="HUE watch"
LABEL maintainer="Rogier Lommers <rogier@lommers.org>"

# install dependencies
RUN apt-get update  
RUN apt-get install -y ca-certificates curl

# add binary
COPY --chown=1000:1000 binary/hue-watch-amd64 /hue-watch-amd64

# change to data dir and run bianry
WORKDIR "/"
CMD ["/hue-watch-amd64"]
