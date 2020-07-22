package main

import (
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func getPHPVersion(t *testing.T) string {
	image := os.Getenv("PHP_VERSION")
	assert.NotEmpty(t, image, "Environment variable `PHP_VERSION` must not be empty")
	return image
}

func getPHPImageTag(version string) string {
	return "php:" + version + "-fpm-alpine"
}

func TestPHPFPMDev(t *testing.T) {
	version := getPHPVersion(t)
	tag := "magento2-php-fpm-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"PHP_FPM_IMAGE=" + getPHPImageTag(version),
	}

	buildOptions := &docker.BuildOptions{
		Tags:         []string{tag},
		BuildArgs:    args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	// PHP version
	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -r 'echo PHP_VERSION;'"}}
	output := docker.Run(t, tag, opts)
	regex := "^" + strings.ReplaceAll(version, ".", "\\.")
	assert.Regexp(t, regex, output)

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

func TestPHPFPMDevHasXdebug(t *testing.T) {
	version := getPHPVersion(t)
	tag := "magento2-php-fpm-test"
	args := []string{
		"BUILD_ENVIRONMENT_IMAGE=magento2-php-fpm-development-onbuild",
		"BUILD_PHP_XDEBUG_ENABLE=1",
		"PHP_FPM_IMAGE=" + getPHPImageTag(version),
	}
	buildOptions := &docker.BuildOptions{
		Tags:         []string{tag},
		BuildArgs:    args,
		OtherOptions: []string{"-f", "../build/php-fpm/Dockerfile"},
	}

	docker.Build(t, "..", buildOptions)

	opts := &docker.RunOptions{Command: []string{"bash", "-c", "php -m | grep xdebug"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "xdebug", output)
}
