#!/bin/bash -eux

function build() {
    local version=$1
    local platform=$2
    local arch=$3
    local plugin_name=$4

    echo "building for $platform $arch"
    GOOS=$platform GOARCH=$arch go build \
        -ldflags "-X main.Version=${version}" \
        -o "${plugin_name}"
}

function buildstatic() {
    local version=$1
    local platform=$2
    local arch=$3
    local plugin_name=$4

    echo "building static for $platform $arch"
    CGO_ENABLED=0 GOOS=$platform GOARCH=$arch go build -a -tags netgo \
        -ldflags "-w -extldflags \"-static\" -X main.Version=${version}" \
        -o "${plugin_name}"
}

function movePluginsToBuildFolder() {
    local folder=$1
    mv $PLUGIN_NAME_NON_STATIC_WIN_32 $folder
    mv $PLUGIN_NAME_NON_STATIC_WIN_64 $folder
    mv $PLUGIN_NAME_NON_STATIC_LINUX_32 $folder
    mv $PLUGIN_NAME_NON_STATIC_LINUX_64 $folder
    mv $PLUGIN_NAME_NON_STATIC_LINUX_ARM64 $folder
    mv $PLUGIN_NAME_NON_STATIC_OSX $folder
    mv $PLUGIN_NAME_NON_STATIC_APPLE_ARM64 $folder
    mv $PLUGIN_NAME_STATIC_WIN_32 $folder
    mv $PLUGIN_NAME_STATIC_WIN_64 $folder
    mv $PLUGIN_NAME_STATIC_LINUX_32 $folder
    mv $PLUGIN_NAME_STATIC_LINUX_64 $folder
    mv $PLUGIN_NAME_STATIC_LINUX_ARM64 $folder
    mv $PLUGIN_NAME_STATIC_OSX $folder
    mv $PLUGIN_NAME_STATIC_APPLE_ARM64 $folder
}

function createBuildMetadataFiles() {
    local version=$1
    local folder=$2
	echo -n "v${version}" > ${folder}/version
	grep -Pzoa "(?s)## v${version}(.*?)##" CHANGELOG.md | grep -va "##" | tr -s '\n' '\n' > ${folder}/changelog
}

script_dir="$(dirname -- "$(realpath -- "${BASH_SOURCE[0]}")")"
cd "${script_dir}"

BUILD_FOLDER=build
PLUGIN_NAME_NON_STATIC_WIN_32=multiapps-plugin-non-static.win32.exe
PLUGIN_NAME_NON_STATIC_WIN_64=multiapps-plugin-non-static.win64.exe
PLUGIN_NAME_NON_STATIC_LINUX_32=multiapps-plugin-non-static.linux32
PLUGIN_NAME_NON_STATIC_LINUX_64=multiapps-plugin-non-static.linux64
PLUGIN_NAME_NON_STATIC_LINUX_ARM64=multiapps-plugin-non-static.linuxarm64
PLUGIN_NAME_NON_STATIC_OSX=multiapps-plugin-non-static.osx
PLUGIN_NAME_NON_STATIC_APPLE_ARM64=multiapps-plugin-non-static.osxarm64

PLUGIN_NAME_STATIC_WIN_32=multiapps-plugin.win32.exe
PLUGIN_NAME_STATIC_WIN_64=multiapps-plugin.win64.exe
PLUGIN_NAME_STATIC_LINUX_32=multiapps-plugin.linux32
PLUGIN_NAME_STATIC_LINUX_64=multiapps-plugin.linux64
PLUGIN_NAME_STATIC_LINUX_ARM64=multiapps-plugin.linuxarm64
PLUGIN_NAME_STATIC_OSX=multiapps-plugin.osx
PLUGIN_NAME_STATIC_APPLE_ARM64=multiapps-plugin.osxarm64

version=$(<cfg/VERSION)
build "$version" linux 386 $PLUGIN_NAME_NON_STATIC_LINUX_32
build "$version" linux amd64 $PLUGIN_NAME_NON_STATIC_LINUX_64
build "$version" linux arm64 $PLUGIN_NAME_NON_STATIC_LINUX_ARM64
build "$version" windows 386 $PLUGIN_NAME_NON_STATIC_WIN_32
build "$version" windows amd64 $PLUGIN_NAME_NON_STATIC_WIN_64
build "$version" darwin amd64 $PLUGIN_NAME_NON_STATIC_OSX
build "$version" darwin arm64 $PLUGIN_NAME_NON_STATIC_APPLE_ARM64

buildstatic "$version" linux 386 $PLUGIN_NAME_STATIC_LINUX_32
buildstatic "$version" linux amd64 $PLUGIN_NAME_STATIC_LINUX_64
buildstatic "$version" linux arm64 $PLUGIN_NAME_STATIC_LINUX_ARM64
buildstatic "$version" windows 386 $PLUGIN_NAME_STATIC_WIN_32
buildstatic "$version" windows amd64 $PLUGIN_NAME_STATIC_WIN_64
buildstatic "$version" darwin amd64 $PLUGIN_NAME_STATIC_OSX
buildstatic "$version" darwin arm64 $PLUGIN_NAME_STATIC_APPLE_ARM64

mkdir -p $BUILD_FOLDER
createBuildMetadataFiles $version $BUILD_FOLDER
movePluginsToBuildFolder $BUILD_FOLDER
