package mobilex

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"code.olapie.com/sugar/mobilex/nomobile"
	"code.olapie.com/sugar/testx"
	"github.com/google/uuid"
)

type FileInfo interface {
	ParentID() string
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

	AddSub(sub FileInfo, appendFlag bool)
	RemoveSub(sub FileInfo)

	Move(id, dirID string)
	Find(id string) FileInfo

	SortSubsByModTime(asc bool)
	SortSubsByName(asc bool)
}

var _ FileInfo = (*FileTreeNode)(nil)

var _ nomobile.FileEntry = (*virtualEntry)(nil)

type virtualEntry struct {
	ID          string
	EntryName   string
	SubEntryIDs []string
}

func (v *virtualEntry) GetID() string {
	return v.ID
}

func (v *virtualEntry) Name() string {
	return v.EntryName
}

func (v *virtualEntry) IsDir() bool {
	return true
}

func (v *virtualEntry) Size() int64 {
	return 0
}

func (v *virtualEntry) ModTime() int64 {
	return 0
}

func (v *virtualEntry) MIMEType() string {
	return ""
}

func (v *virtualEntry) SubIDs() []string {
	return v.SubEntryIDs
}

type FileTreeNode struct {
	entry  nomobile.FileEntry
	parent *FileTreeNode
	dirs   []*FileTreeNode
	files  []*FileTreeNode
}

func (f *FileTreeNode) Entry() nomobile.FileEntry {
	return f.entry
}

func (f *FileTreeNode) ParentID() string {
	if f.parent == nil {
		return ""
	}
	return f.parent.GetID()
}

func (f *FileTreeNode) GetID() string {
	return f.entry.GetID()
}

func (f FileTreeNode) Name() string {
	return f.entry.Name()
}

func (f FileTreeNode) IsDir() bool {
	return f.entry.IsDir()
}

func (f FileTreeNode) Size() int64 {
	return f.entry.Size()
}

func (f FileTreeNode) ModTime() int64 {
	return f.entry.ModTime()
}

func (f FileTreeNode) MIMEType() string {
	return f.entry.MIMEType()
}

func (f *FileTreeNode) SubsCount() int {
	return len(f.dirs) + len(f.files)
}

func (f *FileTreeNode) GetSub(i int) FileInfo {
	if i < len(f.dirs) {
		return f.dirs[i]
	}
	return f.files[i-len(f.dirs)]
}

func (f *FileTreeNode) DirsCount() int {
	return len(f.dirs)
}

func (f *FileTreeNode) GetDir(i int) FileInfo {
	return f.dirs[i]
}

func (f *FileTreeNode) FilesCount() int {
	return len(f.files)
}

func (f *FileTreeNode) GetFile(i int) FileInfo {
	return f.files[i]
}

func (f *FileTreeNode) AddSub(sub FileInfo, appendFlag bool) {
	node := sub.(*FileTreeNode)
	if node.parent != nil {
		panic(fmt.Sprintf("%s is in dir %s", node.GetID(), node.parent.GetID()))
	}
	if node.IsDir() {
		if appendFlag {
			f.dirs = append(f.dirs, node)
		} else {
			f.dirs = append([]*FileTreeNode{node}, f.dirs...)
		}
	} else {
		if appendFlag {
			f.files = append(f.files, node)
		} else {
			f.files = append([]*FileTreeNode{node}, f.files...)
		}
	}
	node.parent = f
}

func (f *FileTreeNode) RemoveSub(sub FileInfo) {
	node := sub.(*FileTreeNode)
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
func (f *FileTreeNode) Move(id, dirID string) {
	fi := f.Find(id)
	f.RemoveSub(fi)
	f.Find(dirID).(*FileTreeNode).AddSub(fi, true)
}

// Find searches a descendant node
func (f *FileTreeNode) Find(id string) FileInfo {
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

func (f *FileTreeNode) SortSubsByModTime(asc bool) {
	for _, dir := range f.dirs {
		dir.SortSubsByModTime(asc)
	}
	sort.Slice(f.dirs, func(i, j int) bool {
		fi, fj := f.dirs[i], f.dirs[j]
		if fi.ModTime() == fj.ModTime() {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		}
		return asc == (fi.ModTime() < fj.ModTime())
	})

	sort.Slice(f.files, func(i, j int) bool {
		fi, fj := f.files[i], f.files[j]
		if fi.ModTime() == fj.ModTime() {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		}
		return asc == (fi.ModTime() < fj.ModTime())
	})
}

func (f *FileTreeNode) SortSubsByName(asc bool) {
	for _, dir := range f.dirs {
		dir.SortSubsByName(asc)
	}
	sort.Slice(f.dirs, func(i, j int) bool {
		fi, fj := f.dirs[i], f.dirs[j]
		if fi.Name() == fj.Name() {
			return asc == (fi.ModTime() < fj.ModTime())
		}
		return asc == (fi.Name() == fj.Name())
	})

	sort.Slice(f.files, func(i, j int) bool {
		fi, fj := f.files[i], f.files[j]
		if fi.Name() == fj.Name() {
			return asc == (fi.ModTime() < fj.ModTime())
		}
		return asc == (fi.Name() == fj.Name())
	})
}

func NewVirtualDir(id, name string) FileInfo {
	return &FileTreeNode{
		entry: &virtualEntry{
			ID:        id,
			EntryName: name,
		},
	}
}

func FileInfoFromEntry(entry nomobile.FileEntry) FileInfo {
	return &FileTreeNode{
		entry: entry,
	}
}

func BuildFileTree(entries []nomobile.FileEntry) FileInfo {
	root := NewVirtualDir("", "").(*FileTreeNode)
	if len(entries) == 0 {
		return root
	}

	idToEntry := make(map[string]nomobile.FileEntry)
	for _, f := range entries {
		idToEntry[f.GetID()] = f
	}

	idToNode := make(map[string]*FileTreeNode)
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

	root.SortSubsByModTime(false)
	return root
}

func buildFileTreeNode(parent *FileTreeNode, id string, idToEntry map[string]nomobile.FileEntry, result map[string]*FileTreeNode) {
	node := result[id]
	if node != nil {
		if parent != nil {
			parent.AddSub(node, true)
		}
		return
	}

	entry, exists := idToEntry[id]
	if !exists {
		fmt.Println("Warn: no entry for file id", id)
		return
	}
	delete(idToEntry, id)

	node = &FileTreeNode{
		entry: entry,
	}
	if parent != nil {
		parent.AddSub(node, true)
	}
	result[node.GetID()] = node
	if !node.IsDir() {
		return
	}

	for _, subID := range entry.SubIDs() {
		buildFileTreeNode(node, subID, idToEntry, result)
	}
}

var _ nomobile.FileEntry = (*mockFileEntry)(nil)

type mockFileEntry struct {
	id       string
	name     string
	isDir    bool
	size     int64
	modTime  int64
	mimeType string
	subIDs   []string
}

func (m *mockFileEntry) GetID() string {
	return m.id
}

func (m *mockFileEntry) Name() string {
	return m.name
}

func (m *mockFileEntry) IsDir() bool {
	return m.isDir
}

func (m *mockFileEntry) Size() int64 {
	return m.size
}

func (m *mockFileEntry) ModTime() int64 {
	return m.modTime
}

func (m *mockFileEntry) MIMEType() string {
	return m.mimeType
}

func (m *mockFileEntry) SubIDs() []string {
	return m.subIDs
}

func NewMockFileInfo(isDir bool) FileInfo {
	if isDir {
		return &FileTreeNode{
			entry: &mockFileEntry{
				id:      uuid.NewString(),
				name:    "dir" + testx.RandomString(10),
				isDir:   true,
				modTime: time.Now().Unix(),
				subIDs:  []string{uuid.NewString(), uuid.NewString()},
			},
		}
	}
	return &FileTreeNode{
		entry: &mockFileEntry{
			id:      uuid.NewString(),
			name:    "file" + testx.RandomString(10),
			modTime: time.Now().Unix(),
		},
	}
}
