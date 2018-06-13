FROM alpine:edge

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up GOPATH & PATH

ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/irisnet
ENV REPO_PATH    $BASE_PATH/irishub-sync
ENV LOG_DIR      /sync-iris/log
ENV PATH         $GOPATH/bin:$PATH

# Link expected Go repo path

RUN mkdir -p $LOG_DIR $GOPATH/pkg $GOPATH/bin $BASE_PATH $REPO_PATH

# Add source files

COPY . $REPO_PATH

# Install minimum necessary dependencies, build irishub-sync
RUN apk add --no-cache $PACKAGES && \
    mv $REPO_PATH/conf/db/types.go.example $REPO_PATH/conf/db/types.go && \
    mv $REPO_PATH/conf/server/types.go.example $REPO_PATH/conf/server/types.go && \
    cd $REPO_PATH && \
    make all && \
    mv $REPO_PATH/sync-irishub $GOPATH/bin && \
    rm -rf $REPO_PATH/vendor && \
    rm -rf $GOPATH/src/github.com/golang $GOPATH/bin/dep $GOPATH/pkg/* && \
    apk del $PACKAGES

# Set volumes
VOLUME ["$LOG_DIR"]

CMD sync-irishub > $LOG_DIR/debug.log && tail -f $LOG_DIR/debug.log