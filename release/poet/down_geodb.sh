#!/bin/bash
[ $(id -u) != "0" ] && { echo "Error: You must be root to run this script!"; exit 1; }


echo "Download geoip and geosite..."
#当前目录
POET_DIR="/etc/sing-poet"

if [ ! -d "$POET_DIR" ]; then
    mkdir -p "$POET_DIR"
fi

#下载文件
wget --no-check-certificate --show-progress -O "${POET_DIR}/geoip.db" -c -N https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db
wget --no-check-certificate --show-progress -O "${POET_DIR}/geosite.db" -c -N https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db

ls -lha "${POET_DIR}"/*.db


