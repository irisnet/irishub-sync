FROM alpine:3.8

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up GOPATH & PATH

ENV PROJECT_NAME irishub-sync
ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/irisnet
ENV REPO_PATH    $BASE_PATH/$PROJECT_NAME
ENV LOG_DIR      /$PROJECT_NAME/log
ENV PATH         $GOPATH/bin:$PATH

# Set volumes

VOLUME $LOG_DIR

# Link expected Go repo path

RUN mkdir -p $GOPATH/pkg $GOPATH/bin $BASE_PATH $REPO_PATH $LOG_DIR

# Add source files

COPY . $REPO_PATH

# Install minimum necessary dependencies, build irishub-server
RUN apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all && \
    mv $REPO_PATH/$PROJECT_NAME $GOPATH/bin && \
    rm -rf $REPO_PATH/vendor && \
    rm -rf $GOPATH/src/github.com/golang $GOPATH/bin/dep $GOPATH/pkg/* && \
    apk del $PACKAGES

VOLUME ["$LOG_DIR"]


CMD irishub-sync > $LOG_DIR/debug.log && tail -f $LOG_DIR/debug.log
