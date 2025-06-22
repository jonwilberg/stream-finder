#!/bin/bash
set -e

apt-get update
apt-get upgrade -y

apt-get install -y apt-transport-https openjdk-11-jdk wget
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo 'deb https://artifacts.elastic.co/packages/7.x/apt stable main' > /etc/apt/sources.list.d/elastic-7.x.list
apt-get update && apt-get install -y elasticsearch

# Allow Elasticsearch port through firewall
ufw allow 9200/tcp

# Set the bootstrap password
echo "${ELASTICSEARCH_PASSWORD}" | /usr/share/elasticsearch/bin/elasticsearch-keystore add --stdin --force bootstrap.password

systemctl enable elasticsearch.service
systemctl start elasticsearch.service 