package main

import (
    "errors"
    "fmt"
    "bufio"
    "os"
    "log"
    "strings"
    "strconv"
    "math"
    "container/heap"
    "time"
)

/*
 * Min Heap
 */
type DistanceNode struct {
    to *Node
    index int
}
type PriorityQueue []*DistanceNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].to.dist < pq[j].to.dist
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*DistanceNode)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    item.index = -1 // for safety
    *pq = old[0 : n-1]

    return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(node *DistanceNode) {
    heap.Fix(pq, node.index)
}


/*
 * Vertex & Edges
 */

type Link struct {
    begin *Node
    end *Node
    description string
}

type Node struct {
    id string
    longitude int32
    latitude int32
    state string
    description string
    dist float64
    adjacency []*Node
    parent *Node
}

func (n Node) distance(n2 *Node) (float64) {
    latitude_diff := n.latitude - n2.latitude
    longitude_diff := n.longitude - n2.longitude

    lat := math.Pow(float64(latitude_diff), 2)
    lng := math.Pow(float64(longitude_diff), 2)

    return math.Pow(lat + lng, 0.5)
}


func findByCityAndState(nodes []*Node, city, state string) (*Node, error) {
    for _, node := range nodes {
        if node.state == state && strings.Contains(node.description, city) {
            return node, nil
        }
    }

    return nil, errors.New("Couldn't find the node")
}

// Load Nodes from file
func loadNodes () ([]*Node, map[string]*Node) {
    nodes := make([]*Node, 0, 100000)
    mapNode := make(map[string]*Node)
    file, err := os.Open("data/nhpn.nod")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        // Line parsing
        nodeId := strings.TrimSpace(line[23:33])
        lng := strings.TrimSpace(line[33:43])
        lat := strings.TrimSpace(line[45:53])
        state := strings.TrimSpace(line[53:55])
        description := strings.TrimSpace(line[55:88])

        // String -> Int
        longitude, _ := strconv.ParseInt(lng, 10, 32)
        latitude, _ := strconv.ParseInt(lat, 10, 32)

        node := Node{
            id : nodeId,
            longitude: int32(longitude),
            latitude: int32(latitude),
            state: state,
            description: description,
            adjacency: make([]*Node, 0),
        }

        nodes = append(nodes, &node)
        mapNode[nodeId] = &node;
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return nodes, mapNode;
}

// Load links from file
func loadLinks(mapNode map[string]*Node) ([]Link) {
    links := make([]Link, 0, 126000)

    file, err := os.Open("data/nhpn.lnk")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        begin := strings.TrimSpace(line[33:43])
        end := strings.TrimSpace(line[43:53])
        description := strings.TrimSpace(line[53:88])

        link := Link{begin: mapNode[begin], end: mapNode[end], description: description}
        links = append(links, link)
    }

    return links
}

func reverse_path(nodes []*Node) ([]*Node) {
    for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
        tmp := nodes[i]
        nodes[i] = nodes[j]
        nodes[j] = tmp
    }

    return nodes
}

func dijkstra (nodes []*Node, source, destination *Node) ([]*Node) {
    // Store DistanceNode in a map to be able
    // to retrieve the index of a given in the queue
    distNodes := make(map[string]*DistanceNode)

    // Set all nodes to be at '+infinity' distance from source
    for _, node := range nodes {
        node.dist = math.MaxInt32
        node.parent = nil
    }

    // Push source node to min heap
    source.dist = 0
    pq := make(PriorityQueue, 1)
    pq[0] = &DistanceNode{
        to: source,
        index: 0,
    }
    heap.Init(&pq)

    for len(pq) != 0 {
        node := heap.Pop(&pq).(*DistanceNode).to

        if node.id == destination.id { break }

        for _, next_node := range node.adjacency {
            next_dist := node.dist + node.distance(next_node)

            // Push it to min heap if not already in
            if next_node.dist == math.MaxInt32 {
                next_node.dist = next_dist
                next_node.parent = node

                tmp := &DistanceNode{to: next_node}

                distNodes[next_node.id] = tmp
                heap.Push(&pq, tmp)
            } else if next_node.dist > next_dist {
                // Relaxation
                next_node.dist = next_dist

                // Decrease_key
                pq.update(distNodes[next_node.id])

                next_node.parent = node
            }
        }
    }

    // Create path from destination to source
    path := make([]*Node, 0, 100)

    for destination != nil {
        path = append(path, destination)
        destination = destination.parent
    }

    return reverse_path(path)

}

func FilesLoader () ([]*Node) {
    nodes, mapNode := loadNodes()

    // Populate nodes adjacency list
    for _, link := range loadLinks(mapNode) {
        link.begin.adjacency = append(link.begin.adjacency, link.end)
        link.end.adjacency = append(link.end.adjacency, link.begin)
    }

    return nodes
}



func main() {
    nodes := FilesLoader()

    start, err := findByCityAndState(nodes, "PASADENA", "CA")
    if err != nil { log.Fatal(err); return }

    end, err := findByCityAndState(nodes, "CAMBRIDGE", "MA")
    if err != nil { log.Fatal(err); return }

    startTime := time.Now()

    dijkstra(nodes, start, end)

    elapsed := time.Since(startTime)
    fmt.Println(elapsed)
}
