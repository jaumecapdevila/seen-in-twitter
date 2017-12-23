FROM alpine:3.7
# Define the current working dir
ENV WORKDIR=/app
# Define the name of binary that is going to be used as the entrypoint of the image
ENV BIN=twitter-votes
RUN apk -U add ca-certificates
# Copy the configuration file
COPY ./config/config.yml ${WORKDIR}/config/config.yml
# Copy the generated binary inside the container
COPY ./bin/${BIN} ${WORKDIR}/${BIN}
# Establish the working dir defined in the env variable
WORKDIR $WORKDIR
# Establish the binary as the entrypoint
ENTRYPOINT "${WORKDIR}/${BIN}"
