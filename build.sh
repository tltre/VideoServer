#! /bin/bash

# build Web UI
cd D:/gopath/src/VideoServer/web || exit
go install
cp D:/gopath/bin/web D:/gopath/bin/video_server_web_ui/web
cp -R D:/gopath/src/VideoServer/templates D:/gopath/bin/video_server_web_ui/