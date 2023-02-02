FROM centos:7
COPY _build/server /liqiang/server
COPY ./client.html /liqiang/client.html
RUN chmod +x /liqiang/server
WORKDIR /liqiang

# PORT
EXPOSE 6000

CMD ["/liqiang/server", "-listen", "0.0.0.0:6000"]