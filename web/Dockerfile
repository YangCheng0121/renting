FROM alpine:3.2
ADD conf /conf
ADD html /html
ADD renting-web /renting-web
WORKDIR /
ENTRYPOINT [ "/renting-web" ]
EXPOSE 8888
