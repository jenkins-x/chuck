
# FROM golang:1.15 AS builder

# RUN groupadd -g 999 user && \
#     useradd -r -u 999 -g user user

# FROM scratch
# COPY --from=builder /etc/passwd /etc/passwd
# USER user
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# COPY build/linux/chuck /

# ENTRYPOINT ["/chuck"]

FROM centos:7
COPY build/linux/chuck /
ENTRYPOINT ["/chuck"]