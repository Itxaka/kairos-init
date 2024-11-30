# Dockerfile to test the init feature

FROM golang AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/kairos-init .

FROM ubuntu:24.04 as base-kairos
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=build /app/kairos-init /kairos-init
RUN /kairos-init -l debug -f all
RUN rm /kairos-init
# Move to init
RUN kernel=$(ls /lib/modules | head -n1) && depmod -a "${kernel}"
RUN echo "" > /etc/machine-id || true
RUN rm -rf /boot/initramfs-* || true
RUN rm /var/lib/dbus/machine-id || true
RUN rm /etc/hostname || true
RUN systemctl disable systemd-pcrlock-make-policy || true; \
      systemctl mask systemd-pcrlock-make-policy || true; \
      journalctl --vacuum-size=1K || true;
