package main

import (
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
	"testing"
)

const PHP_MODULES = `[PHP Modules]
apcu
bcmath
bz2
calendar
Core
ctype
curl
date
dom
exif
fileinfo
filter
ftp
gd
gettext
gmp
hash
iconv
igbinary
imap
intl
json
ldap
libxml
lzf
mbstring
mysqli
mysqlnd
newrelic
openssl
pcntl
pcre
PDO
pdo_mysql
pdo_sqlite
Phar
posix
readline
redis
Reflection
session
shmop
SimpleXML
soap
sockets
sodium
SPL
sqlite3
standard
sysvmsg
sysvsem
sysvshm
tidy
tokenizer
xml
xmlreader
xmlrpc
xmlwriter
xsl
yaml
Zend OPcache
zip
zlib

[Zend Modules]
Zend OPcache
`

func TestDev72HasAllPHPModules(t *testing.T) {
	t.Parallel()
	tag := "magento2-php-fpm-72-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -r 'echo PHP_VERSION;'"}}
	output := docker.Run(t, tag, opts)
	assert.Regexp(t, "^7\\.2\\.", output)

	opts = &docker.RunOptions{Command: []string{"php", "-m"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, PHP_MODULES, output)
}

func TestDev72HasCronSupport(t *testing.T) {
	tag := "magento2-php-fpm-72-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"which", "supercronic"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "/usr/local/bin/supercronic", output)

	opts = &docker.RunOptions{Command: []string{"which", "tini"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, "/sbin/tini", output)
}

func TestDev72HasXdebug(t *testing.T) {
	tag := "magento2-php-fpm-72-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"BUILD_PHP_XDEBUG_ENABLE=1",
	}
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c",  "php -m | grep xdebug"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "xdebug", output)
}

func TestDev73HasAllPHPModules(t *testing.T) {
	t.Parallel()
	tag := "magento2-php-fpm-73-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.3-fpm-alpine",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -r 'echo PHP_VERSION;'"}}
	output := docker.Run(t, tag, opts)
	assert.Regexp(t, "^7\\.3\\.", output)

	opts = &docker.RunOptions{Command: []string{"php", "-m"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, PHP_MODULES, output)
}

func TestDev73HasCronSupport(t *testing.T) {
	tag := "magento2-php-fpm-73-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.3-fpm-alpine",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"which", "supercronic"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "/usr/local/bin/supercronic", output)

	opts = &docker.RunOptions{Command: []string{"which", "tini"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, "/sbin/tini", output)
}

func TestDev73HasXdebug(t *testing.T) {
	tag := "magento2-php-fpm-73-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.3-fpm-alpine",
		"BUILD_PHP_XDEBUG_ENABLE=1",
	}
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c",  "php -m | grep xdebug"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "xdebug", output)
}

func TestDev74HasAllPHPModules(t *testing.T) {
	t.Parallel()
	tag := "magento2-php-fpm-74-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.4-fpm-alpine",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -r 'echo PHP_VERSION;'"}}
	output := docker.Run(t, tag, opts)
	assert.Regexp(t, "^7\\.4\\.", output)

	opts = &docker.RunOptions{Command: []string{"php", "-m"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, PHP_MODULES, output)
}

func TestDev74HasCronSupport(t *testing.T) {
	tag := "magento2-php-fpm-74-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.4-fpm-alpine",
	}

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"which", "supercronic"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "/usr/local/bin/supercronic", output)

	opts = &docker.RunOptions{Command: []string{"which", "tini"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, "/sbin/tini", output)
}

func TestDev74HasXdebug(t *testing.T) {
	tag := "magento2-php-fpm-74-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=php:7.4-fpm-alpine",
		"BUILD_PHP_XDEBUG_ENABLE=1",
	}
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
		BuildArgs: args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c",  "php -m | grep xdebug"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "xdebug", output)
}
