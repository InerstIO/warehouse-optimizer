package warehouse

import (
	"sort"
)

const (
	max_item = 15
)

// By is the type of a "less" function that defines the ordering of its Order arguments.
type By func(o1, o2 *Order, m map[int]Product) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(orders []Order, m map[int]Product) {
	ps := &orderSorter{
		orders: orders,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
		m:      m,
	}
	sort.Sort(ps)
}

// orderSorter joins a By function and a slice of Orders to be sorted.
type orderSorter struct {
	orders []Order
	by     func(o1, o2 *Order, m map[int]Product) bool // Closure used in the Less method.
	m      map[int]Product
}

// Len is part of sort.Interface.
func (s *orderSorter) Len() int {
	return len(s.orders)
}

// Swap is part of sort.Interface.
func (s *orderSorter) Swap(i, j int) {
	s.orders[i], s.orders[j] = s.orders[j], s.orders[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *orderSorter) Less(i, j int) bool {
	return s.by(&s.orders[i], &s.orders[j], s.m)
}

// ByWeightReverse is a closure that order the Orders by reverse weight
func ByWeightReverse(o1, o2 *Order, m map[int]Product) bool {
	return OrderWeight(*o1, m) > OrderWeight(*o2, m)
}

// MergeOrders returns the reconbined order IDs that has total weight not larger than max
func MergeOrders(orders []Order, m map[int]Product, max float64) []Order {
	reOrders := make([]Order, 1)
	ordersWeight := []float64{0.0}
	By(ByWeightReverse).Sort(orders, m)
	var fit bool
	for _, o := range orders {
		fit = false
		for j := range reOrders {
			ow := OrderWeight(o, m)
			if (ordersWeight[j]+ow <= max && len(reOrders[j]) + len(o) <= max_item ) || len(reOrders[j]) == 0 {
				reOrders[j] = append(reOrders[j], o...)
				ordersWeight[j] += ow
				fit = true
				break
			}
		}
		if !fit {
			var newOrder Order
			newOrder = append(newOrder, o...)
			reOrders = append(reOrders, newOrder)
			ordersWeight = append(ordersWeight, OrderWeight(o, m))
		}
	}
	return reOrders
}
