#!/bin/bash


#当前目录
CURRENT_DIR=$(cd $(dirname $0); pwd)
#项目目录
WORKSPACE_DIR=$(dirname $(dirname "$CURRENT_DIR"))

cd $WORKSPACE_DIR
PROGRAME_NAME="sing-poet"
outfile="sing-box"
outfile1="${CURRENT_DIR}/test/${PROGRAME_NAME}"

# check tuic service bug
# GOMODIR=$(go env GOMODCACHE)
# bugPkg="${GOMODIR}/github.com/sagernet/sing-quic"
# srcPkgPtah=$(ls -td "${bugPkg}"@* | head -n 1)
# bugStr=$(grep "authentication: unknown user" $srcPkgPtah/tuic/service.go | grep -v string | grep UUID)
# if [[ -n $bugStr ]]; then
#     echo "${srcPkgPtah}/tuic/service.go"
#     echo -e "\n\tquic src code has bug, please check!!!"
#     exit 1
# fi





# packaging...
rm $outfile $outfile1
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -trimpath -ldflags="-s -w -buildid=" -tags with_gvisor,with_dhcp,with_quic,with_acme,with_v2ray_api -o $outfile ./cmd/sing-box
# 有需要则修改 Makefile::TAGS_POET2
make poet_build
chmod +x $outfile
mv $outfile $outfile1

cd $CURRENT_DIR/test/
tar zcvf $PROGRAME_NAME-lasted.tar.gz $PROGRAME_NAME
echo -e "\n\t${PROGRAME_NAME} is ready to go!!!"
ls -lha $outfile1 $PROGRAME_NAME-lasted.tar.gz





