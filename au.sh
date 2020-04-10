#!/bin/sh

/usr/sbin/addgroup --group $2
/usr/sbin/adduser --quiet --disabled-password --gecos "" --disabled-login $1 --ingroup $2

(/usr/bin/echo "$3"; /usr/bin/echo "$3")  | /usr/bin/smbpasswd -s -a $1

/usr/bin/cat <<EOF >> /etc/samba/smb.conf

[$1]
    comment = whatever
    path = /home/$1
    browsable = yes
    read only = no
    guest ok = no

EOF

/usr/bin/cat <<EOF >> /etc/ppp/chap-secrets

$1 pptpd $3 *

EOF

/usr/bin/docker exec docker-jitsi-meet_prosody_1 /usr/bin/prosodyctl --config /config/prosody.cfg.lua register $1 meet.jitsi $3