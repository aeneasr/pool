// Copyright [yyyy] [name of copyright owner]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"context"
	"runtime"
	"sync"
)

// Options defines how concurrent tasks should be executed.
type Options struct {
	// Workers is the maximum amount of workers processing tasks.
	Workers int
}

func newOptions(o []option) *Options {
	opt := &Options{Workers: runtime.NumCPU()}

	for _, f := range o {
		f(opt)
	}

	if opt.Workers < 1 {
		opt.Workers = 1
	}

	return opt
}

type option func(o *Options)

// WithWorkers defines how many workers will be used to process tasks at most.
func WithWorkers(count uint) option {
	return func(o *Options) {
		o.Workers = int(count)
	}
}

// Map manipulates a slice and transforms it to a slice of another type concurrently.
func Map[T any, R any](ctx context.Context, collection []T, iteratee func(context.Context, T, int) (R, error), opts ...option) ([]R, error) {
	o := newOptions(opts)
	result := make([]R, len(collection))
	ic := make(chan int)
	ec := make(chan error)

	var wg sync.WaitGroup
	for i := 0; i < o.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range ic {

				r, err := iteratee(ctx, collection[index], index)
				if err != nil {
					ec <- err
					return
				}

				result[index] = r
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ec)
	}()

	for i := range collection {
		ic <- i
	}

	close(ic)

	for {
		select {
		case err, ok := <-ec:
			if ok {
				return nil, err
			}

			return result, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
