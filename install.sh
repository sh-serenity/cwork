#!/bin/sh
apt-get update
apt-get -y install gpg
set -ex;
    key='A4A9406876FCBD3C456770C88C718D3B5072E1F5'; \
    export GNUPGHOME="$(mktemp -d)"; \
    gpg --batch --keyserver ha.pool.sks-keyservers.net --recv-keys "$key"; \
    gpg --batch --export "$key" > /etc/apt/trusted.gpg.d/mysql.gpg; \
    gpgconf --kill all; \
    rm -rf "$GNUPGHOME"; \
    apt-key list > /dev/null

echo "deb http://repo.mysql.com/apt/debian/ buster mysql-5.7" > /etc/apt/sources.list.d/mysql.list
apt-get update && apt-get install -y mysql-server nginx ppp pptpd git samba sudo
/usr/sbin/iptables -A INPUT -s $2.0/24 -j ACCEPT
/usr/sbin/iptables -A INPUT -p tcp -m multiport --dports 137,138,139,445 -j DROP

/usr/bin/cat <<EOF >> /etc/rc.local

/usr/sbin/iptables -A INPUT -s $2 -j ACCEPT
/usr/sbin/iptables -A INPUT -p tcp -m multiport --dports 137,138,139,445 -j DROP

EOF
/usr/bin/cat <<EOF >> /etc/pptpd.conf
localip  $2.1
remoteip $2.10-250
EOF
/etc/init.d/pptpd restart
mkdir -p /opt/go/cwork
/usr/bin/cat <<EOF > /etc/nginx/conf.d/cwork.conf
server {
        server_name $1;
        root /opt/go/cwork;
        location = /favicon.ico {
                log_not_found off;
                access_log off;
        }
        location = /robots.txt {
                allow all;
                log_not_found off;
                access_log off;
        }
        location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        
                expires 5m;
                log_not_found off;
        }

        location / {
                fastcgi_intercept_errors on;
                include fastcgi_params;
                fastcgi_pass 127.0.0.1:9001;
        }
}
EOF

/usr/bin/cat <<EOF >> /etc/sudoers

go      ALL=(ALL) NOPASSWD: /opt/go/cwork/au.sh
EOF

cd /opt/go
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/opt/go
/usr/local/go/bin/go get github.com/go-sql-driver/mysql
/usr/local/go/bin/go get github.com/gorilla/sessions

git clone https://github.com/sh-serenity/cwork.git
mv /opt/go/cwork/adduser.conf /etc

/usr/sbin/addgroup --group go
/usr/sbin/adduser --quiet --disabled-password --gecos "" --disabled-login go --ingroup go
cd cwork
/usr/bin/cat <<EOF > /opt/go/cwork/tmpl/menu.html
{{ define "menu" }}

<div><center><a href="/reg2/">Зарегистрировать юзера</a><a href="https://$2.1:8443/">Початится</a><a href="/exit/">Выйти</a></center></div>

{{ end }}
EOF

mysql -p  < shwork.sql
go build
chown -R go:go /opt/go
/etc/init.d/nginx restart

apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
apt-key fingerprint 0EBFCD88
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get -y install docker-ce docker-ce-cli containerd.io

cd /opt
#git clone https://github.com/jitsi/docker-jitsi-meet.git
#cp /opt/docker-jitsi-meet/env.example /opt/docker-jitsi-meet/.env
#/usr/bin/replace '#DOCKER_HOST_ADDRESS=192.168.1.1' DOCKER_HOST_ADDRESS=$2.1  -- /opt/docker-jitsi-meet/.env
wget https://github.com/jitsi/docker-jitsi-meet/archive/stable-3547.tar.gz
tar xvzf stable-3547.tar.gz
mv docker-jitsi-meet-stable-3547 docker-jitsi-meet

/usr/bin/cat <<EOF > /opt/docker-jitsi-meet/.env
JICOFO_COMPONENT_SECRET=azwsdcrf321
JICOFO_AUTH_PASSWORD=azwsdcrf321
JVB_AUTH_PASSWORD=azwsdcrf321
JIGASI_XMPP_PASSWORD=azwsdcrf321
JIBRI_RECORDER_PASSWORD=azwsdcrf321
JIBRI_XMPP_PASSWORD=azwsdcrf321
CONFIG=~/.jitsi-meet-cfg
HTTP_PORT=8000
HTTPS_PORT=8443
TZ=Europe/Kiev
#PUBLIC_URL=https://meet.example.com
DOCKER_HOST_ADDRESS=$2.1
ENABLE_AUTH=1
AUTH_TYPE=internal
XMPP_DOMAIN=meet.jitsi
XMPP_SERVER=xmpp.meet.jitsi
XMPP_BOSH_URL_BASE=http://xmpp.meet.jitsi:5280
XMPP_AUTH_DOMAIN=auth.meet.jitsi
XMPP_MUC_DOMAIN=muc.meet.jitsi
XMPP_INTERNAL_MUC_DOMAIN=internal-muc.meet.jitsi
XMPP_GUEST_DOMAIN=guest.meet.jitsi
XMPP_MODULES=
XMPP_MUC_MODULES=
XMPP_INTERNAL_MUC_MODULES=
JVB_BREWERY_MUC=jvbbrewery
JVB_AUTH_USER=jvb
JVB_STUN_SERVERS=meet-jit-si-turnrelay.jitsi.net:443
JVB_PORT=10000
JVB_TCP_HARVESTER_DISABLED=true
JVB_TCP_PORT=4443
#JVB_ENABLE_APIS=rest,colibri
JICOFO_AUTH_USER=focus
JIGASI_XMPP_USER=jigasi
JIGASI_BREWERY_MUC=jigasibrewery
JIGASI_PORT_MIN=20000
JIGASI_PORT_MAX=20050
#JIGASI_ENABLE_SDES_SRTP=1
XMPP_RECORDER_DOMAIN=recorder.meet.jitsi
JIBRI_RECORDER_USER=recorder
JIBRI_RECORDING_DIR=/config/recordings
JIBRI_FINALIZE_RECORDING_SCRIPT_PATH=/config/finalize.sh
JIBRI_XMPP_USER=jibri
JIBRI_BREWERY_MUC=jibribrewery
JIBRI_PENDING_TIMEOUT=90
JIBRI_STRIP_DOMAIN_JID=muc
JIBRI_LOGS_DIR=/config/logs
RESTART_POLICY=unless-stopped
JVB_TCP_MAPPED_PORT=4443
EOF


curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

cd /opt/docker-jitsi-meet
docker-compose up -d

cp /opt/go/cwork/cwork.service /etc/systemd/system/
/usr/bin/systemctl start cwork

