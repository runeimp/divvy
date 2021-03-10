PROJECT_NAME := "Divvy up the Loot"
DISTRO_NAME := "DivvyUpTheLoot"
BINARY_NAME := "divvy"
SOURCE_PATH := "cmd/divvy/main.go"

alias _build-mac := _build-macos
alias _build-win := _build-windows
alias ver := version


@_default:
	just _term-wipe
	just --list



# Build compiled app
build target='':
	#!/usr/bin/env bash
	set -euo pipefail
	just _term-wipe
	if [[ '{{target}}' = '' ]]; then
		target='{{os()}}'
	else
		target='{{target}}'
	fi
	
	if [[ "${target}" = 'macos' ]] || [[ "${target}" = 'mac' ]]; then
		just _build-macos
		just distro macos-x86_64 bin/macos_amd64/{{BINARY_NAME}}
		# just distro macos_apple_silicon bin/apple_silicon/{{BINARY_NAME}}
	elif [[ "${target}" = 'windows' ]] || [[ "${target}" = 'win' ]]; then
		just _build-windows
		just distro windows_x64 bin/windows/{{BINARY_NAME}}.exe
	else
		just "_build-${target}"
		just distro "${target}" "bin/${target}/{{BINARY_NAME}}"
	fi

@_build-linux:
	echo "Building Linux app"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/linux/{{BINARY_NAME}} {{SOURCE_PATH}}

@_build-macos:
	echo "Building macOS app"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/macos_amd64/{{BINARY_NAME}} {{SOURCE_PATH}}
	# GOOS=darwin GOARCH=arm64 go build -o bin/apple_silicon/{{BINARY_NAME}} main.go # Not until Go 1.16

@_build-windows:
	echo "Building Windows app"
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/windows/{{BINARY_NAME}}.exe {{SOURCE_PATH}}
	just distro windows_x64 bin/windows/{{BINARY_NAME}}.exe


# Setup distrobution archive
distro arch file:
	#!/bin/sh
	files=( CHANGELOG.md LICENSE README.md )
	path="$(dirname "{{file}}")"
	name="$(basename "{{file}}")"
	ver="$(just version)"
	# echo "path = ${path}"
	# echo "name = ${name}"
	# echo " ver = ${ver}"
	cp ${files[@]} "${path}/"
	mkdir -p "distro/{{DISTRO_NAME}}-v${ver}"
	echo "cd ${path}"
	cd "${path}"
	echo "zip ../../distro/{{DISTRO_NAME}}-v${ver}/{{BINARY_NAME}}-v${ver}-{{arch}}.zip '${name}' ${files[@]}"
	zip "../../distro/{{DISTRO_NAME}}-v${ver}/{{BINARY_NAME}}-v${ver}-{{arch}}.zip" "${name}" ${files[@]}
	lsd -hl "../../distro/{{DISTRO_NAME}}-v${ver}"


# Run app
run +args='':
	just _term-wipe
	go run {{SOURCE_PATH}} {{args}} Giveaway-Points.csv
	@echo
	@echo "CSV Files:"
	@lsd -al *.csv
	@echo
	csvtk pretty $(ls -1 DivvyUpTheLoot_*.csv | tail -1)



# Wipes the terminal buffer for a clean start
_term-wipe:
	#!/usr/bin/env bash
	set -exo pipefail
	if [[ ${#VISUAL_STUDIO_CODE} -gt 0 ]]; then
		clear
	elif [[ ${KITTY_WINDOW_ID} -gt 0 ]] || [[ ${#TMUX} -gt 0 ]] || [[ "${TERM_PROGRAM}" = 'vscode' ]]; then
		printf '\033c'
	elif [[ "$(uname)" == 'Darwin' ]] || [[ "${TERM_PROGRAM}" = 'Apple_Terminal' ]] || [[ "${TERM_PROGRAM}" = 'iTerm.app' ]]; then
		osascript -e 'tell application "System Events" to keystroke "k" using command down'
	elif [[ -x "$(which tput)" ]]; then
		tput reset
	elif [[ -x "$(which reset)" ]]; then
		reset
	else
		clear
	fi


# Output the app version
@version:
	grep '^\tappVersion' {{SOURCE_PATH}} | cut -d'"' -f2

