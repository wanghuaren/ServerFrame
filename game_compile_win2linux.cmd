chcp 65001

echo off

rmdir server /q /s
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

go env -w GOOS=windows

mkdir server

copy gamegate\gamegate server
copy gamegate\gateconf.ini server

copy gameserver\gameserver server
copy gameserver\serverconf.ini server

copy gamedb\gamedb server
copy gamedb\dbconf.ini server

del gamegate\gamegate
del gameserver\gameserver
del gamedb\gamedb

echo "编译完成,按任意键结束"

pause