package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShift(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	v := Shift(&l)
	assert.NotNil(t, v)
	assert.Equal(t, 1, *v)
	assert.Equal(t, []int{2, 3, 4, 5}, l)
}

func TestPop(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	v := Pop(&l)
	assert.NotNil(t, v)
	assert.Equal(t, 5, *v)
	assert.Equal(t, []int{1, 2, 3, 4}, l)
}

func TestPush(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	a := Push(&l, 6)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, l)
	assert.Equal(t, &[]int{1, 2, 3, 4, 5, 6}, a)
}

func TestUnshift(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	a := Unshift(&l, 0)
	assert.Equal(t, &[]int{0, 1, 2, 3, 4, 5}, a)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, l)
}
