package pool

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

var intSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
var expectedErr = errors.New("err")

func mapTimesTwo(ctx context.Context, t int, i int) (int, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return t * 2, nil
}

func ExampleMap() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := Map(ctx,
		[]int{1, 2, 3, 4, 5},
		func(ctx context.Context, element int, index int) (string, error) {
			return fmt.Sprintf("element: %d", element), nil
		}, WithWorkers(4))
	if err != nil {
		// ...
	}
	fmt.Printf("%v", result)
	// Output: [element: 1 element: 2 element: 3 element: 4 element: 5]
}

func TestMap_Success(t *testing.T) {
	actual, err := Map(context.Background(), intSlice, mapTimesTwo, WithWorkers(4))
	assertNoError(t, err)
	assertEqual(t, []int{2, 4, 6, 8, 10, 12, 14, 16, 18}, actual)
}

func TestMap_WithError(t *testing.T) {
	expectedErr := errors.New("err")
	actual, err := Map(context.Background(), intSlice, func(ctx context.Context, t int, i int) (int, error) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		if t > 7 {
			return 0, expectedErr
		}
		return t * 2, nil
	}, WithWorkers(4))
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	assertEqual(t, []int{}, actual)
}

func TestMap_WithCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	actual, err := Map(ctx, intSlice, mapTimesTwo, WithWorkers(4))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	assertEqual(t, []int{}, actual)
}
