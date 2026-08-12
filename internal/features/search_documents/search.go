package search_documents

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// MatchResult は検索のヒット情報を格納する構造体です。
type MatchResult struct {
	Path       string `json:"path"`
	FileName   string `json:"fileName"`
	LineNumber int    `json:"lineNumber,omitempty"`
	Snippet    string `json:"snippet,omitempty"`
	IsTitle    bool   `json:"isTitle,omitempty"`
}

// SearchDocuments は指定ディレクトリ配下の Markdown からキーワードに一致する箇所を探索します。
func SearchDocuments(rootDir, query string) ([]*MatchResult, error) {
	if strings.TrimSpace(query) == "" {
		return []*MatchResult{}, nil
	}

	queryLower := strings.ToLower(query)
	var results []*MatchResult

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		name := info.Name()
		// 隠しフォルダ・ドットファイルはスキップ
		if info.IsDir() {
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}
		relPathSlash := filepath.ToSlash(relPath)

		// 1. ファイル名へのマッチングチェック
		if strings.Contains(strings.ToLower(name), queryLower) {
			results = append(results, &MatchResult{
				Path:     relPathSlash,
				FileName: name,
				IsTitle:  true,
			})
		}

		// 2. 本文へのマッチングチェック
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			lineText := scanner.Text()
			if strings.Contains(strings.ToLower(lineText), queryLower) {
				snippet := strings.TrimSpace(lineText)
				if len(snippet) > 120 {
					snippet = snippet[:120] + "..."
				}
				results = append(results, &MatchResult{
					Path:       relPathSlash,
					FileName:   name,
					LineNumber: lineNum,
					Snippet:    snippet,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
