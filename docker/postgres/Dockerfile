FROM postgres:15
LABEL authors="ebkopec"

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y postgresql-15-pgagent \
 && rm -rf /var/lib/apt/lists/*