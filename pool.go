// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package sizedpool

import "sync"

// TODO: create bucketing functions in vars like power-of-2

// A Pool ...
type Pool struct {
	sync.Mutex
	factory func(n int) interface{}
	sizer   func(x interface{}) int
	bucket  func(n int) uint
	pools   map[uint]*sync.Pool
}

// New ...
func New(factory func(n int) interface{}, sizer func(x interface{}) int, bucket func(n int) uint) *Pool {
	p := &Pool{factory: factory, bucket: bucket, sizer: sizer}
	p.pools = make(map[uint]*sync.Pool)
	return p
}

func (p *Pool) getPool(n int) *sync.Pool {
	i := p.bucket(n)
	sp, ok := p.pools[i]
	if !ok {
		sp = &sync.Pool{New: func() interface{} { return p.factory(n) }}
		p.pools[i] = sp
	}
	return sp
}

// Get ...
func (p *Pool) Get(n int) interface{} {
	return p.getPool(n).Get()
}

// Put ...
func (p *Pool) Put(x interface{}) {
	p.getPool(p.sizer(x)).Put(x)
	return
}
