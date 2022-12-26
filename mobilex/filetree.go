package mobilex

import (
	"fmt"
	"sort"
	"strings"
)

type FileEntry interface {
	GetID() string
	Name() string
	IsDir() bool
	Size() int64
	ModTime() int64
	MIMEType() string
	SubIDs() []string
}

type FileInfo interface {
	GetID() string
	Name() string
	IsDir() bool
	Size() int64
	ModTime() int64
	MIMEType() string

	SubsCount() int
	GetSub(i int) FileInfo

	DirsCount() int
	GetDir(i int) FileInfo

	FilesCount() int
	GetFile(i int) FileInfo
}

var _ FileInfo = (*fileTreeNode)(nil)

type fileTreeNode struct {
	info   FileEntry
	parent *fileTreeNode
	dirs   []*fileTreeNode
	files  []*fileTreeNode
}

func (f *fileTreeNode) GetID() string {
	if f.info == nil {
		return ""
	}
	return f.info.GetID()
}

func (f fileTreeNode) Name() string {
	if f.info == nil {
		return ""
	}
	return f.info.Name()
}

func (f fileTreeNode) IsDir() bool {
	if f.info == nil {
		return true
	}
	return f.info.IsDir()
}

func (f fileTreeNode) Size() int64 {
	if f.info == nil {
		return 0
	}
	return f.info.Size()
}

func (f fileTreeNode) ModTime() int64 {
	if f.info == nil {
		return 0
	}
	return f.info.ModTime()
}

func (f fileTreeNode) MIMEType() string {
	if f.info == nil {
		return ""
	}
	return f.info.MIMEType()
}

func (f *fileTreeNode) SubsCount() int {
	return len(f.dirs) + len(f.files)
}

func (f *fileTreeNode) GetSub(i int) FileInfo {
	if i < len(f.dirs) {
		return f.dirs[i]
	}
	return f.files[i-len(f.dirs)]
}

func (f *fileTreeNode) DirsCount() int {
	return len(f.dirs)
}

func (f *fileTreeNode) GetDir(i int) FileInfo {
	return f.dirs[i]
}

func (f *fileTreeNode) FilesCount() int {
	return len(f.files)
}

func (f *fileTreeNode) GetFile(i int) FileInfo {
	return f.files[i]
}

func (f *fileTreeNode) AddSub(sub FileInfo) {
	node := sub.(*fileTreeNode)
	if node.parent != nil {
		panic(fmt.Sprintf("%s is in dir %s", node.GetID(), node.parent.GetID()))
	}
	if node.IsDir() {
		f.dirs = append(f.dirs, node)
	} else {
		f.files = append(f.files, node)
	}
	node.parent = f
}

func (f *fileTreeNode) RemoveSub(sub FileInfo) {
	node := sub.(*fileTreeNode)
	node.parent = nil
	for i, v := range f.dirs {
		if v.GetID() == node.GetID() {
			f.dirs = append(f.dirs[:i], f.dirs[i+1:]...)
			return
		}
	}

	for i, v := range f.files {
		if v.GetID() == node.GetID() {
			f.files = append(f.files[:i], f.files[i+1:]...)
			return
		}
	}
}

// Move file id to directory
func (f *fileTreeNode) Move(id, dirID string) {
	fi := f.Find(id)
	f.RemoveSub(fi)
	f.Find(dirID).(*fileTreeNode).AddSub(fi)
}

// Find searches a descendant node
func (f *fileTreeNode) Find(id string) FileInfo {
	if f.GetID() == id {
		return f
	}

	for _, fi := range f.files {
		if fi.GetID() == id {
			return fi
		}
	}

	for _, dir := range f.dirs {
		if fi := dir.Find(id); fi != nil {
			return fi
		}
	}

	return nil
}

func (f *fileTreeNode) SortSubs() {
	for _, dir := range f.dirs {
		dir.SortSubs()
	}
	sort.Slice(f.dirs, func(i, j int) bool {
		fi, fj := f.dirs[i], f.dirs[j]
		if fi.ModTime() == fj.ModTime() {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		}
		return fi.ModTime() < fj.ModTime()
	})

	sort.Slice(f.files, func(i, j int) bool {
		fi, fj := f.files[i], f.files[j]
		if fi.ModTime() == fj.ModTime() {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		}
		return fi.ModTime() < fj.ModTime()
	})
}

func BuildFileTree(entries []FileEntry) *fileTreeNode {
	root := new(fileTreeNode)
	if len(entries) == 0 {
		return root
	}

	idToEntry := make(map[string]FileEntry)
	for _, f := range entries {
		idToEntry[f.GetID()] = f
	}

	idToNode := make(map[string]*fileTreeNode)
	for _, e := range entries {
		buildFileTreeNode(nil, e.GetID(), idToEntry, idToNode)
	}

	for _, node := range idToNode {
		if node.parent == nil {
			if node.IsDir() {
				root.dirs = append(root.dirs, node)
			} else {
				root.files = append(root.files, node)
			}
		}
	}

	root.SortSubs()
	return root
}

func buildFileTreeNode(parent *fileTreeNode, id string, idToEntry map[string]FileEntry, result map[string]*fileTreeNode) {
	node := result[id]
	if node != nil {
		if parent != nil {
			parent.AddSub(node)
		}
		return
	}

	entry, exists := idToEntry[id]
	if !exists {
		fmt.Println("Warn: no entry for file id", id)
		return
	}
	delete(idToEntry, id)

	node = &fileTreeNode{
		info: entry,
	}
	if parent != nil {
		parent.AddSub(node)
	}
	result[node.GetID()] = node
	if !node.IsDir() {
		return
	}

	for _, subID := range entry.SubIDs() {
		buildFileTreeNode(node, subID, idToEntry, result)
	}
}
