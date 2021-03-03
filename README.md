# Magento 2 Containers

This repository is a first building block to run Magento 2 in a cloud (GCP, AWS, Azure) and is suitable for all environments (development, UAT, production).

[![Build Status](https://travis-ci.com/vova-tarasov/magento2-containers.svg?branch=master)](https://travis-ci.com/vova-tarasov/magento2-containers)


- [Goal](#goal)
- [Motivation](#motivation)
- [Development containers](#development-containers)
    * [Initial setup](#initial-setup)
        * [Prerequisites](#prerequisites)
        * [Getting started with a new project](docs/getting-started-new.md)
        * [Getting started with an existing project](docs/getting-started-existing.md)
    * [Components](#components)
        * [PHP-FPM](#php-fpm)
        * [Database](#database)
        * [Nginx](#nginx)
        * [Varnish](#varnish)
        * [Elasticsearch](#elasticsearch)
        * [Redis](#redis)
        * [CRON](#cron)
        * [Email](#email)
    * [Debugging and profiling](#debugging-and-profiling)
        * [Xdebug](#xdebug)
        * [New Relic](#new-relic)
- [Production containers](#production-containers)
    * [Quick setup](#quick-setup)
        * [Build](#build)

## Goal
To create a reproducible and reusable containerized environment that works on a local machine and in production.      

## Motivation
Keeping production and development environments in sync and up to date is essential for productive project development.

Each developer should have an ability to easily create and spin up the same or as close as possible to production setup on a local machine.
Given that containers are a standard nowadays, 
propagating changes to code, containers and their dependencies should be 
a pretty straight-forward task.

## Development containers
### Initial setup
#### Prerequisites
Ensure the following conditions are met:

1. Locally installed [Docker](https://www.docker.com/products/docker-desktop) (Docker Engine >=18.09 version)
2. Docker has enough resources allocated:
    - at least 10 GB of free disk space for hosting containers
    - at least 2 CPU cores (4 CPU recommended)
    - at least 4 GB of RAM (6 GB recommended)

### Components
> When building or modifying a component, prioritize production environment first and override it for your development needs via Docker build args or ENV variables.

The project should compile for a production environment with little to no intervention to configuration.

Below you will an answer to most of the use cases, but if you spot a missing gap, feel free to create a PR and contribute to the document.
#### PHP-FPM
> PHP-FPM allows you to configure most of the values via ENV variables which makes it handy to tune the same Docker image for both performance and development needs.

Following [the DRY rule](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself), the same Dockerfile should be used for all environments.
However, some instructions like installation of Xdebug should not slip into production. 

Therefore, [multi-stage build](https://docs.docker.com/develop/develop-images/multistage-build/) will have [Docker ONBUILD](https://docs.docker.com/engine/reference/builder/#onbuild) instructions explicitly needed to keep the resulting image similar but with small differences for each environment.
The resulting [Dockerfile](build/php-fpm/Dockerfile) will have the following view:

```dockerfile
ARG PHP_FPM_IMAGE=php:7.2-fpm-alpine
ARG BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-production-onbuild

# Base image
FROM ${PHP_FPM_IMAGE} as magento2-php-fpm-base
RUN ... software installation for all environments

# Production image
FROM magento2-php-fpm-base as magento2-php-fpm-production-onbuild
ONBUILD ... production specific instructions 

# Development image
FROM magento2-php-fpm-base as magento2-php-fpm-development-onbuild
ONBUILD ... development specific instructions 

# Final image compiles based on $BUILD_ENVIRONMENT_IMAGE argument
# "Base image" + "Production image" instructions
# OR
# "Base image" + "Development image" instructions
FROM ${BUILD_ENVIRONMENT_IMAGE} as magento2-php-fpm
```

By default, php.ini and PHP-FPM are configured to run as a container for production use. They may still be optimal for a development purpose. 

A snippet of [php.ini](build/php-fpm/etc/php.ini) file

   ```ini
    ; Maximum amount of memory a script may consume (128MB)
    ; http://php.net/memory-limit
    memory_limit = ${PHP_MEMORY_LIMIT}
   ```

To redefine PHP_MEMORY_LIMIT value *(or any other)*, set it in [docker-compose.yaml](docker-compose.yaml) or [.env](.env) file

   ```dotenv
    MAGENTO_RUN_MODE=development
    PHP_DISPLAY_ERRORS=1
    PHP_OPCACHE_CONSISTENCY_CHECK=1
    PHP_MEMORY_LIMIT=4G

   ```

To change PHP version, modify the `PHP_FPM_IMAGE` argument of the build section. Currently, 7.2, 7.3 and 7.4 versions supported

 ```yaml
version: "3.7"
  services:
   php-fpm:
    build:
      context: .
      dockerfile: ./build/php-fpm/Dockerfile
      args:
        PHP_FPM_IMAGE: "php:7.4-fpm-alpine"
  ```

#### Database
The Magento software recommends [MariaDB](https://mariadb.org/). It can be replaced with [MySQL](https://www.mysql.com/) or [Percona](https://www.percona.com/) based on eventual production setup. 

MariaDB version of the config file by default has a low amount of resource set.
In order to speed things up for local development [my.cnf](build/mysql/etc/my.cnf) has `innodb_buffer_pool_size` set to `1G`.
You may need to adjust it based on your database size and machine capacity.

`slow_query_log` is enabled to log queries that are slower than 2 seconds.
To watch the file for changes, you may use the following command:


  ```shell script
    docker exec -it $(docker ps -f name=mysql -q) tail -f /var/log/mysql/mariadb-slow.log
  ```

#### Nginx
> Due to security reason, containers running in production must not use systems ports (0-1023). 

Nginx runs under the same UID = 1000, and GID = 1000 as PHP-FPM and listens on port `8080`.

1. To run local environment without FPC, open the website directly on http://magento2.local:8080 or change the Nginx port mapping to `80:8080` in [docker-compose.yaml](docker-compose.yaml) as shown below 

    ```yaml
      version: "3.7"
      services:
        varnish:
          ports:
          - "6081"
        nginx:
          ...
          ports:
          - "80:8080"
     ``` 

2. Changing anything in Nginx [configuration](build/nginx/etc/default.conf) requires a restart
   ```shell script
      docker exec -it $(docker ps -f name=nginx -q) nginx -s reload
   ```

#### Varnish
> Due to security reason, containers running in production must not use systems ports (0-1023). 

Varnish listens on its default port `6081`.

To run the local environment with FPC on, ensure you have the following configuration in [docker-compose.yaml](docker-compose.yaml) as shown below

  ```yaml
    version: "3.7"
    services:
      varnish:
        ports:
          - "80:6081"
    nginx:
        ...
        ports:
          - "8080"
  ```

#### Elasticsearch
Elasticsearch is a default search engine for Magento starting with version 2.3. MySQL search engine support is removed in Magento 2.4     

For simplicity' sake and to keep resource consumption low, Elasticsearch runs in a single-node configuration.

#### Redis
The same instance of Redis may be used to store user sessions and cache.

Connecting to Redis is pretty straight-forward, from your command line run 

  ```shell script
    docker exec -it $(docker ps -f name=redis -q) redis-cli
  ```

To clear all the data, use the following command

  ```shell script
    docker exec -it $(docker ps -f name=redis -q) redis-cli FLUSHALL
  ```

#### Cron 
> Due to security reason, containers must not run under `root`

Due to the nature of `cron`, it runs only under `root` user and not a good fit for the containerized solution. [Supercronic](https://github.com/aptible/supercronic) was designed specifically for containers to replace standard cron. Many features come out of the box, including graceful shutdown and logging to STDOUT.

#### Emails
Emails configured via `sendmail` command in [php.ini](build/php-fpm/etc/php.ini) file

  ```ini
    ; For Unix only.  You may supply arguments as well (default: "sendmail -t -i").
    ; http://php.net/sendmail-path
    sendmail_path = sendmail -t -i -S ${PHP_SENDMAIL_PATH}
  ``` 

External email providers:
 - AWS SES
 - Sendinblue
 - Sendgrid
 - MailHog (development setup)

  ```ini
    PHP_SENDMAIL_PATH=mail:1025
  ```
     
### Debugging and profiling
#### Xdebug
It is known Xdebug can kill the performance and make your setup work slow. It also drastically degrades the speed of `composer install` command, so it's recommended to turn it off during initial project import and for Frontend work.

To turn Xdebug on in [docker-compose.yaml](docker-compose.yaml) change `BUILD_PHP_XDEBUG_ENABLE` to `1` and ensure you have `BUILD_PHP_XDEBUG_REMOTE_HOST` defined  


  ```ini
    services: 
      ... 
      php-fpm:
        build:
          ...
          args:
            BUILD_PHP_XDEBUG_ENABLE: 1
            BUILD_PHP_XDEBUG_REMOTE_HOST: "host.docker.internal."
  ```

then rebuild the container

  ```shell script
    docker-compose up --build
  ```

#### New Relic
To enable it locally, set corresponding values in [.env](.env)


  ```dotenv
    NEWRELIC_ENABLED=1
    NEWRELIC_LICENSE=your license goes here
    NEWRELIC_APPNAME="your awesome project name"
  ```

uncomment New Relic agent in [docker-compose.yaml](docker-compose.yaml)

  ```yaml
    services:
      ... 
      newrelic-agent:
        image: newrelic/php-daemon
        restart: always
        networks:
          - backend
  ```
start PHP-FPM, CRON and New Relic containers

## Production containers
### Quick setup
##### Build
1. [Follow these steps](#starting-a-project-from-scratch) to setup the project

2. Ensure you have [app/etc/config.php](src/app/etc/config.php) file with `modules` and `scopes` configurations


  ```php
    <?php
    return [
        'modules' => [
            'Magento_AdminAnalytics' => 1,
            ...
        ],
        'scopes' => [
            'websites' => [
                'admin' => [
                    'website_id' => '0',
                    'code' => 'admin',
                    'name' => 'Admin',
                    'sort_order' => '0',
                    'default_group_id' => '0',
                    'is_default' => '0',
                ],
                'base' => [
                    'website_id' => '1',
                    'code' => 'base',
                    'name' => 'Main Website',
                    'sort_order' => '0',
                    'default_group_id' => '1',
                    'is_default' => '1',
                ],
            ],
            'groups' => [
                0 => [
                    'group_id' => '0',
                    'website_id' => '0',
                    'name' => 'Default',
                    'root_category_id' => '0',
                    'default_store_id' => '0',
                    'code' => 'default',
                ],
                1 => [
                    'group_id' => '1',
                    'website_id' => '1',
                    'name' => 'Main Website Store',
                    'root_category_id' => '2',
                    'default_store_id' => '1',
                    'code' => 'main_website_store',
                ],
            ],
            'stores' => [
                'admin' => [
                    'store_id' => '0',
                    'code' => 'admin',
                    'website_id' => '0',
                    'group_id' => '0',
                    'name' => 'Admin',
                    'sort_order' => '0',
                    'is_active' => '1',
                ],
                'default' => [
                    'store_id' => '1',
                    'code' => 'default',
                    'website_id' => '1',
                    'group_id' => '1',
                    'name' => 'Default Store View',
                    'sort_order' => '0',
                    'is_active' => '1',
                ],
            ],
        ]
    ];

  ```

3. Prepare [access keys to Magento 2](https://marketplace.magento.com/customer/accessKeys/) and optionally to your GitHub account by
    copying [auth.json.sample](auth.json.sample) to `auth.json` and replacing credentials.
    > `auth.json` won't be added to the final image nor leave a trace in docker build history 

4. Run `make` to build images

   ```shell script
    make
   ```

5. Now you're ready to publish them into a cloud of your choice