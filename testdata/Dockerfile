FROM centos:latest
COPY ./tmp/logaxe /
RUN dd if=/dev/random of=/sample.log bs=1G count=1
CMD ["/logaxe", "/"]
