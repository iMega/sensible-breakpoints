FROM golang:1.10-alpine3.8

RUN echo '@testing http://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    echo '@main http://dl-cdn.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories && \
    echo '@community http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories && \
    apk add --upd pkgconfig \
        libgsf-dev \
        graphicsmagick-dev@community \
        openexr-dev@community \
        poppler-dev \
        librsvg-dev \
        fftw-dev@main \
        vips-dev@testing

ARG OPENSLIDE_VERSION=3.4.1
ARG OPENSLIDE_URL=https://github.com/openslide/openslide/releases/download

COPY openslide-init.patch /usr/local/src

RUN more /usr/local/src/openslide-init.patch

RUN wget ${OPENSLIDE_URL}/v${OPENSLIDE_VERSION}/openslide-${OPENSLIDE_VERSION}.tar.gz && \
	tar xf openslide-${OPENSLIDE_VERSION}.tar.gz && \
	patch -p0 </usr/local/src/openslide-init.patch && \
    cd openslide-${OPENSLIDE_VERSION} && \
	./configure && \
	make &&\
    make install
