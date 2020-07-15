FROM golang:1.14.5-alpine3.11

RUN echo '@testing http://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    echo '@main http://dl-cdn.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories && \
    apk add --upd build-base \
        pkgconfig \
        libgsf-dev \
        openjpeg-dev \
        graphicsmagick-dev \
        openexr-dev \
        poppler-dev \
        librsvg-dev \
        fftw-dev@main \
        sqlite-dev \
        vips-dev@testing

ARG OPENSLIDE_VERSION=3.4.1
ARG OPENSLIDE_URL=https://github.com/openslide/openslide/releases/download

COPY openslide-init.patch /usr/local/share

RUN wget ${OPENSLIDE_URL}/v${OPENSLIDE_VERSION}/openslide-${OPENSLIDE_VERSION}.tar.gz && \
	tar xf openslide-${OPENSLIDE_VERSION}.tar.gz && \
	patch -p0 </usr/local/share/openslide-init.patch && \
    cd openslide-${OPENSLIDE_VERSION} && \
	./configure && \
	make &&\
    make install
