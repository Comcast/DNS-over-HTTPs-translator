#!/bin/sh

# Script to deploy translator as a systemd service
# The script assumes that it is being run from the 
# translator's home directory (current dir).

make clean
make build

sudo mkdir /etc/doh-translator
sudo cp ./config-doh-translator.yaml /etc/doh-translator/

sudo cp ./.build/doh-translator-linux-amd64 /usr/local/bin
sudo cp ./systemd/doh-translator.service /lib/systemd/system/
