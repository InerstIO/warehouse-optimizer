package warehouse

import (
	"sort"
)

// ByI is the type of a "less" function that defines the ordering of its Planet arguments.
type ByI func(i1, i2 *Item, m map[int]Product) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by ByI) Sort(items Order, m map[int]Product) {
	is := &itemSorter{
		items: items,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
		m:      m,
	}
	sort.Sort(is)
}

// itemSorter joins a By function and a slice of Items to be sorted.
type itemSorter struct {
	items Order
	by      func(i1, i2 *Item, m map[int]Product) bool // Closure used in the Less method.
	m      map[int]Product
}

// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(&s.items[i], &s.items[j], s.m)
}

// ByItemWeightReverse is a closure that order the Items by reverse weight
func ByItemWeightReverse(i1, i2 *Item, m map[int]Product) bool {
	return m[i1.ProdID].w > m[i2.ProdID].w
}

// SplitOrder splits the order that has total weight larger than max
func SplitOrder(order Order, m map[int]Product, max float64) []Order {
	if OrderWeight(order, m) <= max {
		return []Order{order}
	}
	reOrders := make([]Order, 1)
	ordersWeight := []float64{0.0}
	ByI(ByItemWeightReverse).Sort(order, m)
	var fit bool
	for _, i := range order {
		fit = false
		for j := range reOrders {
			iw := m[i.ProdID].w
			if ordersWeight[j]+iw <= max || len(reOrders[j]) == 0 {
				reOrders[j] = append(reOrders[j], i)
				ordersWeight[j] += iw
				fit = true
				break
			}
		}
		if !fit {
			var newOrder Order
			newOrder = append(newOrder, i)
			reOrders = append(reOrders, newOrder)
			ordersWeight = append(ordersWeight, m[i.ProdID].w)
		}
	}
	return reOrders
}