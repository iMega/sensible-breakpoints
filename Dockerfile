FROM golang:1.14.5-alpine3.11

RUN apk add --upd build-base \
        pkgconfig \
        libgsf-dev \
        openjpeg-dev \
        graphicsmagick-dev \
        openexr-dev \
        poppler-dev \
        librsvg-dev \
        fftw-dev \
        sqlite-dev \
        vips-dev

ARG OPENSLIDE_VERSION=3.4.1
ARG OPENSLIDE_URL=https://github.com/openslide/openslide/releases/download

RUN wget ${OPENSLIDE_URL}/v${OPENSLIDE_VERSION}/openslide-${OPENSLIDE_VERSION}.tar.gz && \
	tar xf openslide-${OPENSLIDE_VERSION}.tar.gz && \
    cd openslide-${OPENSLIDE_VERSION} && \
	./configure && \
	make &&\
    make install
