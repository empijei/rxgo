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

type BehaviorSubject[T any] struct {
	val T

	completed bool
	cur       int
	subs      map[int]Observer[T]
}

var (
	_ Observable[any] = &BehaviorSubject[any]{}
	_ Observer[any]   = &BehaviorSubject[any]{}
)

func NewBehaviorSubject[T any](val T) *BehaviorSubject[T] {
	return &BehaviorSubject[T]{
		val:  val,
		subs: map[int]Observer[T]{},
	}
}

func (s *BehaviorSubject[T]) Next(elm T) {
	if s.completed {
		return
	}
	s.val = elm
	for _, sub := range s.subs {
		sub.Next(elm)
	}
}

func (s *BehaviorSubject[T]) Complete() {
	s.completed = true
	for _, sub := range s.subs {
		sub.Complete()
	}
	s.subs = map[int]Observer[T]{}
}

func (s *BehaviorSubject[T]) Error(err error) {
	for _, sub := range s.subs {
		sub.Error(err)
	}
}

func (s *BehaviorSubject[T]) Subscribe(o Observer[T]) Subscription {
	key := s.cur
	s.subs[key] = o
	o.Next(s.val)
	s.cur++
	return subscription{
		key:   key,
		unsub: func() { delete(s.subs, key) }, // Should unsub cause completion?
	}
}
