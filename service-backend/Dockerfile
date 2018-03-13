FROM alpine
MAINTAINER 1ambda

ARG APP

COPY ./bin/${APP} /opt/service/bin/
RUN mv /opt/service/bin/* /opt/service/bin/app && \
    chmod +x /opt/service/bin/app

# gRPC Port
EXPOSE 50001

ENTRYPOINT [ "/opt/service/bin/app" ]