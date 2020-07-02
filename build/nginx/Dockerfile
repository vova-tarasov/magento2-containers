# Copy pub/ folder with generated static assets
ARG PHP_FPM_IMAGE_TAG=magento2-php-fpm:latest

FROM ${PHP_FPM_IMAGE_TAG} as php-fpm

FROM nginx:1-alpine

MAINTAINER Volodymyr Tarasov <volodymyr.tarasov@gmail.com>

COPY build/nginx/etc/default.conf /etc/nginx/conf.d/default.conf
COPY build/nginx/etc/nginx.conf /etc/nginx/nginx.conf

RUN apk add shadow \
  # Change uid to 1000 and gid to 1000
  && usermod -u 1000 nginx \
  && groupmod -g 1000 nginx \
  && chown nginx:nginx /var/log/nginx

WORKDIR /var/www/html

USER nginx

COPY --from=php-fpm /var/www/html/pub /var/www/html/pub