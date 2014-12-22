package main

import (
    "container/heap"
    "time"
    "bytes"
    "strconv"
)

type CacheItem struct {
    timestamp int64
    data *bytes.Reader
    index int
    path string
}

type ImageCache struct {
    items []*CacheItem
    imap map[string]*CacheItem
}

var start int64

func NewImageCache() (*ImageCache) {
    ic := &ImageCache{imap: make( map[string]*CacheItem, 100)}
    heap.Init( ic )
    return ic
}

func NewCacheItem(path string, data *bytes.Reader) (*CacheItem) {
    if start == 0 {
        //FIXME: remove this
        start = time.Now().UnixNano()
    }
    return &CacheItem{ data: data, path: path, timestamp: time.Now().UnixNano() }
}

func (pq ImageCache) Find( path string ) (*bytes.Reader, bool)  {
    if item, ok := pq.imap[path]; ok {
        return item.data, ok
    }
    return nil, false
}

func (pq ImageCache) Len() int { return len(pq.items) }

func (pq ImageCache) Top() *CacheItem {
    n := len(pq.items)
    if n == 0 {
        return nil
    }
    return pq.items[0]
}

func (pq ImageCache) Less( i, j int ) bool {
    return pq.items[i].timestamp < pq.items[j].timestamp
}

func (pq ImageCache) Swap( i, j int ) {
    pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
    pq.items[i].index = i
    pq.items[j].index = j
}

func (pq ImageCache) GetPaths() ([]string) {
    paths := make([]string, len(pq.items))
    for i, v := range pq.items {
        paths[i] = v.path + ":" + strconv.FormatInt(v.timestamp - start, 10)
    }

    return paths
}

func (pq *ImageCache) Push( x interface{} ) {
    num := len(pq.items)
    item := x.(*CacheItem)
    item.index = num
    pq.items = append( pq.items, item )
    pq.imap[item.path] = item
}

func (pq *ImageCache) Pop() interface{} {
    old := pq.items
    n := len(old)
    item := old[n-1]
    item.index = -1 // for safety
    pq.items = old[0 : n-1]
    delete( pq.imap, item.path )
    return item
}

func (pq *ImageCache) Update( path string ) {
    item, ok := pq.imap[ path ]
    if !ok { return }
    heap.Remove( pq, item.index )
    item.timestamp = time.Now().UnixNano()
    //heap.Fix( pq, item.index )
    heap.Push(pq, item)
}
