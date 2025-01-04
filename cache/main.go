package main

import(
	"fmt"
)

type Node struct{
  val string
  left *Node
  right *Node
}

type Queue struct{
  head *Node
  tail *Node
  length int
  capacity int
}

type Cache struct{
	Queue Queue
	Hash Hash
}

type Hash map[string]*Node

func NewCache(cap int) Cache{
	return Cache{Queue: NewQueue(cap), Hash: Hash{}}
}

func NewQueue(cap int) Queue{
	head := &Node{}
	tail := &Node{}
	head.right = tail
	tail.left = head

   return Queue{head: head,tail: tail,length: 0,capacity:cap }
}

func (q *Queue) InsertAtEnd(val string){
	if q.length == q.capacity{q.RemoveFromEnd()}
	lastnode_ptr := q.tail.left
	newnode_ptr := &Node{left:lastnode_ptr,right:q.tail,val: val}
	lastnode_ptr.right = newnode_ptr
	q.tail.left = newnode_ptr

}

func (q *Queue) InsertAtStart(val string){
	if q.length == q.capacity{q.RemoveFromEnd()}
	first_node := q.head.right
	newnode := &Node{left: q.head,right: first_node,val: val}
	q.head.right = newnode
	first_node.left = newnode
	q.length++
}

func (q *Queue) RemoveFromEnd(){
	lastnode := q.tail.left
	last_lastnode := lastnode.left
	lastnode.left = nil
	lastnode.right = nil
	last_lastnode.right = q.tail
	q.tail.left = last_lastnode
	q.length--
}

func (q *Queue) BringAtStart(n *Node){
	left_node := n.left
	right_node := n.right
	first_node := q.head.right
	left_node.right = right_node
	right_node.left = left_node
	q.head.right = n
	n.left = q.head
	n.right = first_node
	first_node.left = n
}

func (c *Cache) Access(val string) *Node{
	addr,ok := c.Hash[val]
	if ok{
		c.Queue.BringAtStart(c.Hash[val])
		return addr
	}
	c.Queue.InsertAtStart(val)
	c.Hash[val] = c.Queue.head.right
	return c.Hash[val]
}

func (c *Cache) Print() {
	ptr := c.Queue.head.right
	for ptr != c.Queue.tail{
       fmt.Printf("%v ",ptr.val)
	   ptr = ptr.right
	}
	fmt.Print("\n")
}

func main(){
	first_cache := NewCache(3)
	first_cache.Access("one")
	first_cache.Print()

	first_cache.Access("two")
	first_cache.Print()

	first_cache.Access("three")
	first_cache.Print()

	first_cache.Access("one")
	first_cache.Print()

	first_cache.Access("four")
	first_cache.Print()



}