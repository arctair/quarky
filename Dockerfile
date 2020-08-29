FROM debian
COPY bin/quarky /bin/quarky
ENTRYPOINT ["/bin/quarky"]
