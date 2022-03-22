FROM alpine:3.10
LABEL MAINTAINER="tttlkkkl <tttlkkkl@aliyun.com>"
ENV TZ "Asia/Shanghai"
ENV TERM xterm
RUN echo 'https://mirrors.aliyun.com/alpine/v3.10/main/' > /etc/apk/repositories && \
    echo 'https://mirrors.aliyun.com/alpine/v3.10/community/' >> /etc/apk/repositories
COPY app /usr/local/bin/app
WORKDIR /app
RUN chmod +x /usr/local/bin/app \
    && apk update && apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \ 
    && echo "Asia/Shanghai" > /etc/timezone 
EXPOSE 9051 9052
CMD app