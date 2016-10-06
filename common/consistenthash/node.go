package consistenthash

type Node struct {
	Id     string
	HashId uint32
}

func NewNode(id string) *Node {
	return &Node{
		Id:     id,
		HashId: hashId(id),
	}
}

type Nodes []*Node

// The following methods are the implementation of interface below:
// type Interface interface {
//         // Len is the number of elements in the collection.
//         Len() int
//         // Less reports whether the element with
//         // index i should sort before the element with index j.
//         Less(i, j int) bool
//         // Swap swaps the elements with indexes i and j.
//         Swap(i, j int)
// }
// This is required to run sort.Sort(r.Nodes) which expects the parameter to
// implement the interface.

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].HashId < n[j].HashId }
