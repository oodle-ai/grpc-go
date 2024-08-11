/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package leakcheck contains functions to check leaked goroutines.
//
// Call "defer leakcheck.Check(t)" at the beginning of tests.
package leakcheck

import (
	"runtime"
	"sort"
	"strings"
	"time"
)

<<<<<<< HEAD
=======
// failTestsOnLeakedBuffers is a special flag that will cause tests to fail if
// leaked buffers are detected, instead of simply logging them as an
// informational failure. This can be enabled with the "checkbuffers" compile
// flag, e.g.:
//
//	go test -tags=checkbuffers
var failTestsOnLeakedBuffers = false

func init() {
	defaultPool := mem.DefaultBufferPool()
	globalPool.Store(&defaultPool)
	(internal.SetDefaultBufferPoolForTesting.(func(mem.BufferPool)))(&globalPool)
}

var globalPool swappableBufferPool

type swappableBufferPool struct {
	atomic.Pointer[mem.BufferPool]
}

func (b *swappableBufferPool) Get(length int) []byte {
	return (*b.Load()).Get(length)
}

func (b *swappableBufferPool) Put(buf *[]byte) {
	(*b.Load()).Put(buf)
}

// SetTrackingBufferPool upgrades the default buffer pool in the mem package to
// one that tracks where buffers are allocated. CheckTrackingBufferPool should
// then be invoked at the end of the test to validate that all buffers pulled
// from the pool were returned.
func SetTrackingBufferPool(efer Errorfer) {
	if !mem.PoolingEnabled {
		return
	}
	newPool := mem.BufferPool(&trackingBufferPool{
		pool:             *globalPool.Load(),
		efer:             efer,
		allocatedBuffers: make(map[*byte][]uintptr),
	})
	globalPool.Store(&newPool)
}

// CheckTrackingBufferPool undoes the effects of SetTrackingBufferPool, and fails
// unit tests if not all buffers were returned. It is invalid to invoke this
// method without previously having invoked SetTrackingBufferPool.
func CheckTrackingBufferPool() {
	if !mem.PoolingEnabled {
		return
	}
	p := (*globalPool.Load()).(*trackingBufferPool)
	p.lock.Lock()
	defer p.lock.Unlock()

	globalPool.Store(&p.pool)

	type uniqueTrace struct {
		stack []uintptr
		count int
	}

	var uniqueTraces []uniqueTrace
	for _, stack := range p.allocatedBuffers {
		idx, ok := slices.BinarySearchFunc(uniqueTraces, stack, func(trace uniqueTrace, stack []uintptr) int {
			return slices.Compare(trace.stack, stack)
		})
		if !ok {
			uniqueTraces = slices.Insert(uniqueTraces, idx, uniqueTrace{stack: stack})
		}
		uniqueTraces[idx].count++
	}

	for _, ut := range uniqueTraces {
		frames := runtime.CallersFrames(ut.stack)
		var trace strings.Builder
		for {
			f, ok := frames.Next()
			if !ok {
				break
			}
			trace.WriteString(f.Function)
			trace.WriteString("\n\t")
			trace.WriteString(f.File)
			trace.WriteString(":")
			trace.WriteString(strconv.Itoa(f.Line))
			trace.WriteString("\n")
		}
		format := "%d allocated buffers never freed:\n%s"
		args := []any{ut.count, trace.String()}
		if failTestsOnLeakedBuffers {
			p.efer.Errorf(format, args...)
		} else {
			p.efer.Logf("WARNING "+format, args...)
		}
	}
}

type trackingBufferPool struct {
	pool mem.BufferPool
	efer Errorfer

	lock             sync.Mutex
	allocatedBuffers map[*byte][]uintptr
}

func (p *trackingBufferPool) Get(length int) []byte {
	p.lock.Lock()
	defer p.lock.Unlock()
	if length == 0 {
		return nil
	}

	buf := p.pool.Get(length)

	var stackBuf [16]uintptr
	var stack []uintptr
	skip := 2
	for {
		n := runtime.Callers(skip, stackBuf[:])
		stack = append(stack, stackBuf[:n]...)
		if n < len(stackBuf) {
			break
		}
		skip += len(stackBuf)
	}
	p.allocatedBuffers[unsafe.SliceData(buf)] = stack

	return buf
}

func (p *trackingBufferPool) Put(buf *[]byte) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(*buf) == 0 {
		return
	}

	key := unsafe.SliceData(*buf)
	if _, ok := p.allocatedBuffers[key]; !ok {
		p.efer.Errorf("Unknown buffer freed:\n%s", string(debug.Stack()))
	} else {
		delete(p.allocatedBuffers, key)
	}
	p.pool.Put(buf)
}

>>>>>>> 35f35e4d (Enable test failures with tag)
var goroutinesToIgnore = []string{
	"testing.Main(",
	"testing.tRunner(",
	"testing.(*M).",
	"runtime.goexit",
	"created by runtime.gc",
	"created by runtime/trace.Start",
	"interestingGoroutines",
	"runtime.MHeap_Scavenger",
	"signal.signal_recv",
	"sigterm.handler",
	"runtime_mcall",
	"(*loggingT).flushDaemon",
	"goroutine in C code",
	// Ignore the http read/write goroutines. gce metadata.OnGCE() was leaking
	// these, root cause unknown.
	//
	// https://github.com/grpc/grpc-go/issues/5171
	// https://github.com/grpc/grpc-go/issues/5173
	"created by net/http.(*Transport).dialConn",
}

// RegisterIgnoreGoroutine appends s into the ignore goroutine list. The
// goroutines whose stack trace contains s will not be identified as leaked
// goroutines. Not thread-safe, only call this function in init().
func RegisterIgnoreGoroutine(s string) {
	goroutinesToIgnore = append(goroutinesToIgnore, s)
}

func ignore(g string) bool {
	sl := strings.SplitN(g, "\n", 2)
	if len(sl) != 2 {
		return true
	}
	stack := strings.TrimSpace(sl[1])
	if strings.HasPrefix(stack, "testing.RunTests") {
		return true
	}

	if stack == "" {
		return true
	}

	for _, s := range goroutinesToIgnore {
		if strings.Contains(stack, s) {
			return true
		}
	}

	return false
}

// interestingGoroutines returns all goroutines we care about for the purpose of
// leak checking. It excludes testing or runtime ones.
func interestingGoroutines() (gs []string) {
	buf := make([]byte, 2<<20)
	buf = buf[:runtime.Stack(buf, true)]
	for _, g := range strings.Split(string(buf), "\n\n") {
		if !ignore(g) {
			gs = append(gs, g)
		}
	}
	sort.Strings(gs)
	return
}

// Errorfer is the interface that wraps the Errorf method. It's a subset of
// testing.TB to make it easy to use Check.
type Errorfer interface {
	Errorf(format string, args ...any)
}

func check(efer Errorfer, timeout time.Duration) {
	// Loop, waiting for goroutines to shut down.
	// Wait up to timeout, but finish as quickly as possible.
	deadline := time.Now().Add(timeout)
	var leaked []string
	for time.Now().Before(deadline) {
		if leaked = interestingGoroutines(); len(leaked) == 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	for _, g := range leaked {
		efer.Errorf("Leaked goroutine: %v", g)
	}
}

// Check looks at the currently-running goroutines and checks if there are any
// interesting (created by gRPC) goroutines leaked. It waits up to 10 seconds
// in the error cases.
func Check(efer Errorfer) {
	check(efer, 10*time.Second)
}
