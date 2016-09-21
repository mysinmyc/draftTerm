#!/bin/bash

vHOME=$(dirname $0)/..


if [ -n "${GUEST_PASSWORD}" ]; then
	echo Chaging guest password...
	echo guest:${GUEST_PASSWORD} | chpasswd
fi 

if [ ! -e ${vHOME}/etc/key.pem ]; then
	echo "Generating dummy certificates..."
	cd ${vHOME}/etc
	${vHOME}/bin/generate_cert --host '*'
fi

vCOMMAND="$*"
if [ -z "${vCOMMAND}" ]; then
	vCOMMAND="/sbin/agetty -"
fi


${vHOME}/bin/draftTermd --secure --key ${vHOME}/etc/key.pem --cert ${vHOME}/etc/cert.pem --listen 0.0.0.0:8443 --cmd "${vCOMMAND}"
