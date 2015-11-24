package main

import (
    "fmt"
)

type Vertex struct {
    name string
    pi *Vertex
    distance int
    status bool
}

type Graph struct {
    Vertices map[string][]*Vertex
}


func (g Graph) addVertex(vtx Vertex) bool {
    if _, ok := g.Vertices[vtx.name]; !ok {
        g.Vertices[vtx.name] = []*Vertex{}
        return true
    }
    return false
}

func (g Graph) addEdge(vtx1, vtx2 *Vertex) {
    if g.Vertices[vtx1.name] != nil && g.Vertices[vtx2.name] != nil {
        g.Vertices[vtx1.name] = append(g.Vertices[vtx1.name], vtx2)
        g.Vertices[vtx2.name] = append(g.Vertices[vtx2.name], vtx1)
    } else {
        fmt.Println("Ca va peter !!!")
    }
}

func (g Graph) showGraph() {
    for k := range g.Vertices {
        fmt.Print("Vertex : ", k, " Edges : ")
        for i := range g.Vertices[k] {
            fmt.Print(g.Vertices[k][i].name, " ")
        }
        fmt.Println("")
    }
}

func shortestPath(s *Vertex) {
    fmt.Print(s.name)
    if s.pi != nil {
        fmt.Print(" --> ")
        shortestPath(s.pi)
    }
}


// Given a graph and a source return the shorten path
func BFS(g *Graph, s *Vertex) {
    neighbour := make([]*Vertex, 0)
    s.status = true

    for _, v := range g.Vertices[s.name] {
        if v.status == false {
            fmt.Println("Adding")
            neighbour = append(neighbour, v)
        }
    }

    for _, v := range neighbour {
        v.distance = s.distance + 1
        v.pi = s
        v.status = true
        BFS(g, v)
    }
}

func main() {
    var g Graph

    a := Vertex{name: "A"}
    b := Vertex{name: "B"}
    c := Vertex{name: "C"}
    d := Vertex{name: "D"}
    e := Vertex{name: "E"}

    g.Vertices = make(map[string][]*Vertex)

    g.addVertex(a)
    g.addVertex(b)
    g.addVertex(c)
    g.addVertex(d)
    g.addVertex(e)
    g.addEdge(&a, &b)
    g.addEdge(&a, &c)
    g.addEdge(&c, &e)
    g.addEdge(&e, &d)
    g.addEdge(&b, &d)
    g.showGraph()

    BFS(&g, &a)

    g.showGraph()

    fmt.Println(a)
    fmt.Println(b)
    fmt.Println(c)
    fmt.Println(e)
    fmt.Println(d)

    shortestPath(&d)
}
