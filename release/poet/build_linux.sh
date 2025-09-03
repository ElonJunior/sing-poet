#!/bin/bash


#当前目录
CURRENT_DIR=$(cd $(dirname $0); pwd)
#项目目录
WORKSPACE_DIR=$(dirname $(dirname "$CURRENT_DIR"))

cd $WORKSPACE_DIR
PROGRAME_NAME="sing-poet"
outfile="sing-box"
outfile1="${CURRENT_DIR}/test/${PROGRAME_NAME}"

# packaging...
rm $outfile $outfile1
# 有需要则修改 Makefile::TAGS_POET2
make poet_build
chmod +x $outfile
mv $outfile $outfile1

cd $CURRENT_DIR/test/
tar zcvf $PROGRAME_NAME-lasted.tar.gz $PROGRAME_NAME
echo -e "\n\t${PROGRAME_NAME} is ready to go!!!"
ls -lha $outfile1 $CURRENT_DIR/test/$PROGRAME_NAME-lasted.tar.gz





