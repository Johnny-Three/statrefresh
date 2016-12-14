package util

import (
	//"fmt"
	//"reflect"
	"sync"
)

type BeeMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

func NewBeeMap() *BeeMap {
	return &BeeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

//Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Returns true if v is exist in the map.
func (m *BeeMap) GetByValue(vin interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	//fmt.Println("v is ", vin)
	for key, value := range m.bm {

		/*
			fmt.Printf("k is %v,v is %v\n", key, value)

			switch v := value.(type) {
			case string:
				fmt.Println("string")
				fmt.Println(v)
			case int32, int64:
				fmt.Println("int")
				fmt.Println(v)

			default:
				fmt.Println("unknown")
			}

			fmt.Println("===============")
			switch v := vin.(type) {
			case string:
				fmt.Println("string")
				fmt.Println(v)
			case int32, int64:
				fmt.Println("int64")
				fmt.Println(v)
			case int:
				fmt.Println("int")
				fmt.Println(v)
			default:
				fmt.Println("unknown")
			}
		*/

		if value.(int) == vin.(int) {

			//fmt.Println("kkk is ", key)
			return key
		}
	}
	return nil
}

// Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *BeeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

// Items returns all items in safemap.
func (m *BeeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}
