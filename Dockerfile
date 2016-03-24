FROM scratch

EXPOSE 8080
ADD ./main /main
ADD ./dockerbuild /
CMD ["/main", "-server"]
