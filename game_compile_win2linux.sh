#!/usr/bin/env bash
# -*- coding: utf-8 -*-
stty -echo

rm -rf server

go env -w GOOS=linux

cd gamegate
go build
cd ..

cd gameserver
go build
cd ..

cd gamedb
go build
cd ..

go env -w GOOS=darwin

mkdir server

cp gamegate/gamegate server
cp gamegate/gateconf.ini server

cp gameserver/gameserver server
cp gameserver/serverconf.ini server

cp gamedb/gamedb server
cp gamedb/dbconf.ini server

rm gamegate/gamegate
rm gameserver/gameserver
rm gamedb/gamedb

echo "编译完成"

sleep 3
