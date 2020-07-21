package main

import (
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDev74(t *testing.T) {
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

  // PHP version
	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -r 'echo PHP_VERSION;'"}}
	output := docker.Run(t, tag, opts)
	assert.Regexp(t, "^7\\.4\\.", output)

  // Contains required PHP modules
	opts = &docker.RunOptions{Command: []string{"php", "-m"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, PHP_MODULES, output)

  // Contains supercronic
	opts = &docker.RunOptions{Command: []string{"which", "supercronic"}}
	output = docker.Run(t, tag, opts)
	assert.Equal(t, "/usr/local/bin/supercronic", output)

  // Contains tini
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
