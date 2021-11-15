killall exchaind
rm -rf multi_run.log
make mainnet WITH_ROCKSDB=true
echo 3 > /proc/sys/vm/drop_caches 
export EXCHAIND_PATH=~/.exchaind
rm -rf ${EXCHAIND_PATH}/
exchaind init multi_run --chain-id exchain-66 --home ${EXCHAIND_PATH}
cp /Users/oker/scf/data/genesis.json ${EXCHAIND_PATH}/config/genesis.json
rm -rf ${EXCHAIND_PATH}/data
cp -rf /Users/oker/scf/data/s0-5810700-rocksdb/data ${EXCHAIND_PATH}
#exchaind replay -d /Users/oker/scf/data/sx-5811000-5813000-rocksdb/data --home ~/.exchaind  --save_block > multi_run.log 2>&1 &
exchaind replay -d /Users/oker/scf/data/sx-5811000-5813000-rocksdb/data --home ~/.exchaind  --paralleled-tx > multi_run.log 2>&1 &
