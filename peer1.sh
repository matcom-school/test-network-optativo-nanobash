#!/usr/bin/env sh
set -eu


CCADDR="host.docker.internal"

echo ${CCADDR}
# look for binaries in local dev environment /build/bin directory and then in local samples /bin directory
export PATH="${PWD}"/../../fabric/build/bin:"${PWD}"/../bin:"$PATH"
export FABRIC_CFG_PATH="${PWD}"/../config

export FABRIC_LOGGING_SPEC=debug:cauthdsl,policies,msp,grpc,peer.gossip.mcs,gossip,leveldbhelper=info
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ID=peer0.org1.example.com
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LISTENADDRESS=127.0.0.1:7051
export CORE_PEER_CHAINCODEADDRESS="${CCADDR}":7052
export CORE_PEER_CHAINCODELISTENADDRESS=127.0.0.1:7052
# bootstrap peer is the other peer in the same org
export CORE_PEER_GOSSIP_BOOTSTRAP=127.0.0.1:7051
export CORE_PEER_GOSSIP_EXTERNALENDPOINT=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/msp
export CORE_PEER_FILESYSTEMPATH="${PWD}"/data/peer0.org1.example.com
export CORE_LEDGER_SNAPSHOTS_ROOTDIR="${PWD}"/data/peer0.org1.example.com/snapshots
# used in metrics
export CORE_METRICS_PROVIDER=prometheus
export CORE_OPERATIONS_LISTENADDRESS=127.0.0.1:8446

# uncomment the lines below to utilize couchdb state database, when done with the environment you can stop the couchdb container with "docker rm -f couchdb1"
export CORE_LEDGER_STATE_STATEDATABASE=CouchDB
export CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=127.0.0.1:5984
export CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=portainer
export CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=N@rut0
# remove the container only if it exists
CONTAINER_NAME=$(docker ps -aqf "NAME=worldstate_org1")
if [ -z "$CONTAINER_NAME" -o "$CONTAINER_NAME" == " " ]; then
	 echo "---- No world-state container available for deletion ----"
else
 	 docker rm -f $CONTAINER_NAME
fi
docker run --publish 5984:5984 --detach -e COUCHDB_USER=portainer -e COUCHDB_PASSWORD=N@rut0 --name worldstate_org1 couchdb:3.1.1

# start peer
peer node start
