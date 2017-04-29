FROM centurylink/ca-certs
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

EXPOSE 80

ADD entrypoint entrypoint
ADD html html
ENTRYPOINT ["./entrypoint"]
