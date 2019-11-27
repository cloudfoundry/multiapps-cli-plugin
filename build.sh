#!/bin/bash -eux

function build() {
    local version=$1
    local platform=$2
    local arch=$3
    local plugin_name=$4

    echo calling to build for $platform $arch
    CGO_ENABLED=0 GOOS=$platform GOARCH=$arch go build \
        -ldflags "-X main.Version=${version}" \
        -o ${plugin_name}
}

function copyPluginsWithOldNames() {
    cp $PLUGIN_NAME_WIN_64 $OLD_PLUGIN_NAME_WIN_64
    cp $PLUGIN_NAME_LINUX_64 $OLD_PLUGIN_NAME_LINUX_64
    cp $PLUGIN_NAME_OSX $OLD_PLUGIN_NAME_OSX
}

function movePluginsToBuildFolder() {
    local folder=$1
    mv $PLUGIN_NAME_WIN_32 $folder
    mv $PLUGIN_NAME_WIN_64 $folder
    mv $PLUGIN_NAME_LINUX_32 $folder
    mv $PLUGIN_NAME_LINUX_64 $folder
    mv $PLUGIN_NAME_OSX $folder
}

function createBuildMetadataFiles() {
    local version=$1
    local folder=$2
	echo -n "v${version}" > ${folder}/version
	grep -Pzoa "(?s)## v${version}(.*?)##" CHANGELOG.md | grep -va "##" | tr -s '\n' '\n' > ${folder}/changelog
}

if [ ! -f "build.sh" ] | [ ! -f "LICENSE" ]; then
   echo "Please run with the project root as working directory";
   exit 1 ; 
fi

BUILD_FOLDER=build
PLUGIN_NAME_WIN_32=multiapps-plugin.win32
PLUGIN_NAME_WIN_64=multiapps-plugin.win64
PLUGIN_NAME_LINUX_32=multiapps-plugin.linux32
PLUGIN_NAME_LINUX_64=multiapps-plugin.linux64
PLUGIN_NAME_OSX=multiapps-plugin.osx
OLD_PLUGIN_NAME_WIN_64=mta_plugin_windows_amd64.exe
OLD_PLUGIN_NAME_LINUX_64=mta_plugin_linux_amd64
OLD_PLUGIN_NAME_OSX=mta_plugin_darwin_amd64

version=$(<cfg/VERSION)
build $version linux 386 $PLUGIN_NAME_LINUX_32
build $version linux amd64 $PLUGIN_NAME_LINUX_64
build $version windows 386 $PLUGIN_NAME_WIN_32
build $version windows amd64 $PLUGIN_NAME_WIN_64
build $version darwin amd64 $PLUGIN_NAME_OSX

mkdir -p $BUILD_FOLDER
createBuildMetadataFiles $version $BUILD_FOLDER
movePluginsToBuildFolder $BUILD_FOLDER
cp Dockerfile $BUILD_FOLDER/
cd $BUILD_FOLDER
copyPluginsWithOldNames