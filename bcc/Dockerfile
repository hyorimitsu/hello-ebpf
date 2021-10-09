FROM ubuntu:18.04

RUN set -ex; \
  echo "deb [trusted=yes] http://repo.iovisor.org/apt/bionic bionic-nightly main" > /etc/apt/sources.list.d/iovisor.list; \
  apt-get update -y; \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
    auditd \
    bcc-tools \
    libelf1 \
    libbcc-examples;

COPY customs /usr/share/bcc/tools
RUN chmod +x /usr/share/bcc/tools/test.py

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/bin/bash"]