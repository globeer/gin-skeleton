# docker build --rm -t docker.bower.co.kr/gogin:0.1.0 . && docker rmi $(docker images -f dangling=true -q)
FROM golang:1.16.2-alpine3.13 as builder

ARG USER_ID=1000
ARG GROUP_ID=1000
ARG APP_USER=bower
# ARG registryUrl
# ARG registryCredentialsId


WORKDIR /tmp/alpine-golang-image
COPY . .

# RUN addgroup -g $GROUP_ID $APP_USER 
RUN adduser -u $USER_ID -g $GROUP_ID -D -H $APP_USER
RUN apk add --no-cache upx
RUN go mod tidy && go get -u -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o main main.go
RUN go get github.com/pwaller/goupx
RUN goupx main

# ---------------------------- 5.59 MB ----------------------------
# useradd 에 따른 디스크 사용량 증가
FROM busybox:1.33.0
# ---------------------------- 4.52 MB ----------------------------
# FROM scratch
# ---------------------------- common ----------------------------
ENV APP_USER bower
ENV APP_HOME /home/$APP_USER
ENV APP_PORT=8080
# for busybox
# RUN adduser -u 1000 -D -H $APP_USER
# RUN mkdir -p $APP_HOME/app/log $APP_HOME/app/config
# work
WORKDIR $APP_HOME
# COPY --from=builder /usr/share/zoneinfo/Asia/Seoul /etc/localtime
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --chown=bower:bower --from=builder /tmp/alpine-golang-image/main $APP_HOME
COPY --chown=bower:bower app/ app/
COPY --chown=bower:bower public/ public/
COPY --chown=bower:bower view/ view/

# for busybox
# RUN rm -f $APP_HOME/app/config/config.yml
# RUN chown -R $APP_USER:$APP_USER $APP_HOME/
USER $APP_USER
# host 와 통신 port
EXPOSE $APP_PORT
VOLUME ["$APP_HOME/app/log", "$APP_HOME/app/config"]
CMD ["./main"]
