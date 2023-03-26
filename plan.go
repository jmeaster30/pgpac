package main

type Plan struct {
	idToObjectMap      map[int64]SchemaObject
	parallelObjectList [][]int64
}

func (p *Plan) BuildPlan(objs []SchemaObject) {

}
