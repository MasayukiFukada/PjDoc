package render_document

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidPath = errors.New("invalid or illegal file path")
	ErrNotFound    = errors.New("file not found")
)

// Document は読み込んだドキュメントの情報を格納する構造体です。
type Document struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// ReadDocument は指定された相対パスの Markdown ファイルを安全に読み込みます。
func ReadDocument(rootDir, relPath string) (*Document, error) {
	cleanRelPath := filepath.Clean(relPath)

	// ディレクトリトラバーサル攻撃を検証・防ぎます
	if strings.HasPrefix(cleanRelPath, "..") || filepath.IsAbs(cleanRelPath) {
		return nil, ErrInvalidPath
	}

	fullPath := filepath.Join(rootDir, cleanRelPath)
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, err
	}

	// ルートディレクトリのシンボリックリンクを解決
	evalRootDir, err := filepath.EvalSymlinks(absRootDir)
	if err != nil {
		return nil, err
	}

	// 対象ファイルのシンボリックリンクを解決して存在確認
	evalFullPath, err := filepath.EvalSymlinks(absFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 解決後の実パスがルートディレクトリ配下に収まっているかを検証 (シンボリックリンク悪用防止)
	relToRoot, err := filepath.Rel(evalRootDir, evalFullPath)
	if err != nil || strings.HasPrefix(relToRoot, "..") || filepath.IsAbs(relToRoot) {
		return nil, ErrInvalidPath
	}

	info, err := os.Stat(evalFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if info.IsDir() {
		return nil, ErrInvalidPath
	}

	data, err := os.ReadFile(evalFullPath)
	if err != nil {
		return nil, err
	}

	return &Document{
		Path:    filepath.ToSlash(cleanRelPath),
		Content: string(data),
	}, nil
}
