FROM alpine:3.10 AS builder

RUN apk add --no-cache \
        varnish \
        varnish-dev \
        dpkg-dev \
        make \
        automake \
        autoconf \
        libexecinfo-dev \
        libtool \
        pcre-dev \
        python3 \
        py-docutils \
    && wget "https://github.com/nigoroll/libvmod-dynamic/archive/6.2.zip" -O /tmp/libvmod-dynamic-6.2.zip \
    && unzip -d /tmp /tmp/libvmod-dynamic-6.2.zip \
    && cd /tmp/libvmod-dynamic-6.2 \
    && ./autogen.sh \
    && ./configure \
    && make \
    && make install

FROM alpine:3.10

RUN apk add --no-cache varnish shadow \
    && usermod -u 1000 varnish \
    && groupmod -g 1000 varnish \
    && chown -R varnish:varnish /etc/varnish/

COPY --from=builder /usr/lib/varnish/vmods/libvmod_dynamic.* /usr/lib/varnish/vmods/

WORKDIR /etc/varnish

COPY etc/varnish.vcl /etc/varnish/

USER varnish

CMD [ "varnishd", \
    "-F", \
    "-a", ":6081,HTTP", \
    "-s", "default,512m", \
    "-f", "/etc/varnish/varnish.vcl", \
    "-p", "feature=+esi_disable_xml_check", \
    "-p", "feature=+esi_ignore_other_elements" ]