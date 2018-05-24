FROM alpine:edge

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up GOPATH & PATH

ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/irisnet
ENV REPO_PATH    $BASE_PATH/iris-sync-server
ENV LOG_DIR      /sync-iris/log
ENV PATH         $GOPATH/bin:$PATH

# Set volumes

VOLUME $LOG_DIR:sync-iris-log

# Link expected Go repo path

RUN mkdir -p $LOG_DIR $GOPATH/pkg $GOPATH/bin $BASE_PATH $REPO_PATH

# Add source files

COPY . $REPO_PATH

# Install minimum necessary dependencies, build iris-sync-server
RUN apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all && \
    mv $REPO_PATH/sync-iris $GOPATH/bin && \
    rm -rf $REPO_PATH/vendor && \
    apk del $PACKAGES

CMD sync-iris > $LOG_DIR/debug.log && tail -f $LOG_DIR/debug.log