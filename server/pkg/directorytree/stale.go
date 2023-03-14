package directorytree

type Stale struct {
	dirs  []string
	files []string
}

func FindStales(db *DirectoryNode) *Stale {
	result := Stale{}
	findStale(db, &result)
	return &result
}

func findStale(node *DirectoryNode, result *Stale) {
	if node.Excluded {
		addToStale(node, result)
	} else {
		for _, child := range node.Children {
			findStale(child, result)
		}
	}
}

func addToStale(node *DirectoryNode, result *Stale) {
	result.dirs = append(result.dirs, node.getPath())

	for _, file := range node.Files {
		result.files = append(result.files, file.getPath())
	}

	for _, child := range node.Children {
		addToStale(child, result)
	}
}
