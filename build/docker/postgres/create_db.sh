#!/bin/sh

for db in userapi_accounts userapi_devices pushserver mediaapi syncapi roomserver keyserver federationapi appservice naffka; do
    createdb -U dendrite -O dendrite dendrite_$db
done
