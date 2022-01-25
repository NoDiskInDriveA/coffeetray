#! /usr/bin/env bash
# no path correction, so run with proj dir as cwd

mkdir -p build/coffeetray.app/Contents/Resources build/coffeetray.app/Contents/MacOS

mkdir -p icon.iconset
rm icon.iconset/*

sips -z 16 16 resources/Coffeetray.png --out icon.iconset/icon_16x16.png

sips -z 32 32 resources/Coffeetray.png --out icon.iconset/icon_16x16@2x.png
sips -z 32 32 resources/Coffeetray.png --out icon.iconset/icon_32x32.png

sips -z 64 64 resources/Coffeetray.png --out icon.iconset/icon_32x32@2x.png
sips -z 64 64 resources/Coffeetray.png --out icon.iconset/icon_64x64.png

sips -z 128 128 resources/Coffeetray.png --out icon.iconset/icon_64x64@2x.png
sips -z 128 128 resources/Coffeetray.png --out icon.iconset/icon_128x128.png

sips -z 256 256 resources/Coffeetray.png --out icon.iconset/icon_128x128@2x.png
sips -z 256 256 resources/Coffeetray.png --out icon.iconset/icon_256x256.png

sips -z 512 512 resources/Coffeetray.png --out icon.iconset/icon_256x256@2x.png
sips -z 512 512 resources/Coffeetray.png --out icon.iconset/icon_512x512.png

sips -z 1024 1024 resources/Coffeetray.png --out icon.iconset/icon_512x512@2x.png
sips -z 1024 1024 resources/Coffeetray.png --out icon.iconset/icon_1024x1024.png

iconutil -c icns -o build/coffeetray.app/Contents/Resources/icon.icns icon.iconset/
cp resources/EyeCon* build/coffeetray.app/Contents/Resources/
cp resources/Info.plist build/coffeetray.app/Contents/

go build -v -o build/coffeetray.app/Contents/MacOS/ cmd/coffeetray.go

rm -rf icon.iconset/