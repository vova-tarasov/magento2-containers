package main

import (
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

