ARG PHP_FPM_IMAGE=php:7.4-fpm-alpine
ARG BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-production-onbuild

FROM ${PHP_FPM_IMAGE} as magento2-php-fpm-base

ARG SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.1.9/supercronic-linux-amd64
ARG SUPERCRONIC=supercronic-linux-amd64
ARG SUPERCRONIC_SHA1SUM=5ddf8ea26b56d4a7ff6faecdd8966610d5cb9d85

ARG TOOLS="tini bash gzip lsof sed tar shadow"
ARG MODULES="\
  bcmath \
  bz2 \
  calendar \
  exif \
  gd \
  gettext \
  gmp \
  imap \
  intl \
  ldap \
  mysqli \
  opcache \
  pcntl \
  pdo_mysql \
  redis \
  shmop \
  soap \
  sockets \
  sysvmsg \
  sysvsem \
  sysvshm \
  tidy \
  xmlrpc \
  xsl \
  zip"

ENV PHP_FPM_MODE static
ENV PHP_FPM_MAX_CHILDREN 20
ENV PHP_FPM_START_SERVERS 10
ENV PHP_FPM_MIN_SPARE_SERVERS 5
ENV PHP_FPM_MAX_SPARE_SERVERS 10
ENV PHP_FPM_MAX_REQUESTS 0
ENV PHP_MEMORY_LIMIT 2G
ENV PHP_DISPLAY_ERRORS Off
ENV PHP_TIMEZONE UTC
ENV PHP_SENDMAIL_PATH_SERVER 'localhost:1025'

ENV PHP_OPCACHE_ENABLE 1
ENV PHP_OPCACHE_MEMORY_CONSUMPTION 512
ENV PHP_OPCACHE_MAX_ACCELERATED_FILES 100000
ENV PHP_OPCACHE_VALIDATE_TIMESTAMPS 0
ENV PHP_OPCACHE_REVALIDATE_FREQ 0
ENV PHP_OPCACHE_CONSISTENCY_CHECK 0

ENV NEWRELIC_ENABLED 1
ENV NEWRELIC_LICENSE ''
ENV NEWRELIC_APPNAME 'Magento'
ENV NEWRELIC_DAEMON_ADDRESS 'newrelic-agent:31339'

# Install build dependencies
RUN apk add --no-cache --virtual .build-deps ${PHPIZE_DEPS} \
    icu-dev \
    gettext-dev \
    gmp-dev \
    openldap-dev \
    bzip2-dev \
    freetype-dev \
    libpng-dev \
    libjpeg-turbo-dev \
    libxml2-dev \
    tidyhtml-dev \
    libxslt-dev \
    libzip-dev \
    yaml-dev \
    ca-certificates \
  # Install runtime dependencies
  && apk add --no-cache \
    ${TOOLS} \
    icu-libs \
    gettext \
    gmp \
    imap-dev \
    openldap \
    bzip2 \
    libpng \
    libjpeg-turbo \
    freetype \
    libxml2 \
    tidyhtml \
    libxslt \
    libzip \
    yaml \
  # Change uid to 1000 and gid to 1000
  && usermod -u 1000 www-data \
  && groupmod -g 1000 www-data \
  # Configure GD (for PHP version >=5.4 use different configure args)
  && if [ "${PHP_VERSION::3}" = "7.2" ] || [ "${PHP_VERSION::3}" = "7.3" ]; then \
        docker-php-ext-configure gd \
        --with-gd \
        --with-freetype-dir=/usr/include/ \
        --with-png-dir=/usr/include/ \
        --with-jpeg-dir=/usr/include/; \
      else \
        docker-php-ext-configure gd \
        --with-freetype \
        --with-jpeg; \
      fi \
  # Configure igbinary for better compression and Yaml
  && pecl install -o -f igbinary yaml lzf apcu \
  && docker-php-ext-enable igbinary yaml lzf apcu \
  # Configure redis
  && mkdir -p /usr/src/php/ext \
  && cd /usr/src/php/ext \
  && pecl bundle redis \
  && docker-php-ext-configure redis --enable-redis-igbinary --enable-redis-lzf \
  # Install PHP modules
  && docker-php-ext-install -j$(nproc) ${MODULES} \
  # Install a cron runner for containers
  && cd /tmp \
  && chmod +s /bin/busybox \
  && curl -fsSLO "$SUPERCRONIC_URL" \
  && echo "${SUPERCRONIC_SHA1SUM}  ${SUPERCRONIC}" | sha1sum -c - \
  && chmod +x "$SUPERCRONIC" \
  && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
  && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic \
  # Install Composer
  && curl -sSL https://getcomposer.org/installer | php -- --1 --install-dir=/usr/bin --filename=composer \
  && composer global require hirak/prestissimo --no-plugins --no-scripts \
  # Remove build dependencies to minify the image size
  && rm -rf /tmp/* \
  && apk del --purge .build-deps

# Install Newrelic
RUN curl -L https://download.newrelic.com/php_agent/archive/9.10.0.262/newrelic-php5-9.10.0.262-linux-musl.tar.gz | tar -C /tmp -zx \
  && NR_INSTALL_USE_CP_NOT_LN=1 NR_INSTALL_SILENT=1 /tmp/newrelic-php5-*/newrelic-install install \
  && rm -rf /tmp/*

COPY build/php-fpm/etc/www.conf /usr/local/etc/php-fpm.d/www.conf
COPY build/php-fpm/etc/php.ini /usr/local/etc/php/php.ini
COPY build/php-fpm/etc/conf.d/ /usr/local/etc/php/conf.d/

#### Production image ####
FROM magento2-php-fpm-base as magento2-php-fpm-production-onbuild

ONBUILD COPY src/composer.json src/composer.lock ./

# Get auth.json from another container to prevent persisting it in resulting image or layers
ONBUILD RUN COMPOSER_AUTH=$(wget -O - -q http://credentials-server:8080/auth.json) composer install --prefer-dist --no-interaction --no-progress \
    && rm -rf ~/.composer && chown -R www-data:www-data .

ONBUILD USER www-data

ONBUILD COPY --chown=www-data:www-data src/ .

ONBUILD RUN bin/magento setup:di:compile \
        && composer dump-autoload --optimize --apcu \
        && bin/magento setup:static-content:deploy -f -s compact

#### Development image ####
FROM magento2-php-fpm-base as magento2-php-fpm-development-onbuild

ONBUILD ARG BUILD_PHP_XDEBUG_ENABLE=0
ONBUILD ARG BUILD_PHP_XDEBUG_IDEKEY=PHPSTORM
ONBUILD ARG BUILD_PHP_XDEBUG_MAX_NESTING_LEVEL=2000
ONBUILD ARG BUILD_PHP_XDEBUG_REMOTE_ENABLE=1
ONBUILD ARG BUILD_PHP_XDEBUG_CLIENT_HOST=10.254.254.254

ONBUILD RUN if [ "$BUILD_PHP_XDEBUG_ENABLE" = "1" ] ; then \
    apk add --no-cache --virtual .build-deps ${PHPIZE_DEPS} \
    && pecl install -o -f xdebug \
    && echo "zend_extension=xdebug.so" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.idekey = $BUILD_PHP_XDEBUG_IDEKEY" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.max_nesting_level = $BUILD_PHP_XDEBUG_MAX_NESTING_LEVEL" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.remote_enable = $BUILD_PHP_XDEBUG_REMOTE_ENABLE" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.client_host = $BUILD_PHP_XDEBUG_CLIENT_HOST" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && echo "xdebug.mode = debug" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini \
    && rm -rf /tmp/* \
    && apk del --purge .build-deps \
  ;fi

ONBUILD USER www-data

ONBUILD RUN mkdir -p pub

FROM ${BUILD_ENVIRONMENT_IMAGE} as magento2-php-fpm