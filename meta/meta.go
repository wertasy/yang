package meta

import (
	"fmt"
	"strings"

	"github.com/c2stack/c2g/c2"
)

///////////////////
// Interfaces
//////////////////

// Examples: Just about everything
type Identifiable interface {
	GetIdent() string
}

// Examples: Most everything except items like ChoiceCase, RpcInput
type Describable interface {
	Identifiable
	SetDescription(string)
	GetDescription() string
	SetReference(string)
	GetReference() string
}

// Examples: Things that have more than one.
type Meta interface {
	Identifiable
	GetParent() MetaList
	SetParent(MetaList)
	GetSibling() Meta
	SetSibling(Meta)
}

// Examples: Module, Container but not Leaf or LeafList
type MetaList interface {
	Meta
	AddMeta(Meta) error
	GetFirstMeta() Meta
	ReplaceMeta(oldChild Meta, newChild Meta) error
	Clear()
}

type DataDef interface {
	Meta
	NextDataDef() DataDef
}

type HasGroupings interface {
	MetaList
	GetGroupings() MetaList
}

type HasTypedefs interface {
	MetaList
	GetTypedefs() MetaList
}

type HasDataType interface {
	Meta
	GetDataType() *DataType
	SetDataType(dataType *DataType)
}

type MetaProxy interface {
	ResolveProxy() Iterator
}

///////////////////////
// Base structs
///////////////////////

// MetaList implementation helper(s)
type ListBase struct {
	// Parent? - it's normally in MetaBase
	FirstMeta Meta
	LastMeta  Meta
}

func (y *ListBase) Clear() {
	y.FirstMeta = nil
	y.LastMeta = nil
}
func (y *ListBase) linkMeta(impl MetaList, meta Meta) error {
	meta.SetParent(impl)
	if y.LastMeta != nil {
		y.LastMeta.SetSibling(meta)
	}
	y.LastMeta = meta
	if y.FirstMeta == nil {
		y.FirstMeta = meta
	}
	return nil
}
func (y *ListBase) swapMeta(oldChild Meta, newChild Meta) error {
	previousSibling := y.FirstMeta
	for previousSibling != nil && previousSibling.GetSibling() != oldChild {
		previousSibling = previousSibling.GetSibling()
	}
	if previousSibling == nil {
		return &schemaError{"child not found"}
	}
	previousSibling.SetSibling(newChild)
	newChild.SetSibling(oldChild.GetSibling())
	newChild.SetParent(oldChild.GetParent())
	if y.FirstMeta == oldChild {
		y.FirstMeta = newChild
	}
	if y.LastMeta == oldChild {
		y.LastMeta = newChild
	}
	return nil
}

// Meta implementation helpers
type MetaBase struct {
	Parent  MetaList
	Sibling Meta
}

// Meta and MetaList combination helpers
type MetaContainer struct {
	Ident string
	MetaBase
	ListBase
}

// Meta
func (y *MetaContainer) GetIdent() string {
	return y.Ident
}

// Meta
func (y *MetaContainer) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *MetaContainer) GetParent() MetaList {
	return y.Parent
}
func (y *MetaContainer) GetSibling() Meta {
	return y.Sibling
}
func (y *MetaContainer) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *MetaContainer) AddMeta(meta Meta) error {
	return y.linkMeta(y, meta)
}
func (y *MetaContainer) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *MetaContainer) Clear() {
	y.Clear()
}
func (y *MetaContainer) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

type metaError struct {
	Msg string
}

func (e *metaError) Error() string {
	return e.Msg
}

////////////////////////
// Implementations
/////////////////////////

type Module struct {
	Ident        string
	Description  string
	Namespace    string
	Contact      string
	Organization string
	Reference    string
	Revision     *Revision
	Prefix       string
	MetaBase
	Defs      MetaContainer
	Groupings MetaContainer
	Typedefs  MetaContainer
	Imports   map[string]*Import
	Includes  []*Include
}

// Identifiable
func (y *Module) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Module) GetDescription() string {
	return y.Description
}
func (y *Module) SetDescription(d string) {
	y.Description = d
}
func (y *Module) GetReference() string {
	return y.Reference
}
func (y *Module) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Module) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Module) GetParent() MetaList {
	return y.Parent
}
func (y *Module) GetSibling() Meta {
	return y.Sibling
}
func (y *Module) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Module) AddMeta(meta Meta) error {
	switch x := meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, x)
	case *Typedef:
		y.Typedefs.SetParent(y)
		return y.Typedefs.linkMeta(y, x)
	default:
		y.Defs.SetParent(y)
		return y.Defs.linkMeta(y, x)
	}
}

// technically not true, it's the MetaContainers, but we'll see how this pans out
func (y *Module) GetFirstMeta() Meta {
	return y.Defs.GetFirstMeta()
}
func (y *Module) Clear() {
	y.Clear()
}
func (y *Module) DataDefs() MetaList {
	return &y.Defs
}
func (y *Module) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.Defs.ReplaceMeta(oldChild, newChild)
}

// HasGroupings
func (y *Module) GetGroupings() MetaList {
	return &y.Groupings
}
func (y *Module) GetTypedefs() MetaList {
	return &y.Typedefs
}

func (y *Module) AddInclude(i *Include) {
	y.Includes = append(y.Includes, i)
	moveModuleMeta(y, i.Module)
}

func (y *Module) AddImport(i *Import) {
	if y.Imports == nil {
		y.Imports = make(map[string]*Import)
	}
	y.Imports[i.Prefix] = i
}

func moveModuleMeta(dest *Module, src *Module) error {
	iters := []Iterator{
		ChildrenNoResolve(src.GetGroupings()),
		ChildrenNoResolve(src.GetTypedefs()),
		ChildrenNoResolve(src.DataDefs()),
	}
	for _, iter := range iters {
		for iter.HasNext() {
			if m, err := iter.Next(); err != nil {
				return err
			} else {
				dest.AddMeta(m)
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////

type Import struct {
	Prefix      string
	Revision    *Revision
	Description string
	Reference   string
	Module      *Module
}

// Identifiable
func (y *Import) GetIdent() string {
	return y.Module.GetIdent()
}

// Describable
func (y *Import) GetDescription() string {
	return y.Description
}
func (y *Import) SetDescription(d string) {
	y.Description = d
}
func (y *Import) GetReference() string {
	return y.Reference
}
func (y *Import) SetReference(r string) {
	y.Reference = r
}

////////////////////////////////////////////////////

type Include struct {
	Revision    *Revision
	Description string
	Reference   string
	Module      *Module
}

// Identifiable
func (y *Include) GetIdent() string {
	return y.Module.GetIdent()
}

// Describable
func (y *Include) GetDescription() string {
	return y.Description
}
func (y *Include) SetDescription(d string) {
	y.Description = d
}
func (y *Include) GetReference() string {
	return y.Reference
}
func (y *Include) SetReference(r string) {
	y.Reference = r
}

////////////////////////////////////////////////////

type ChoiceDecider func(Choice, ChoiceCase, interface{})

type Choice struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	details Details
}

// Identifiable
func (y *Choice) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Choice) GetDescription() string {
	return y.Description
}
func (y *Choice) SetDescription(d string) {
	y.Description = d
}
func (y *Choice) GetReference() string {
	return y.Reference
}
func (y *Choice) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Choice) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Choice) GetParent() MetaList {
	return y.Parent
}
func (y *Choice) GetSibling() Meta {
	return y.Sibling
}
func (y *Choice) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Choice) AddMeta(meta Meta) error {
	return y.linkMeta(y, meta)
}
func (y *Choice) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Choice) Clear() {
	y.Clear()
}
func (y *Choice) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// Other
func (c *Choice) GetCase(ident string) (*ChoiceCase, error) {
	m, err := FindByPathWithoutResolvingProxies(c, ident)
	if err != nil {
		return nil, err
	}
	return m.(*ChoiceCase), nil
}

// HasDetails
func (c *Choice) Details() *Details {
	return &c.details
}

////////////////////////////////////////////////////

type ChoiceCase struct {
	Ident string
	MetaBase
	ListBase
}

// Identifiable
func (y *ChoiceCase) GetIdent() string {
	return y.Ident
}

// Meta
func (y *ChoiceCase) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *ChoiceCase) GetParent() MetaList {
	return y.Parent
}
func (y *ChoiceCase) GetSibling() Meta {
	return y.Sibling
}
func (y *ChoiceCase) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *ChoiceCase) AddMeta(meta Meta) error {
	return y.linkMeta(y, meta)
}
func (y *ChoiceCase) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *ChoiceCase) Clear() {
	y.Clear()
}
func (y *ChoiceCase) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// MetaProxy
func (y *ChoiceCase) ResolveProxy() Iterator {
	return &iterator{position: y.GetFirstMeta(), resolveProxies: true}
}

////////////////////////////////////////////////////

type Revision struct {
	Ident       string
	Description string
	Reference   string
}

// Identifiable
func (y *Revision) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Revision) GetDescription() string {
	return y.Description
}
func (y *Revision) SetDescription(d string) {
	y.Description = d
}
func (y *Revision) GetReference() string {
	return y.Reference
}
func (y *Revision) SetReference(r string) {
	y.Reference = r
}

////////////////////////////////////////////////////

type Container struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs  MetaContainer
	details   Details
}

// Identifiable
func (y *Container) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Container) GetDescription() string {
	return y.Description
}
func (y *Container) SetDescription(d string) {
	y.Description = d
}
func (y *Container) GetReference() string {
	return y.Reference
}
func (y *Container) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Container) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Container) GetParent() MetaList {
	return y.Parent
}
func (y *Container) GetSibling() Meta {
	return y.Sibling
}
func (y *Container) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Container) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, meta)
	default:
		e := y.linkMeta(y, meta)
		return e
	}
}
func (y *Container) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Container) Clear() {
	y.ListBase.Clear()
}
func (y *Container) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *Container) GetGroupings() MetaList {
	return &y.Groupings
}

// HasTypedefs
func (y *Container) GetTypedefs() MetaList {
	return &y.Typedefs
}

// HasDetails
func (y *Container) Details() *Details {
	return &y.details
}

////////////////////////////////////////////////////

type List struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	Groupings   MetaContainer
	Typedefs    MetaContainer
	details     Details
	listDetails ListDetails
	Key         []string
}

// Identifiable
func (y *List) GetIdent() string {
	return y.Ident
}

// Describable
func (y *List) GetDescription() string {
	return y.Description
}
func (y *List) SetDescription(d string) {
	y.Description = d
}
func (y *List) GetReference() string {
	return y.Reference
}
func (y *List) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *List) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *List) GetParent() MetaList {
	return y.Parent
}
func (y *List) GetSibling() Meta {
	return y.Sibling
}
func (y *List) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *List) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, meta)
	default:
		return y.linkMeta(y, meta)
	}
}
func (y *List) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *List) Clear() {
	y.ListBase.Clear()
}
func (y *List) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *List) GetGroupings() MetaList {
	return &y.Groupings
}

// HasTypedefs
func (y *List) GetTypedefs() MetaList {
	return &y.Typedefs
}

// HasDetails
func (y *List) Details() *Details {
	return &y.details
}

// HasListDetails
func (y *List) ListDetails() *ListDetails {
	return &y.listDetails
}

// List
func (y *List) KeyMeta() (keyMeta []HasDataType) {
	keyMeta = make([]HasDataType, len(y.Key))
	for i, keyIdent := range y.Key {
		km, err := Find(y, keyIdent)
		keyMeta[i] = km.(HasDataType)
		// really shouldn't happen
		if err != nil {
			panic(err)
		}
	}
	return
}

////////////////////////////////////////////////////

type Leaf struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	details  Details
	DataType *DataType
}

func NewLeaf(ident string, dataType string) *Leaf {
	l := &Leaf{Ident: ident}
	l.DataType = NewDataType(l, dataType)
	return l
}

// Distinguishes the concrete type in choice-cases
func (y *Leaf) Leaf() Meta {
	return y
}

// Identifiable
func (y *Leaf) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Leaf) GetDescription() string {
	return y.Description
}
func (y *Leaf) SetDescription(d string) {
	y.Description = d
}
func (y *Leaf) GetReference() string {
	return y.Reference
}
func (y *Leaf) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Leaf) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Leaf) GetParent() MetaList {
	return y.Parent
}
func (y *Leaf) GetSibling() Meta {
	return y.Sibling
}
func (y *Leaf) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// HasDataType
func (y *Leaf) GetDataType() *DataType {
	return y.DataType
}
func (y *Leaf) SetDataType(dataType *DataType) {
	y.DataType = dataType
}
func (y *Leaf) Details() *Details {
	return &y.details
}

////////////////////////////////////////////////////

type LeafList struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	details     Details
	listDetails ListDetails
	DataType    *DataType
}

func NewLeafList(ident string, dataType string) *LeafList {
	l := &LeafList{Ident: ident}
	l.DataType = NewDataType(l, dataType)
	return l
}

// Identifiable
func (y *LeafList) GetIdent() string {
	return y.Ident
}

// Describable
func (y *LeafList) GetDescription() string {
	return y.Description
}
func (y *LeafList) SetDescription(d string) {
	y.Description = d
}
func (y *LeafList) GetReference() string {
	return y.Reference
}
func (y *LeafList) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *LeafList) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *LeafList) GetParent() MetaList {
	return y.Parent
}
func (y *LeafList) GetSibling() Meta {
	return y.Sibling
}
func (y *LeafList) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// HasType
func (y *LeafList) GetDataType() *DataType {
	return y.DataType
}
func (y *LeafList) SetDataType(dataType *DataType) {
	y.DataType = dataType
}
func (y *LeafList) Details() *Details {
	return &y.details
}

// HasListDetails
func (y *LeafList) ListDetails() *ListDetails {
	return &y.listDetails
}

////////////////////////////////////////////////////

type Any struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	details Details
	Type    *DataType
}

func NewAny(ident string) *Any {
	any := &Any{Ident: ident}
	any.Type = NewDataType(any, "any")
	return any
}

// Distinguishes the concrete type in choice-cases
func (y *Any) Any() Meta {
	return y
}

// Identifiable
func (y *Any) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Any) GetDescription() string {
	return y.Description
}
func (y *Any) SetDescription(d string) {
	y.Description = d
}
func (y *Any) GetReference() string {
	return y.Reference
}
func (y *Any) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Any) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Any) GetParent() MetaList {
	return y.Parent
}
func (y *Any) GetSibling() Meta {
	return y.Sibling
}
func (y *Any) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// HasDataType
func (y *Any) GetDataType() *DataType {
	return y.Type
}
func (y *Any) SetDataType(dataType *DataType) {
	panic("Illegal operation")
}
func (y *Any) Details() *Details {
	return &y.details
}

////////////////////////////////////////////////////

type Grouping struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	details   Details
	Groupings MetaContainer
	Typedefs  MetaContainer
}

// Identifiable
func (y *Grouping) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Grouping) GetDescription() string {
	return y.Description
}
func (y *Grouping) SetDescription(d string) {
	y.Description = d
}
func (y *Grouping) GetReference() string {
	return y.Reference
}
func (y *Grouping) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Grouping) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Grouping) GetParent() MetaList {
	return y.Parent
}
func (y *Grouping) GetSibling() Meta {
	return y.Sibling
}
func (y *Grouping) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Grouping) AddMeta(meta Meta) error {
	return y.linkMeta(y, meta)
}
func (y *Grouping) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Grouping) Clear() {
	y.ListBase.Clear()
}
func (y *Grouping) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *Grouping) GetGroupings() MetaList {
	return &y.Groupings
}

// HasTypedefs
func (y *Grouping) GetTypedefs() MetaList {
	return &y.Typedefs
}

// HasDetails
func (y *Grouping) Details() *Details {
	return &y.details
}

////////////////////////////////////////////////////

type RpcInput struct {
	MetaBase
	ListBase
	Typedefs  MetaContainer
	Groupings MetaContainer

	// Hack - not used, schema_data is incorrectly reflecting on this
	Ident       string
	Description string
}

// Identifiable
func (y *RpcInput) GetIdent() string {
	// Not technically true, but works
	return "input"
}

// Meta
func (y *RpcInput) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *RpcInput) GetParent() MetaList {
	return y.Parent
}
func (y *RpcInput) GetSibling() Meta {
	return y.Sibling
}
func (y *RpcInput) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *RpcInput) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, meta)
	default:
		return y.linkMeta(y, meta)
	}
}
func (y *RpcInput) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *RpcInput) Clear() {
	y.ListBase.Clear()
}
func (y *RpcInput) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *RpcInput) GetGroupings() MetaList {
	return &y.Groupings
}

// HasTypedefs
func (y *RpcInput) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type RpcOutput struct {
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs  MetaContainer

	// Hack - not used, schema_data is incorrectly reflecting on this
	Ident       string
	Description string
}

// Identifiable
func (y *RpcOutput) GetIdent() string {
	return "output"
}

// Meta
func (y *RpcOutput) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *RpcOutput) GetParent() MetaList {
	return y.Parent
}
func (y *RpcOutput) GetSibling() Meta {
	return y.Sibling
}
func (y *RpcOutput) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *RpcOutput) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, meta)
	default:
		return y.linkMeta(y, meta)
	}
}
func (y *RpcOutput) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *RpcOutput) Clear() {
	y.ListBase.Clear()
}
func (y *RpcOutput) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *RpcOutput) GetGroupings() MetaList {
	return &y.Groupings
}

// HasTypedefs
func (y *RpcOutput) GetTypedefsGroupings() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Rpc struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	Input  *RpcInput
	Output *RpcOutput
}

// Identifiable
func (y *Rpc) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Rpc) GetDescription() string {
	return y.Description
}
func (y *Rpc) SetDescription(d string) {
	y.Description = d
}
func (y *Rpc) GetReference() string {
	return y.Reference
}
func (y *Rpc) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Rpc) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Rpc) GetParent() MetaList {
	return y.Parent
}
func (y *Rpc) GetSibling() Meta {
	return y.Sibling
}
func (y *Rpc) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Rpc) AddMeta(meta Meta) error {
	switch t := meta.(type) {
	case *RpcInput:
		t.SetParent(y)
		y.Input = t
		return nil
	case *RpcOutput:
		t.SetParent(y)
		y.Output = t
	default:
		return &metaError{"Illegal call to add meta: rpc has fixed input and output children"}
	}
	if y.Output != nil {
		y.Input.Sibling = y.Output
	}
	return nil
}
func (y *Rpc) GetFirstMeta() Meta {
	// input and output are not official "children" of an rpc
	return nil
}
func (y *Rpc) Clear() {
	y.Input = nil
	y.Output = nil
}
func (y *Rpc) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.AddMeta(newChild)
}

////////////////////////////////////////////////////

type Notification struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	Groupings MetaContainer
	Typedefs  MetaContainer
}

// Identifiable
func (y *Notification) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Notification) GetDescription() string {
	return y.Description
}
func (y *Notification) SetDescription(d string) {
	y.Description = d
}
func (y *Notification) GetReference() string {
	return y.Reference
}
func (y *Notification) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Notification) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Notification) GetParent() MetaList {
	return y.Parent
}
func (y *Notification) GetSibling() Meta {
	return y.Sibling
}
func (y *Notification) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Notification) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Grouping:
		y.Groupings.SetParent(y)
		return y.Groupings.linkMeta(y, meta)
	default:
		return y.linkMeta(y, meta)
	}
}
func (y *Notification) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Notification) Clear() {
	y.ListBase.Clear()
}
func (y *Notification) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// HasGroupings
func (y *Notification) GetGroupings() MetaList {
	return &y.Groupings
}

// HasGroupings
func (y *Notification) GetTypedefs() MetaList {
	return &y.Typedefs
}

////////////////////////////////////////////////////

type Typedef struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	DataType *DataType
}

// Identifiable
func (y *Typedef) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Typedef) GetDescription() string {
	return y.Description
}
func (y *Typedef) SetDescription(d string) {
	y.Description = d
}
func (y *Typedef) GetReference() string {
	return y.Reference
}
func (y *Typedef) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Typedef) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Typedef) GetParent() MetaList {
	return y.Parent
}
func (y *Typedef) GetSibling() Meta {
	return y.Sibling
}
func (y *Typedef) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// HasDataType
func (y *Typedef) GetDataType() *DataType {
	return y.DataType
}

func (y *Typedef) SetDataType(dataType *DataType) {
	y.DataType = dataType
}

////////////////////////////////////////////////////

type Uses struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
	grouping *Grouping
	// if-feature
	// status
	// when
}

// Identifiable
func (y *Uses) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Uses) GetDescription() string {
	return y.Description
}
func (y *Uses) SetDescription(d string) {
	y.Description = d
}
func (y *Uses) GetReference() string {
	return y.Reference
}
func (y *Uses) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Uses) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Uses) GetParent() MetaList {
	return y.Parent
}
func (y *Uses) GetSibling() Meta {
	return y.Sibling
}
func (y *Uses) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Uses) AddMeta(meta Meta) error {
	switch meta.(type) {
	case *Refine:
		return y.linkMeta(y, meta)
	default:
		return &metaError{"Illegal call to add meta: uses only allows refine definitions"}
	}
}
func (y *Uses) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Uses) Clear() {
	y.ListBase.Clear()
}
func (y *Uses) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

func (y *Uses) FindGrouping(ident string) (*Grouping, error) {
	// lazy load grouping
	if y.grouping == nil {
		if xMod, xIdent, err := externalModule(y, ident); err != nil {
			return nil, err
		} else if xMod != nil {
			if found, err := FindByPath(xMod.GetGroupings(), xIdent); err != nil {
				return nil, err
			} else if found != nil {
				y.grouping = found.(*Grouping)
			}
		} else {
			p := y.GetParent()
			for p != nil && y.grouping == nil {
				if withGrouping, hasGrouping := p.(HasGroupings); hasGrouping {
					if found, err := FindByPath(withGrouping.GetGroupings(), y.GetIdent()); err != nil {
						return nil, err
					} else if found != nil {
						y.grouping = found.(*Grouping)
					}
				}
				p = p.GetParent()
			}
		}
	}
	return y.grouping, nil
}

func (y *Uses) refinements() map[string]*Refine {
	r := make(map[string]*Refine)
	p := y.ListBase.FirstMeta
	for p != nil {
		r[p.GetIdent()] = p.(*Refine)
		p = p.GetSibling()
	}
	return r
}

// MetaProxy
func (y *Uses) ResolveProxy() Iterator {
	// right now, not caching, but may want to consider doing so
	// as this info really shouldn't change
	if g, err := y.FindGrouping(y.Ident); err != nil {
		return nil
	} else if g != nil {
		if y.GetFirstMeta() != nil {
			if copy, err := y.refinedClone(g, ""); err != nil {
				return ErrIterator(err)
			} else {
				g = copy.(*Grouping)
			}
		}
		return Children(g)
	}
	return nil
}

func (y *Uses) refinedClone(m Meta, path string) (Meta, error) {
	// we try to be efficient here and only clone meta that has refinements
	r, requiredClone := y.requiresClone(path)
	if !requiredClone {
		return m, nil
	}
	copy := Copy(m, false)
	if r != nil {
		if err := r.refine(copy); err != nil {
			return nil, err
		}
	}
	if l, hasChildren := copy.(MetaList); hasChildren {
		i := Children(l)
		l.Clear()
		for i.HasNext() {
			child, err := i.Next()
			if err != nil {
				return nil, err
			}
			childPath := appendPath(path, child.GetIdent())
			childCopy, err := y.refinedClone(child, childPath)
			if err != nil {
				return nil, err
			}
			l.AddMeta(childCopy)
		}
	}

	return copy, nil
}

func (y *Uses) requiresClone(path string) (*Refine, bool) {
	p := y.GetFirstMeta()
	cloneRequired := false
	for p != nil {
		if strings.HasPrefix(p.GetIdent(), path) {
			cloneRequired = true
		}
		if path == p.GetIdent() {
			return p.(*Refine), true
		}
		p = p.GetSibling()
	}
	return nil, cloneRequired
}

func appendPath(parent string, child string) string {
	if parent == "" {
		return child
	}
	return parent + "/" + child
}

/////////////////////

type Refine struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	DefaultPtr  *string
	details     Details
	listDetails ListDetails
}

func (y *Refine) refine(m Meta) error {
	if y.Description != "" {
		m.(Describable).SetDescription(y.Description)
	}
	dt, hasType := m.(HasDataType)
	if y.DefaultPtr != nil {
		if !hasType {
			return c2.NewErr(fmt.Sprintf("Cannot set default on %T", m))
		}
		dt.GetDataType().DefaultPtr = y.DefaultPtr
	}
	de, hasDetails := m.(HasDetails)
	if y.details.ConfigPtr != nil {
		if !hasDetails {
			return c2.NewErr(fmt.Sprintf("Cannot set config on %T", m))
		}
		de.Details().ConfigPtr = y.details.ConfigPtr
	}
	if y.details.MandatoryPtr != nil {
		if !hasDetails {
			return c2.NewErr(fmt.Sprintf("Cannot set mandatory on %T", m))
		}
		de.Details().MandatoryPtr = y.details.MandatoryPtr
	}

	dl, hasListDetails := m.(HasListDetails)
	if y.listDetails.HasMaxElements() {
		if !hasListDetails {
			return c2.NewErr(fmt.Sprintf("Cannot set max-elements on %T", m))
		}
		dl.ListDetails().SetMaxElements(y.listDetails.MaxElements())
	}
	if y.listDetails.HasMinElements() {
		if !hasListDetails {
			return c2.NewErr(fmt.Sprintf("Cannot set min-elements on %T", m))
		}
		dl.ListDetails().SetMinElements(y.listDetails.MinElements())
	}
	if y.listDetails.ExplicitlyUnbounded() {
		if !hasListDetails {
			return c2.NewErr(fmt.Sprintf("Cannot set unbounded on %T", m))
		}
		dl.ListDetails().SetUnbounded(true)
	}

	return nil
}

// Identifiable
func (y *Refine) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Refine) GetDescription() string {
	return y.Description
}
func (y *Refine) SetDescription(d string) {
	y.Description = d
}
func (y *Refine) GetReference() string {
	return y.Reference
}
func (y *Refine) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Refine) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Refine) GetParent() MetaList {
	return y.Parent
}
func (y *Refine) GetSibling() Meta {
	return y.Sibling
}
func (y *Refine) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// HasListDetails
func (y *Refine) ListDetails() *ListDetails {
	return &y.listDetails
}

// HasDetails
func (y *Refine) Details() *Details {
	return &y.details
}

////////////////////////////////////////////////////

type Augment struct {
	Ident       string
	Description string
	Reference   string
	MetaBase
	ListBase
}

// Identifiable
func (y *Augment) GetIdent() string {
	return y.Ident
}

// Describable
func (y *Augment) GetDescription() string {
	return y.Description
}
func (y *Augment) SetDescription(d string) {
	y.Description = d
}
func (y *Augment) GetReference() string {
	return y.Reference
}
func (y *Augment) SetReference(r string) {
	y.Reference = r
}

// Meta
func (y *Augment) SetParent(parent MetaList) {
	y.Parent = parent
}
func (y *Augment) GetParent() MetaList {
	return y.Parent
}
func (y *Augment) GetSibling() Meta {
	return y.Sibling
}
func (y *Augment) SetSibling(sibling Meta) {
	y.Sibling = sibling
}

// MetaList
func (y *Augment) AddMeta(meta Meta) error {
	return y.linkMeta(y, meta)
}
func (y *Augment) GetFirstMeta() Meta {
	return y.FirstMeta
}
func (y *Augment) Clear() {
	y.ListBase.Clear()
}
func (y *Augment) ReplaceMeta(oldChild Meta, newChild Meta) error {
	return y.swapMeta(oldChild, newChild)
}

// Finalizable
func (y *Augment) Finalize() error {
	// RFC7950 Sec 7.17
	// "The target node MUST be either a container, list, choice, case, input,
	//   output, or notification node."
	target, err := FindByPath(y.Parent, y.Ident)
	if err != nil {
		return err
	}
	if target == nil {
		return c2.NewErr("augment target is not found " + y.Ident)
	}
	l, isList := target.(MetaList)
	if !isList {
		return c2.NewErr("augment target must be container, list, choice, case, input, output or notification : " + y.Ident)
	}

	return mixin(l, y)
}
