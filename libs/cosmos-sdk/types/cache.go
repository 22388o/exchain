package types

import (
	"encoding/hex"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"sync"
)

type Cache struct {
	parent    *Cache
	storage   map[ethcmn.Address]map[ethcmn.Hash][]byte
	storageMu sync.Mutex
	cdc       *codec.Codec
}

func NewCache(parent *Cache) *Cache {
	return &Cache{
		parent:  parent,
		storage: make(map[ethcmn.Address]map[ethcmn.Hash][]byte, 0),
	}
}

func (cache *Cache) Copy() *Cache {
	return NewCache(cache)
}

func (c *Cache) SetCDC(cdc *codec.Codec) {
	c.cdc = cdc
}
func (c *Cache) Update(ms CacheWrap) {
	var wg sync.WaitGroup
	ms.IteratorCache(func(key, value []byte, isDirty bool) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			typeKey := hex.EncodeToString(key)[:2]
			switch typeKey {
			case "07":
				//if len(key) == 21 {
				//	c.updateGetValidatorAccumulatedCommission(key, value)
				//}
			case "05":
				if len(key) == 1+20+32 {
					c.updateStorage(key, value)
				}
			}
		}()
		return true
	})
	wg.Wait()
}

//func (c *Cache) updateGetValidatorAccumulatedCommission(key []byte, value []byte) {
//	commis := new(SysCoins)
//	c.cdc.MustUnmarshalBinaryLengthPrefixed(value, commis)
//	c.validatorCommission[string(key[2:])] = *commis
//}
//
//func (c *Cache) GetValidatorAccumulatedCommission(key ValAddress) (SysCoins, bool) {
//	data, ok := c.validatorCommission[string(key)]
//	return data, ok
//}

func (c *Cache) updateStorage(key, value []byte) {
	addr := ethcmn.BytesToAddress(key[1:21])
	stateKey := ethcmn.BytesToHash(key[21:])

	c.storageMu.Lock()
	defer c.storageMu.Unlock()
	if _, ok := c.storage[addr]; !ok {
		c.storage[addr] = make(map[ethcmn.Hash][]byte, 0)
	}
	c.storage[addr][stateKey] = value
}

func (c *Cache) GetStorage(addr ethcmn.Address, key ethcmn.Hash) ([]byte, bool) {
	//return nil, false
	if _, ok := c.storage[addr]; !ok {
		return nil, false
	}

	data, ok := c.storage[addr][key]
	if ok {
		//fmt.Println("haveStorage", addr.String(), key.String())
	}

	return data, ok
}
