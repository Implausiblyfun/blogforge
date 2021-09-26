FROM mysql:5.7

# bootstrap from code rather than migrations for now... 
# because its ez and could mirror the approach that we would take in an integration test
COPY ./*.sql /docker-entrypoint-initdb.d