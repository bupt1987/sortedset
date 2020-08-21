package sortedset

import (
	"math/rand"
	"testing"
)

var s *SortedSet

func init() {
	s = New()
}

func checkOrder(t *testing.T, nodes []*Node, expectedOrder []int64) {
	if len(expectedOrder) != len(nodes) {
		t.Errorf("nodes does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if nodes[i].Key() != expectedOrder[i] {
			t.Errorf("nodes[%d] is %q, but the expected key is %q", i, nodes[i].Key(), expectedOrder[i])
		}

	}
}

func TestCase1(t *testing.T) {
	sortedset := New()

	sortedset.AddOrUpdate(1, 89, "Kelly")
	sortedset.AddOrUpdate(2, 100, "Staley")
	sortedset.AddOrUpdate(3, 100, "Jordon")
	sortedset.AddOrUpdate(4, -321, "Park")
	sortedset.AddOrUpdate(5, 101, "Albert")
	sortedset.AddOrUpdate(6, 99, "Lyman")
	sortedset.AddOrUpdate(7, 99, "Singleton")
	sortedset.AddOrUpdate(8, 70, "Audrey")

	sortedset.AddOrUpdate(5, 99, "ntrnrt")

	sortedset.Remove(2)

	node := sortedset.GetByRank(3, false)
	if node == nil || node.Key() != 1 {
		t.Error("GetByRank() does not return expected value `a`")
	}

	node = sortedset.GetByRank(-3, false)
	if node == nil || node.Key() != 6 {
		t.Error("GetByRank() does not return expected value `f`")
	}

	// get all nodes since the first one to last one
	nodes := sortedset.GetByRankRange(1, -1, false)
	checkOrder(t, nodes, []int64{4, 8, 1, 5, 6, 7, 3})

	// get & remove the 2nd/3rd nodes in reserve order
	nodes = sortedset.GetByRankRange(-2, -3, true)
	checkOrder(t, nodes, []int64{7, 6})

	// get all nodes since the last one to first one
	nodes = sortedset.GetByRankRange(-1, 1, false)
	checkOrder(t, nodes, []int64{3, 5, 1, 8, 4})

}

func TestCase2(t *testing.T) {

	// create a new set
	sortedset := New()

	// fill in new node
	sortedset.AddOrUpdate(1, 89, "Kelly")
	sortedset.AddOrUpdate(2, 100, "Staley")
	sortedset.AddOrUpdate(3, 100, "Jordon")
	sortedset.AddOrUpdate(4, -321, "Park")
	sortedset.AddOrUpdate(5, 101, "Albert")
	sortedset.AddOrUpdate(6, 99, "Lyman")
	sortedset.AddOrUpdate(7, 99, "Singleton")
	sortedset.AddOrUpdate(8, 70, "Audrey")

	// update an existing node
	sortedset.AddOrUpdate(5, 99, "ntrnrt")

	// remove node
	sortedset.Remove(2)

	nodes := sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []int64{4, 8, 1, 5, 6, 7, 3})

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	//t.Logf("%v", nodes)
	checkOrder(t, nodes, []int64{3, 7, 6, 5, 1, 8, 4})

	nodes = sortedset.GetByScoreRange(600, 500, nil)
	checkOrder(t, nodes, []int64{})

	nodes = sortedset.GetByScoreRange(500, 600, nil)
	checkOrder(t, nodes, []int64{})

	rank := sortedset.FindRank(6)
	if rank != 5 {
		t.Error("FindRank() does not return expected value `5`")
	}

	rank = sortedset.FindRank(4)
	if rank != 1 {
		t.Error("FindRank() does not return expected value `1`")
	}

	nodes = sortedset.GetByScoreRange(99, 100, nil)
	checkOrder(t, nodes, []int64{5, 6, 7, 3})

	nodes = sortedset.GetByScoreRange(90, 50, nil)
	checkOrder(t, nodes, []int64{1, 8})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []int64{3})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []int64{7, 6, 5})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []int64{5, 6, 7})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []int64{3})

	nodes = sortedset.GetByScoreRange(50, 100, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []int64{8, 1})

	nodes = sortedset.GetByScoreRange(100, 50, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []int64{3, 7})

	minNode := sortedset.PeekMin()
	if minNode == nil || minNode.Key() != 4 {
		t.Error("PeekMin() does not return expected value `d`")
	}

	minNode = sortedset.PopMin()
	if minNode == nil || minNode.Key() != 4 {
		t.Error("PopMin() does not return expected value `d`")
	}

	nodes = sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []int64{8, 1, 5, 6, 7, 3})

	maxNode := sortedset.PeekMax()
	if maxNode == nil || maxNode.Key() != 3 {
		t.Error("PeekMax() does not return expected value `c`")
	}

	maxNode = sortedset.PopMax()
	if maxNode == nil || maxNode.Key() != 3 {
		t.Error("PopMax() does not return expected value `c`")
	}

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	checkOrder(t, nodes, []int64{7, 6, 5, 1, 8})
}

func BenchmarkSortedSet_Add(b *testing.B) {
	b.StopTimer()
	// data initialization
	scores := make([]float64, b.N)
	IDs := make([]int64, b.N)
	for i := range IDs {
		scores[i] = rand.Float64() + float64(rand.Int31n(99))
		IDs[i] = int64(i) + 100000
	}
	// BCE
	_ = scores[:b.N]
	_ = IDs[:b.N]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.AddOrUpdate(IDs[i], scores[i], nil)
	}
}

func BenchmarkSortedSet_GetRank(b *testing.B) {
	l := s.GetCount()
	for i := 0; i < b.N; i++ {
		s.FindRank(100000 + int64(i)%l)
	}
}

func BenchmarkSortedSet_GetDataByRank(b *testing.B) {
	l := s.GetCount()
	for i := 0; i < b.N; i++ {
		s.GetByKey(int64(i) % l)
	}
}
