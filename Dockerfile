FROM scratch
EXPOSE 9200
ENTRYPOINT ["/quicksearch", "-c", "config.yaml"]
COPY quicksearch /
COPY configs/config.yaml .

