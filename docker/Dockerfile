FROM golang:1.18.2 AS gobuilder

RUN apt-get update && apt-get install -y libpcap-dev

WORKDIR /minimega/

COPY . .

RUN ./scripts/all.bash


FROM ubuntu:22.04

RUN apt update && apt install -y \
  dnsmasq \
  g++ \
  git \
  iproute2 \
  iptables \
  isc-dhcp-client \
  libnl-3-dev \
  libnl-genl-3-dev \
  libpcap-dev \
  make \
  ntfs-3g \
  openssh-client \
  openvswitch-switch \
  qemu-kvm \
  qemu-utils \
  wget \
  && apt autoremove -y \
  && apt clean -y\
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /var/cache/apt/archives/*

RUN git clone https://github.com/Raizo62/vwifi.git -b v6.3 /tmp/vwifi && \
  cd /tmp/vwifi && \
  make gitversion && \
  make update && \
  make && \
  make tools && \
  make install && \
  rm -rf /tmp/vwifi

COPY --from=gobuilder /minimega/bin/minimega /opt/minimega/bin/minimega
COPY --from=gobuilder /minimega/bin/miniweb  /opt/minimega/bin/miniweb

# The version of miniccc used has to match exactly with the version of minimega
# running, so let's include them in the image so we can easily grab a copy that
# is sure to be the same version.
COPY --from=gobuilder /minimega/bin/miniccc     /opt/minimega/bin/miniccc
COPY --from=gobuilder /minimega/bin/miniccc.exe /opt/minimega/bin/miniccc.exe

# For the sake of consistency, let's go ahead and include protonuke and
# minirouter in the image too so we can easily grab a copy if/when necessary.
COPY --from=gobuilder /minimega/bin/protonuke     /opt/minimega/bin/protonuke
COPY --from=gobuilder /minimega/bin/protonuke.exe /opt/minimega/bin/protonuke.exe
COPY --from=gobuilder /minimega/bin/minirouter    /opt/minimega/bin/minirouter

# As the minimega API changes, so does the minimega.py generated file. Given
# this, let's go ahead and also include the lib directory so we can grab the
# updated Python package from the Docker image.
COPY --from=gobuilder /minimega/lib     /opt/minimega/lib
COPY --from=gobuilder /minimega/README  /opt/minimega/lib/README
COPY --from=gobuilder /minimega/VERSION /opt/minimega/lib/VERSION

COPY ./web       /opt/minimega/web
COPY ./docker/mm /usr/local/bin/mm

COPY ./docker/start-minimega.sh /start-minimega.sh

RUN chmod +x /usr/local/bin/mm \
  && chmod +x /start-minimega.sh

WORKDIR /opt/minimega

CMD ["/start-minimega.sh"]
