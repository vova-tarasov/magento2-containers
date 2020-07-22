RELEASE=1
PHP_FPM_IMAGE_TAG=magento2-php-fpm:${RELEASE}
NGINX_IMAGE_TAG=magento2-nginx:${RELEASE}
VARNISH_IMAGE_TAG=magento2-varnish:${RELEASE}

.PHONY all: build-php-fpm build-nginx build-varnish

.PHONY build-php-fpm:
	@echo ........ serving auth.json secret file from a separate container
	docker network inspect secret >/dev/null 2>&1 || docker network create --driver bridge secret
	docker stop credentials-server >/dev/null 2>&1 || true
	docker run --network secret --name credentials-server -d --rm --mount type=bind,source=${CURDIR}/auth.json,target=/creds/auth.json busybox httpd -f -p 8080 -h /creds
	@echo ........ building PHP-FPM container
	docker build --network secret -f ./build/php-fpm/Dockerfile -t ${PHP_FPM_IMAGE_TAG} .
	@echo ........ stopping to serve auth.json file
	docker stop credentials-server >/dev/null 2>&1 || true
	docker network rm secret

.PHONY build-nginx:
	@echo ........ building Nginx container
	docker build -f ./build/nginx/Dockerfile -t ${NGINX_IMAGE_TAG} --cache-from ${PHP_FPM_IMAGE_TAG} \
	    --build-arg PHP_FPM_IMAGE_TAG=${PHP_FPM_IMAGE_TAG} \
	    .

.PHONY build-varnish:
	@echo ........ building Varnish container
	docker build -f ./build/varnish/Dockerfile -t ${VARNISH_IMAGE_TAG} ./build/varnish/
