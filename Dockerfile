FROM alpine:3.16

USER root

WORKDIR /root

COPY ["demo_device_plugin", "/usr/bin/"]
CMD  [ "demo_device_plugin" ]
