package graph

import "fmt"

// Validate checks structural integrity of a workflow definition: no
// duplicate node ids, exactly one start node, every referenced successor
// id exists, and at least one stop node is reachable from start. It does
// not evaluate expressions (e.g. builtin.if/builtin.loop conditions) since
// Validate has no access to a runtime variable set; that is the caller's
// job, typically via the expr package at definition-registration time.
func Validate(d *Definition) error {
	if d.ID == "" {
		return fmt.Errorf("definition id is required")
	}
	if len(d.Nodes) == 0 {
		return fmt.Errorf("definition %q has no nodes", d.ID)
	}

	byID := make(map[string]Node, len(d.Nodes))
	for _, n := range d.Nodes {
		if n.ID == "" {
			return fmt.Errorf("definition %q has a node with an empty id", d.ID)
		}
		if _, dup := byID[n.ID]; dup {
			return fmt.Errorf("definition %q has duplicate node id %q", d.ID, n.ID)
		}
		byID[n.ID] = n
	}

	var startIDs, stopIDs []string
	for _, n := range d.Nodes {
		switch n.Type {
		case TypeStart:
			startIDs = append(startIDs, n.ID)
		case TypeStop, TypeTerminate:
			stopIDs = append(stopIDs, n.ID)
		}
		for _, succ := range n.Successors() {
			if _, ok := byID[succ]; !ok {
				return fmt.Errorf("definition %q: node %q references unknown successor %q", d.ID, n.ID, succ)
			}
		}
	}

	if len(startIDs) != 1 {
		return fmt.Errorf("definition %q must have exactly one %s node, found %d", d.ID, TypeStart, len(startIDs))
	}
	if len(stopIDs) == 0 {
		return fmt.Errorf("definition %q must have at least one %s node", d.ID, TypeStop)
	}

	if !reachableFrom(byID, startIDs[0], stopIDs) {
		return fmt.Errorf("definition %q: no %s node is reachable from %s", d.ID, TypeStop, TypeStart)
	}

	return nil
}

// reachableFrom performs a breadth-first search over the graph starting at
// startID and reports whether any node id in targets is visited. Cycles
// (e.g. from builtin.loop back-edges) are handled via the visited set.
func reachableFrom(byID map[string]Node, startID string, targets []string) bool {
	targetSet := make(map[string]bool, len(targets))
	for _, t := range targets {
		targetSet[t] = true
	}

	visited := map[string]bool{startID: true}
	queue := []string{startID}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if targetSet[cur] {
			return true
		}
		for _, succ := range byID[cur].Successors() {
			if !visited[succ] {
				visited[succ] = true
				queue = append(queue, succ)
			}
		}
	}
	return false
}
