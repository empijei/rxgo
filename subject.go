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

type Subject[T any] struct {
	completed bool
	cur       int
	subs      map[int]Observer[T]
}

var (
	_ Observable[any] = &Subject[any]{}
	_ Observer[any]   = &Subject[any]{}
)

func NewSubject[T any]() *Subject[T] {
	return &Subject[T]{subs: map[int]Observer[T]{}}
}

func (s *Subject[T]) Next(elm T) {
	if s.completed {
		return
	}
	if s.completed { //Is this correct?
		return
	}
	for _, sub := range s.subs {
		sub.Next(elm)
	}
}

func (s *Subject[T]) Complete() {
	if s.completed {
		return
	}
	s.completed = true
	for _, sub := range s.subs {
		sub.Complete()
	}
	s.subs = map[int]Observer[T]{}
}

func (s *Subject[T]) Error(err error) {
	if s.completed {
		return
	}
	for _, sub := range s.subs {
		sub.Error(err)
	}
}

type subscription struct {
	key   int
	unsub func()
}

func (s subscription) Unsubscribe() {
	s.unsub()
}

func (s *Subject[T]) Subscribe(o Observer[T]) Subscription {
	key := s.cur
	s.subs[key] = o
	s.cur++
	// TODO: If completed already send completion
	return subscription{
		key:   key,
		unsub: func() { delete(s.subs, key); o.Complete() }, // Should unsub cause completion?
	}
}
