FROM openjdk:19-jdk-alpine3.16

ARG LAVALINK_VERSION=4.0.8

RUN apk add --no-cache wget libgcc udev

RUN mkdir /opt/lavalink \
    && wget https://github.com/lavalink-devs/Lavalink/releases/download/${LAVALINK_VERSION}/Lavalink.jar -qO /opt/lavalink/Lavalink.jar

COPY config/application.yml /opt/lavalink/application.yml

WORKDIR /opt/lavalink

ENTRYPOINT [ "java", "-jar", "Lavalink.jar"]