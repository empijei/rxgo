// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rxgo

type ReplaySubject[T any] struct {
	size int
	vals []T

	completed bool
	cur       int
	subs      map[int]Observer[T]
}

var (
	_ Observable[any] = &ReplaySubject[any]{}
	_ Observer[any]   = &ReplaySubject[any]{}
)

func NewReplaySubject[T any](size int) *ReplaySubject[T] {
	rs := &ReplaySubject[T]{
		size: size,
		subs: map[int]Observer[T]{},
	}
	return rs
}

func (s *ReplaySubject[T]) Next(elm T) {
	if s.completed { // Is this correct?
		return
	}
	s.vals = append(s.vals, elm)
	if len(s.vals) > s.size {
		s.vals = s.vals[1:]
	}
	for _, sub := range s.subs {
		sub.Next(elm)
	}
}

func (s *ReplaySubject[T]) Complete() {
	s.completed = true
	for _, sub := range s.subs {
		sub.Complete()
	}
	s.subs = map[int]Observer[T]{}
}

func (s *ReplaySubject[T]) Error(err error) {
	// BOOKMARK
	for _, sub := range s.subs {
		sub.Error(err)
	}
}

func (s *ReplaySubject[T]) Subscribe(o Observer[T]) Subscription {
	key := s.cur
	s.subs[key] = o
	o.Next(s.val)
	s.cur++
	return subscription{
		key:   key,
		unsub: func() { delete(s.subs, key) }, // Should unsub cause completion?
	}
}
