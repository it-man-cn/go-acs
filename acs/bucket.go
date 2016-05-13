package main

import (
	"sync"
)

type BucketOptions struct {
	ChannelSize int
}

// Bucket is a channel holder.
type Bucket struct {
	cLock    sync.RWMutex        // protect the channels for chs
	chs      map[string]*Channel // map sn to a channel
	boptions BucketOptions
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(boptions BucketOptions) (b *Bucket) {
	b = new(Bucket)
	b.boptions = boptions
	b.chs = make(map[string]*Channel, boptions.ChannelSize)
	return
}

// Put put a channel according with sub key.
func (b *Bucket) Put(key string, ch *Channel) {
	b.cLock.Lock()
	b.chs[key] = ch
	b.cLock.Unlock()
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(key string) {
	b.cLock.Lock()
	if _, ok := b.chs[key]; ok {
		delete(b.chs, key)
	}
	b.cLock.Unlock()
}
