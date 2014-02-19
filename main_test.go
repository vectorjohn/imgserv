package main

import (
    "testing"
    "container/heap"
    "time"
)

func TestQueue( t *testing.T ) {
    q := NewImageCache()
    heap.Push( q, NewCacheItem( "item1", nil ) )
    time.Sleep( time.Duration(1) * time.Millisecond )
    heap.Push( q, NewCacheItem( "item2", nil ) )
    time.Sleep( time.Duration(1) * time.Millisecond )
    heap.Push( q, NewCacheItem( "item3", nil ) )
    time.Sleep( time.Duration(1) * time.Millisecond )
    heap.Push( q, NewCacheItem( "item4", nil ) )
    time.Sleep( time.Duration(1) * time.Millisecond )
    t.Log( q.GetPaths() )

    q.Update( "item1" )
    t.Log( q.GetPaths() )

    t.Log( q.Top().path )
    heap.Pop( q )
    t.Log( q.GetPaths() )

    t.Log( q.Top().path )
    heap.Pop( q )
    t.Log( q.GetPaths() )

    t.Log( q.Top().path )
    heap.Pop( q )
    t.Log( q.GetPaths() )
}
