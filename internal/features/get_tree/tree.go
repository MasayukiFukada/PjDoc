package get_tree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Node はディレクトリツリーの各要素を表す構造体です。
type Node struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	IsDir    bool    `json:"isDir"`
	Children []*Node `json:"children,omitempty"`
}

// BuildTree は指定されたルートディレクトリから Markdown ファイルを含むツリーを構築します。
func BuildTree(rootDir string) (*Node, error) {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	rootNode := &Node{
		Name:  filepath.Base(absRoot),
		Path:  "",
		IsDir: true,
	}

	err = buildSubTree(absRoot, "", rootNode)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func buildSubTree(absBase, relPath string, parent *Node) error {
	currentDir := filepath.Join(absBase, relPath)
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return err
	}

	var children []*Node

	for _, entry := range entries {
		name := entry.Name()

		// .git などの隠しディレクトリや指定除外フォルダをスキップします
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
			continue
		}

		childRelPath := filepath.Join(relPath, name)
		if entry.IsDir() {
			dirNode := &Node{
				Name:  name,
				Path:  filepath.ToSlash(childRelPath),
				IsDir: true,
			}
			err := buildSubTree(absBase, childRelPath, dirNode)
			if err != nil {
				return err
			}
			// 子要素に Markdown ファイルが含まれるディレクトリのみ保持します
			if len(dirNode.Children) > 0 {
				children = append(children, dirNode)
			}
		} else {
			if strings.HasSuffix(strings.ToLower(name), ".md") {
				children = append(children, &Node{
					Name:  name,
					Path:  filepath.ToSlash(childRelPath),
					IsDir: false,
				})
			}
		}
	}

	// ディレクトリ優先・アルファベット順でソートします
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Name < children[j].Name
	})

	parent.Children = children
	return nil
}
