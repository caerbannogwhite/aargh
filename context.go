package aargh

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/apache/arrow-go/v18/arrow/memory"
)

type Context struct {
	// StringPool is a pool of strings that are used by the series.
	// This is used to reduce the number of allocations and to allow for fast comparisons.
	StringPool *StringPool

	// Allocator is the Arrow memory allocator used for building Arrow arrays.
	// Defaults to memory.DefaultAllocator (backed by Go GC, no manual Retain/Release needed).
	Allocator memory.Allocator

	threadsNumber  int
	naText         string
	eol            string
	quote          string
	dateTimeFormat string
}

func NewContext() *Context {
	eol := EOL
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	return &Context{
		StringPool:     NewStringPool().SetNaText(NA_TEXT),
		Allocator:      memory.DefaultAllocator,
		threadsNumber:  THREADS_NUMBER,
		naText:         NA_TEXT,
		eol:            eol,
		quote:          QUOTE,
		dateTimeFormat: DATE_TIME_FORMAT,
	}
}

func (ctx *Context) GetThreadsNumber() int {
	return ctx.threadsNumber
}

func (ctx *Context) SetThreadsNumber(n int) *Context {
	ctx.threadsNumber = n
	return ctx
}

func (ctx *Context) GetNaText() string {
	return ctx.naText
}

func (ctx *Context) SetNaText(s string) *Context {
	ctx.StringPool.SetNaText(s)
	ctx.naText = s
	return ctx
}

func (ctx *Context) GetDateTimeFormat() string {
	return ctx.dateTimeFormat
}

func (ctx *Context) SetDateTimeFormat(s string) *Context {
	ctx.dateTimeFormat = s
	return ctx
}

func (ctx *Context) GetEol() string {
	return ctx.eol
}

func (ctx *Context) SetEol(s string) *Context {
	ctx.eol = s
	return ctx
}

func (ctx *Context) GetQuote() string {
	return ctx.quote
}

func (ctx *Context) SetQuote(s string) *Context {
	ctx.quote = s
	return ctx
}

// StringPool is a pool of strings that are used by the series.
// This is used to reduce the number of allocations and to allow for fast comparisons.
type StringPool struct {
	sync.RWMutex
	pool      map[string]*string
	naTextPtr *string
}

func NewStringPool() *StringPool {
	pool := &StringPool{pool: make(map[string]*string)}
	return pool
}

func (sp *StringPool) SetNaText(s string) *StringPool {
	sp.Lock()
	defer sp.Unlock()
	sp.naTextPtr = sp.Put(s)

	return sp
}

// Get returns the address of the string if it exists in the pool, otherwise nil.
func (sp *StringPool) Get(s string) *string {
	if entry, ok := sp.pool[s]; ok {
		return entry
	}
	return nil
}

// Put returns the address of the string if it exists in the pool, otherwise it adds it to the pool and returns its address.
func (sp *StringPool) Put(s string) *string {
	if entry, ok := sp.pool[s]; ok {
		return entry
	}

	// Create a new string and add it to the pool
	addr := &s
	sp.pool[s] = addr
	return addr
}

// PutSync returns the address of the string if it exists in the pool, otherwise it adds it to the pool and returns its address.
// This version is thread-safe.
func (sp *StringPool) PutSync(s string) *string {
	sp.RLock()
	entry, ok := sp.pool[s]
	sp.RUnlock()
	if ok {
		return entry
	}

	sp.Lock()
	defer sp.Unlock()
	if entry, ok := sp.pool[s]; ok {
		// Someone else inserted the string while we were waiting
		return entry
	}

	// Create a new string and add it to the pool
	sp.pool[s] = &s
	return sp.pool[s]
}

func (sp *StringPool) Len() int {
	return len(sp.pool)
}

func (sp *StringPool) ToString() string {
	out := "StringPool["
	for _, v := range sp.pool {
		out += *v + ", "
	}
	out = out[:len(out)-2] + "]"
	return out
}

func (sp *StringPool) debugPrint() {
	for k, v := range sp.pool {
		fmt.Printf("%s: %p\n", k, v)
	}
}
