FROM alpine

WORKDIR /app
ENV PORT 8080
EXPOSE 8080

ADD ./main /app/main

CMD ["/app/main", "-server"]
