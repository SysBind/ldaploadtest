#!/bin/bash

UID_=1000
GID_=1000

for i in {1..2000}; do
    USERNAME="demo$i" UID_=$((UID + i)) GID_=$((GID + i)) envsubst < account.ldif.tmpl > account.ldif
    # add record via ldapadd
    ldapadd -x -H ldap://localhost -D "cn=admin,dc=sysbind,dc=test" -f account.ldif -w p4ssw0rd
done

