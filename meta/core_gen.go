package meta

import (
	"fmt"
)

// This is boilerplate functions generated from ./meta/gen/ package. Do not edit
// this file, instead edit ./gen/gen.in and run "cd gen && go generate"
// Ident is identity of Module
func (m *Module) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Module) Parent() Meta {
	return m.parent
}

// Description of Module
func (m *Module) Description() string {
	return m.desc
}

func (m *Module) setDescription(desc string) {
	m.desc = desc
}

func (m *Module) Reference() string {
	return m.ref
}

func (m *Module) setReference(ref string) {
	m.ref = ref
}

func (m *Module) Extensions() []*Extension {
	return m.extensions
}

func (m *Module) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Module) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *Module) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *Module) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *Module) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *Module) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *Module) IsRecursive() bool {
	return false
}

func (m *Module) markRecursive() {
	panic("Cannot mark Module) recursive")
}



func (m *Module) Augments() []*Augment {
	return m.augments
}

func (m *Module) addAugments(a *Augment) {
	m.augments = append(m.augments, a)
}

func (m *Module) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *Module) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *Module) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *Module) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *Module) Actions() map[string]*Rpc {
	return m.actions
}

func (m *Module) addAction(a *Rpc) {
	if m.actions == nil {
		m.actions = make(map[string]*Rpc)
	}
    m.actions[a.Ident()] = a
}

func (m *Module) setActions(actions map[string]*Rpc) {
	m.actions = actions
}

func (m *Module) Notifications() map[string]*Notification {
	return m.notifications
}

func (m *Module) addNotification(n *Notification) {
	if m.notifications == nil {
		m.notifications = make(map[string]*Notification)
	}
    m.notifications[n.Ident()] = n
}

func (m *Module) setNotifications(notifications map[string]*Notification) {
	m.notifications = notifications
}


// Definition can be a data defintion, action or notification
func (m *Module) Definition(ident string) Definition {
	if x, found := m.notifications[ident]; found {
		return x
	}
	
	if x, found := m.actions[ident]; found {
		return x
	}
	
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *Module) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.notifications != nil {
		copy.notifications = make(map[string]*Notification, len(m.notifications))
		for ident, notif := range m.notifications {
			copy.notifications[ident] = notif.clone(&copy).(*Notification)
		}
	}
	
	if m.actions != nil {
		copy.actions = make(map[string]*Rpc, len(m.actions))
		for ident, action := range m.actions {
			copy.actions[ident] = action.clone(&copy).(*Rpc)
		}
	}
	
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	

	return &copy
}


// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Import) Parent() Meta {
	return m.parent
}

// Description of Import
func (m *Import) Description() string {
	return m.desc
}

func (m *Import) setDescription(desc string) {
	m.desc = desc
}

func (m *Import) Reference() string {
	return m.ref
}

func (m *Import) setReference(ref string) {
	m.ref = ref
}

func (m *Import) Extensions() []*Extension {
	return m.extensions
}

func (m *Import) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Include) Parent() Meta {
	return m.parent
}

// Description of Include
func (m *Include) Description() string {
	return m.desc
}

func (m *Include) setDescription(desc string) {
	m.desc = desc
}

func (m *Include) Reference() string {
	return m.ref
}

func (m *Include) setReference(ref string) {
	m.ref = ref
}

func (m *Include) Extensions() []*Extension {
	return m.extensions
}

func (m *Include) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Ident is identity of Choice
func (m *Choice) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Choice) Parent() Meta {
	return m.parent
}

// Description of Choice
func (m *Choice) Description() string {
	return m.desc
}

func (m *Choice) setDescription(desc string) {
	m.desc = desc
}

func (m *Choice) Reference() string {
	return m.ref
}

func (m *Choice) setReference(ref string) {
	m.ref = ref
}

func (m *Choice) Status() Status {
	return m.status
}

func (m *Choice) setStatus(status Status) {
	m.status = status
}

func (m *Choice) Extensions() []*Extension {
	return m.extensions
}

func (m *Choice) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Choice) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Choice) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Choice) When() *When {
	return m.when
}

func (m *Choice) setWhen(w *When) {
    m.when = w
}

func (m *Choice) Config() bool {
	return *m.configPtr
}

func (m *Choice) setConfig(c bool) {
	m.configPtr = &c
}

func (m *Choice) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *Choice) Mandatory() bool {
	return m.mandatory
}

func (m *Choice) setMandatory(b bool) {
	m.mandatory = b
}

func (m *Choice) scopedParent() Meta {
	return m.scope
}

func (m *Choice) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.cases != nil {
		copy.cases = make(map[string]*ChoiceCase, len(m.cases))
		for ident, kase := range m.cases {
			copy.cases[ident] = kase.clone(&copy).(*ChoiceCase)
		}
	}
	

	return &copy
}


// Ident is identity of ChoiceCase
func (m *ChoiceCase) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *ChoiceCase) Parent() Meta {
	return m.parent
}

// Description of ChoiceCase
func (m *ChoiceCase) Description() string {
	return m.desc
}

func (m *ChoiceCase) setDescription(desc string) {
	m.desc = desc
}

func (m *ChoiceCase) Reference() string {
	return m.ref
}

func (m *ChoiceCase) setReference(ref string) {
	m.ref = ref
}

func (m *ChoiceCase) Extensions() []*Extension {
	return m.extensions
}

func (m *ChoiceCase) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *ChoiceCase) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *ChoiceCase) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *ChoiceCase) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *ChoiceCase) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *ChoiceCase) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *ChoiceCase) IsRecursive() bool {
	return m.recursive
}

func (m *ChoiceCase) markRecursive() {
	m.recursive = true
}


func (m *ChoiceCase) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *ChoiceCase) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *ChoiceCase) When() *When {
	return m.when
}

func (m *ChoiceCase) setWhen(w *When) {
    m.when = w
}

// Definition can be a data defintion, action or notification
func (m *ChoiceCase) Definition(ident string) Definition {
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *ChoiceCase) scopedParent() Meta {
	return m.scope
}

func (m *ChoiceCase) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	

	return &copy
}


// Ident is identity of Revision
func (m *Revision) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Revision) Parent() Meta {
	return m.parent
}

// Description of Revision
func (m *Revision) Description() string {
	return m.desc
}

func (m *Revision) setDescription(desc string) {
	m.desc = desc
}

func (m *Revision) Reference() string {
	return m.ref
}

func (m *Revision) setReference(ref string) {
	m.ref = ref
}

func (m *Revision) Extensions() []*Extension {
	return m.extensions
}

func (m *Revision) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Revision) scopedParent() Meta {
	return m.scope
}


// Ident is identity of Container
func (m *Container) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Container) Parent() Meta {
	return m.parent
}

// Description of Container
func (m *Container) Description() string {
	return m.desc
}

func (m *Container) setDescription(desc string) {
	m.desc = desc
}

func (m *Container) Reference() string {
	return m.ref
}

func (m *Container) setReference(ref string) {
	m.ref = ref
}

func (m *Container) Status() Status {
	return m.status
}

func (m *Container) setStatus(status Status) {
	m.status = status
}

func (m *Container) Extensions() []*Extension {
	return m.extensions
}

func (m *Container) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Container) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *Container) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *Container) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *Container) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *Container) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *Container) IsRecursive() bool {
	return m.recursive
}

func (m *Container) markRecursive() {
	m.recursive = true
}


func (m *Container) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *Container) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *Container) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *Container) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *Container) Musts() []*Must {
	return m.musts
}

func (m *Container) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *Container) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Container) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Container) When() *When {
	return m.when
}

func (m *Container) setWhen(w *When) {
    m.when = w
}

func (m *Container) Actions() map[string]*Rpc {
	return m.actions
}

func (m *Container) addAction(a *Rpc) {
	if m.actions == nil {
		m.actions = make(map[string]*Rpc)
	}
    m.actions[a.Ident()] = a
}

func (m *Container) setActions(actions map[string]*Rpc) {
	m.actions = actions
}

func (m *Container) Notifications() map[string]*Notification {
	return m.notifications
}

func (m *Container) addNotification(n *Notification) {
	if m.notifications == nil {
		m.notifications = make(map[string]*Notification)
	}
    m.notifications[n.Ident()] = n
}

func (m *Container) setNotifications(notifications map[string]*Notification) {
	m.notifications = notifications
}


// Definition can be a data defintion, action or notification
func (m *Container) Definition(ident string) Definition {
	if x, found := m.notifications[ident]; found {
		return x
	}
	
	if x, found := m.actions[ident]; found {
		return x
	}
	
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *Container) Config() bool {
	return *m.configPtr
}

func (m *Container) setConfig(c bool) {
	m.configPtr = &c
}

func (m *Container) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *Container) Mandatory() bool {
	return m.mandatory
}

func (m *Container) setMandatory(b bool) {
	m.mandatory = b
}

func (m *Container) scopedParent() Meta {
	return m.scope
}

func (m *Container) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.notifications != nil {
		copy.notifications = make(map[string]*Notification, len(m.notifications))
		for ident, notif := range m.notifications {
			copy.notifications[ident] = notif.clone(&copy).(*Notification)
		}
	}
	
	if m.actions != nil {
		copy.actions = make(map[string]*Rpc, len(m.actions))
		for ident, action := range m.actions {
			copy.actions[ident] = action.clone(&copy).(*Rpc)
		}
	}
	
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of List
func (m *List) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *List) Parent() Meta {
	return m.parent
}

// Description of List
func (m *List) Description() string {
	return m.desc
}

func (m *List) setDescription(desc string) {
	m.desc = desc
}

func (m *List) Reference() string {
	return m.ref
}

func (m *List) setReference(ref string) {
	m.ref = ref
}

func (m *List) Extensions() []*Extension {
	return m.extensions
}

func (m *List) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *List) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *List) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *List) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *List) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *List) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *List) IsRecursive() bool {
	return m.recursive
}

func (m *List) markRecursive() {
	m.recursive = true
}


func (m *List) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *List) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *List) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *List) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *List) Musts() []*Must {
	return m.musts
}

func (m *List) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *List) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *List) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *List) When() *When {
	return m.when
}

func (m *List) setWhen(w *When) {
    m.when = w
}

func (m *List) Actions() map[string]*Rpc {
	return m.actions
}

func (m *List) addAction(a *Rpc) {
	if m.actions == nil {
		m.actions = make(map[string]*Rpc)
	}
    m.actions[a.Ident()] = a
}

func (m *List) setActions(actions map[string]*Rpc) {
	m.actions = actions
}

func (m *List) Notifications() map[string]*Notification {
	return m.notifications
}

func (m *List) addNotification(n *Notification) {
	if m.notifications == nil {
		m.notifications = make(map[string]*Notification)
	}
    m.notifications[n.Ident()] = n
}

func (m *List) setNotifications(notifications map[string]*Notification) {
	m.notifications = notifications
}


// Definition can be a data defintion, action or notification
func (m *List) Definition(ident string) Definition {
	if x, found := m.notifications[ident]; found {
		return x
	}
	
	if x, found := m.actions[ident]; found {
		return x
	}
	
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *List) Config() bool {
	return *m.configPtr
}

func (m *List) setConfig(c bool) {
	m.configPtr = &c
}

func (m *List) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *List) Mandatory() bool {
	return m.mandatory
}

func (m *List) setMandatory(b bool) {
	m.mandatory = b
}

func (m *List) MinElements() int { 
	return m.minElements
}

func (m *List) setMinElements(i int) {
	m.minElements = i
}

func (m *List) MaxElements() int { 
	return m.maxElements
}

func (m *List) setMaxElements(i int) {
	m.maxElements = i
}

func (m *List) Unbounded() bool { 
	if m.unboundedPtr != nil {
		return *m.unboundedPtr
	}
	return m.maxElements == 0
}

func (m *List) setUnbounded(b bool) {
	m.unboundedPtr = &b
}


func (m *List) scopedParent() Meta {
	return m.scope
}

func (m *List) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.notifications != nil {
		copy.notifications = make(map[string]*Notification, len(m.notifications))
		for ident, notif := range m.notifications {
			copy.notifications[ident] = notif.clone(&copy).(*Notification)
		}
	}
	
	if m.actions != nil {
		copy.actions = make(map[string]*Rpc, len(m.actions))
		for ident, action := range m.actions {
			copy.actions[ident] = action.clone(&copy).(*Rpc)
		}
	}
	
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of Leaf
func (m *Leaf) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Leaf) Parent() Meta {
	return m.parent
}

// Description of Leaf
func (m *Leaf) Description() string {
	return m.desc
}

func (m *Leaf) setDescription(desc string) {
	m.desc = desc
}

func (m *Leaf) Reference() string {
	return m.ref
}

func (m *Leaf) setReference(ref string) {
	m.ref = ref
}

func (m *Leaf) Extensions() []*Extension {
	return m.extensions
}

func (m *Leaf) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Leaf) Musts() []*Must {
	return m.musts
}

func (m *Leaf) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *Leaf) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Leaf) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Leaf) When() *When {
	return m.when
}

func (m *Leaf) setWhen(w *When) {
    m.when = w
}

func (m *Leaf) Config() bool {
	return *m.configPtr
}

func (m *Leaf) setConfig(c bool) {
	m.configPtr = &c
}

func (m *Leaf) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *Leaf) Mandatory() bool {
	return m.mandatory
}

func (m *Leaf) setMandatory(b bool) {
	m.mandatory = b
}

func (m *Leaf) Type() *Type { 
	return m.dtype
}

func (m *Leaf) setType(t *Type) {
	m.dtype = t
}

func (m *Leaf) Units() string{
	return m.units
}

func (m *Leaf) setUnits(u string) {
    m.units = u
}

func (m *Leaf) Default() interface{} {
	return m.defaultVal
}

func (m *Leaf) HasDefault() bool {
	return m.defaultVal != nil
}

func (m *Leaf) setDefault(d interface{}) {
    m.defaultVal = d
}

func (m *Leaf) scopedParent() Meta {
	return m.scope
}

func (m *Leaf) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of LeafList
func (m *LeafList) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *LeafList) Parent() Meta {
	return m.parent
}

// Description of LeafList
func (m *LeafList) Description() string {
	return m.desc
}

func (m *LeafList) setDescription(desc string) {
	m.desc = desc
}

func (m *LeafList) Reference() string {
	return m.ref
}

func (m *LeafList) setReference(ref string) {
	m.ref = ref
}

func (m *LeafList) Extensions() []*Extension {
	return m.extensions
}

func (m *LeafList) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *LeafList) Musts() []*Must {
	return m.musts
}

func (m *LeafList) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *LeafList) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *LeafList) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *LeafList) When() *When {
	return m.when
}

func (m *LeafList) setWhen(w *When) {
    m.when = w
}

func (m *LeafList) Config() bool {
	return *m.configPtr
}

func (m *LeafList) setConfig(c bool) {
	m.configPtr = &c
}

func (m *LeafList) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *LeafList) Mandatory() bool {
	return m.mandatory
}

func (m *LeafList) setMandatory(b bool) {
	m.mandatory = b
}

func (m *LeafList) MinElements() int { 
	return m.minElements
}

func (m *LeafList) setMinElements(i int) {
	m.minElements = i
}

func (m *LeafList) MaxElements() int { 
	return m.maxElements
}

func (m *LeafList) setMaxElements(i int) {
	m.maxElements = i
}

func (m *LeafList) Unbounded() bool { 
	if m.unboundedPtr != nil {
		return *m.unboundedPtr
	}
	return m.maxElements == 0
}

func (m *LeafList) setUnbounded(b bool) {
	m.unboundedPtr = &b
}


func (m *LeafList) Type() *Type { 
	return m.dtype
}

func (m *LeafList) setType(t *Type) {
	m.dtype = t
}

func (m *LeafList) Units() string{
	return m.units
}

func (m *LeafList) setUnits(u string) {
    m.units = u
}

func (m *LeafList) Default() interface{} {
	return m.defaultVal
}

func (m *LeafList) HasDefault() bool {
	return m.defaultVal != nil
}

func (m *LeafList) setDefault(d interface{}) {
    m.defaultVal = d
}

func (m *LeafList) scopedParent() Meta {
	return m.scope
}

func (m *LeafList) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of Any
func (m *Any) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Any) Parent() Meta {
	return m.parent
}

// Description of Any
func (m *Any) Description() string {
	return m.desc
}

func (m *Any) setDescription(desc string) {
	m.desc = desc
}

func (m *Any) Reference() string {
	return m.ref
}

func (m *Any) setReference(ref string) {
	m.ref = ref
}

func (m *Any) Extensions() []*Extension {
	return m.extensions
}

func (m *Any) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Any) Musts() []*Must {
	return m.musts
}

func (m *Any) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *Any) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Any) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Any) When() *When {
	return m.when
}

func (m *Any) setWhen(w *When) {
    m.when = w
}

func (m *Any) Config() bool {
	return *m.configPtr
}

func (m *Any) setConfig(c bool) {
	m.configPtr = &c
}

func (m *Any) isConfigSet() bool {
	return m.configPtr != nil
}

func (m *Any) Mandatory() bool {
	return m.mandatory
}

func (m *Any) setMandatory(b bool) {
	m.mandatory = b
}

func (m *Any) scopedParent() Meta {
	return m.scope
}

func (m *Any) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of Grouping
func (m *Grouping) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Grouping) Parent() Meta {
	return m.parent
}

// Description of Grouping
func (m *Grouping) Description() string {
	return m.desc
}

func (m *Grouping) setDescription(desc string) {
	m.desc = desc
}

func (m *Grouping) Reference() string {
	return m.ref
}

func (m *Grouping) setReference(ref string) {
	m.ref = ref
}

func (m *Grouping) Extensions() []*Extension {
	return m.extensions
}

func (m *Grouping) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Grouping) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *Grouping) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *Grouping) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *Grouping) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *Grouping) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *Grouping) IsRecursive() bool {
	return false
}

func (m *Grouping) markRecursive() {
	panic("Cannot mark Grouping) recursive")
}



func (m *Grouping) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *Grouping) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *Grouping) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *Grouping) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *Grouping) Actions() map[string]*Rpc {
	return m.actions
}

func (m *Grouping) addAction(a *Rpc) {
	if m.actions == nil {
		m.actions = make(map[string]*Rpc)
	}
    m.actions[a.Ident()] = a
}

func (m *Grouping) setActions(actions map[string]*Rpc) {
	m.actions = actions
}

func (m *Grouping) Notifications() map[string]*Notification {
	return m.notifications
}

func (m *Grouping) addNotification(n *Notification) {
	if m.notifications == nil {
		m.notifications = make(map[string]*Notification)
	}
    m.notifications[n.Ident()] = n
}

func (m *Grouping) setNotifications(notifications map[string]*Notification) {
	m.notifications = notifications
}


// Definition can be a data defintion, action or notification
func (m *Grouping) Definition(ident string) Definition {
	if x, found := m.notifications[ident]; found {
		return x
	}
	
	if x, found := m.actions[ident]; found {
		return x
	}
	
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *Grouping) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.notifications != nil {
		copy.notifications = make(map[string]*Notification, len(m.notifications))
		for ident, notif := range m.notifications {
			copy.notifications[ident] = notif.clone(&copy).(*Notification)
		}
	}
	
	if m.actions != nil {
		copy.actions = make(map[string]*Rpc, len(m.actions))
		for ident, action := range m.actions {
			copy.actions[ident] = action.clone(&copy).(*Rpc)
		}
	}
	
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	

	return &copy
}


// Ident is identity of Uses
func (m *Uses) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Uses) Parent() Meta {
	return m.parent
}

// Description of Uses
func (m *Uses) Description() string {
	return m.desc
}

func (m *Uses) setDescription(desc string) {
	m.desc = desc
}

func (m *Uses) Reference() string {
	return m.ref
}

func (m *Uses) setReference(ref string) {
	m.ref = ref
}

func (m *Uses) Extensions() []*Extension {
	return m.extensions
}

func (m *Uses) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Uses) Augments() []*Augment {
	return m.augments
}

func (m *Uses) addAugments(a *Augment) {
	m.augments = append(m.augments, a)
}

func (m *Uses) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Uses) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Uses) When() *When {
	return m.when
}

func (m *Uses) setWhen(w *When) {
    m.when = w
}

func (m *Uses) scopedParent() Meta {
	return m.scope
}

func (m *Uses) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent

	return &copy
}


// Ident is identity of Refine
func (m *Refine) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Refine) Parent() Meta {
	return m.parent
}

// Description of Refine
func (m *Refine) Description() string {
	return m.desc
}

func (m *Refine) setDescription(desc string) {
	m.desc = desc
}

func (m *Refine) Reference() string {
	return m.ref
}

func (m *Refine) setReference(ref string) {
	m.ref = ref
}

func (m *Refine) Extensions() []*Extension {
	return m.extensions
}

func (m *Refine) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Refine) Musts() []*Must {
	return m.musts
}

func (m *Refine) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *Refine) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Refine) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}


// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *RpcInput) Parent() Meta {
	return m.parent
}

// Description of RpcInput
func (m *RpcInput) Description() string {
	return m.desc
}

func (m *RpcInput) setDescription(desc string) {
	m.desc = desc
}

func (m *RpcInput) Reference() string {
	return m.ref
}

func (m *RpcInput) setReference(ref string) {
	m.ref = ref
}

func (m *RpcInput) Extensions() []*Extension {
	return m.extensions
}

func (m *RpcInput) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *RpcInput) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *RpcInput) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *RpcInput) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *RpcInput) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *RpcInput) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *RpcInput) IsRecursive() bool {
	return false
}

func (m *RpcInput) markRecursive() {
	panic("Cannot mark RpcInput) recursive")
}



func (m *RpcInput) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *RpcInput) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *RpcInput) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *RpcInput) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *RpcInput) Musts() []*Must {
	return m.musts
}

func (m *RpcInput) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *RpcInput) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *RpcInput) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

// Definition can be a data defintion, action or notification
func (m *RpcInput) Definition(ident string) Definition {
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *RpcInput) scopedParent() Meta {
	return m.scope
}

func (m *RpcInput) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *RpcOutput) Parent() Meta {
	return m.parent
}

// Description of RpcOutput
func (m *RpcOutput) Description() string {
	return m.desc
}

func (m *RpcOutput) setDescription(desc string) {
	m.desc = desc
}

func (m *RpcOutput) Reference() string {
	return m.ref
}

func (m *RpcOutput) setReference(ref string) {
	m.ref = ref
}

func (m *RpcOutput) Extensions() []*Extension {
	return m.extensions
}

func (m *RpcOutput) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *RpcOutput) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *RpcOutput) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *RpcOutput) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *RpcOutput) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *RpcOutput) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *RpcOutput) IsRecursive() bool {
	return false
}

func (m *RpcOutput) markRecursive() {
	panic("Cannot mark RpcOutput) recursive")
}



func (m *RpcOutput) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *RpcOutput) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *RpcOutput) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *RpcOutput) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *RpcOutput) Musts() []*Must {
	return m.musts
}

func (m *RpcOutput) addMust(x *Must) {
    m.musts = append(m.musts, x)
}

func (m *RpcOutput) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *RpcOutput) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

// Definition can be a data defintion, action or notification
func (m *RpcOutput) Definition(ident string) Definition {
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *RpcOutput) scopedParent() Meta {
	return m.scope
}

func (m *RpcOutput) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	
	if m.musts != nil {
		copy.musts = make([]*Must, len(m.musts))
		for i, must := range m.musts {
			copy.musts[i] = must.clone(&copy).(*Must)
		}
	}
	

	return &copy
}


// Ident is identity of Rpc
func (m *Rpc) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Rpc) Parent() Meta {
	return m.parent
}

// Description of Rpc
func (m *Rpc) Description() string {
	return m.desc
}

func (m *Rpc) setDescription(desc string) {
	m.desc = desc
}

func (m *Rpc) Reference() string {
	return m.ref
}

func (m *Rpc) setReference(ref string) {
	m.ref = ref
}

func (m *Rpc) Extensions() []*Extension {
	return m.extensions
}

func (m *Rpc) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Rpc) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *Rpc) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *Rpc) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *Rpc) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *Rpc) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Rpc) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Rpc) scopedParent() Meta {
	return m.scope
}

func (m *Rpc) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.input != nil {
		copy.input = m.input.clone(&copy).(*RpcInput)
	}
	if m.output != nil {
		copy.output = m.output.clone(&copy).(*RpcOutput)
	}
	

	return &copy
}


// Ident is identity of Notification
func (m *Notification) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Notification) Parent() Meta {
	return m.parent
}

// Description of Notification
func (m *Notification) Description() string {
	return m.desc
}

func (m *Notification) setDescription(desc string) {
	m.desc = desc
}

func (m *Notification) Reference() string {
	return m.ref
}

func (m *Notification) setReference(ref string) {
	m.ref = ref
}

func (m *Notification) Extensions() []*Extension {
	return m.extensions
}

func (m *Notification) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Notification) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *Notification) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *Notification) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *Notification) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *Notification) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *Notification) IsRecursive() bool {
	return false
}

func (m *Notification) markRecursive() {
	panic("Cannot mark Notification) recursive")
}



func (m *Notification) Groupings() map[string]*Grouping {
	return m.groupings
}

func (m *Notification) addGrouping(g *Grouping) {
	if m.groupings == nil {
		m.groupings = make(map[string]*Grouping)
	}
    m.groupings[g.Ident()] = g
}

func (m *Notification) Typedefs() map[string]*Typedef {
	return m.typedefs
}

func (m *Notification) addTypedef(t *Typedef) {
	if m.typedefs == nil {
		m.typedefs = make(map[string]*Typedef)
	}
    m.typedefs[t.Ident()] = t
}

func (m *Notification) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Notification) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

// Definition can be a data defintion, action or notification
func (m *Notification) Definition(ident string) Definition {
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *Notification) scopedParent() Meta {
	return m.scope
}

func (m *Notification) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	

	return &copy
}


// Ident is identity of Typedef
func (m *Typedef) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Typedef) Parent() Meta {
	return m.parent
}

// Description of Typedef
func (m *Typedef) Description() string {
	return m.desc
}

func (m *Typedef) setDescription(desc string) {
	m.desc = desc
}

func (m *Typedef) Reference() string {
	return m.ref
}

func (m *Typedef) setReference(ref string) {
	m.ref = ref
}

func (m *Typedef) Extensions() []*Extension {
	return m.extensions
}

func (m *Typedef) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Typedef) Type() *Type { 
	return m.dtype
}

func (m *Typedef) setType(t *Type) {
	m.dtype = t
}

func (m *Typedef) Units() string{
	return m.units
}

func (m *Typedef) setUnits(u string) {
    m.units = u
}

func (m *Typedef) Default() interface{} {
	return m.defaultVal
}

func (m *Typedef) HasDefault() bool {
	return m.defaultVal != nil
}

func (m *Typedef) setDefault(d interface{}) {
    m.defaultVal = d
}


// Ident is identity of Augment
func (m *Augment) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Augment) Parent() Meta {
	return m.parent
}

// Description of Augment
func (m *Augment) Description() string {
	return m.desc
}

func (m *Augment) setDescription(desc string) {
	m.desc = desc
}

func (m *Augment) Reference() string {
	return m.ref
}

func (m *Augment) setReference(ref string) {
	m.ref = ref
}

func (m *Augment) Extensions() []*Extension {
	return m.extensions
}

func (m *Augment) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Augment) DataDefinitions() []Definition {
	return m.dataDefs
}

func (m *Augment) DataDefinition(ident string) Definition {
	return m.dataDefsIndex[ident]
}

func (m *Augment) addDataDefinition(d Definition) {
	if c, isChoice := d.(*Choice); isChoice {
		for _, k := range c.Cases() {
			for _, kdef := range k.DataDefinitions() {
				// recurse in case it's another choice
				m.indexDataDefinition(kdef)
			}
		}
 	} else {
		m.indexDataDefinition(d)
	 }
	m.dataDefs = append(m.dataDefs, d)
}

func (m *Augment) indexDataDefinition(def Definition) {
	if m.dataDefsIndex == nil {
		m.dataDefsIndex = make(map[string]Definition)
	}
	if _, exists := m.dataDefsIndex[def.Ident()]; exists {
		// TODO: make this an error
		panic(fmt.Sprintf("Conflict adding add %s to %s. ", def.Ident(), m.Ident()))
	}	
	m.dataDefsIndex[def.Ident()] = def
}

func (m *Augment) popDataDefinitions() []Definition {
	orig := m.dataDefs
	m.dataDefs = make([]Definition, 0, len(orig))
	for key := range m.dataDefsIndex {
		delete(m.dataDefsIndex, key)
	}
	return orig
}
func (m *Augment) IsRecursive() bool {
	return false
}

func (m *Augment) markRecursive() {
	panic("Cannot mark Augment) recursive")
}



func (m *Augment) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Augment) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}

func (m *Augment) When() *When {
	return m.when
}

func (m *Augment) setWhen(w *When) {
    m.when = w
}

func (m *Augment) Actions() map[string]*Rpc {
	return m.actions
}

func (m *Augment) addAction(a *Rpc) {
	if m.actions == nil {
		m.actions = make(map[string]*Rpc)
	}
    m.actions[a.Ident()] = a
}

func (m *Augment) setActions(actions map[string]*Rpc) {
	m.actions = actions
}

func (m *Augment) Notifications() map[string]*Notification {
	return m.notifications
}

func (m *Augment) addNotification(n *Notification) {
	if m.notifications == nil {
		m.notifications = make(map[string]*Notification)
	}
    m.notifications[n.Ident()] = n
}

func (m *Augment) setNotifications(notifications map[string]*Notification) {
	m.notifications = notifications
}


// Definition can be a data defintion, action or notification
func (m *Augment) Definition(ident string) Definition {
	if x, found := m.notifications[ident]; found {
		return x
	}
	
	if x, found := m.actions[ident]; found {
		return x
	}
	
	if x, found := m.dataDefsIndex[ident]; found {
		return x
	}
	
	return nil
}

func (m *Augment) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent
	if m.notifications != nil {
		copy.notifications = make(map[string]*Notification, len(m.notifications))
		for ident, notif := range m.notifications {
			copy.notifications[ident] = notif.clone(&copy).(*Notification)
		}
	}
	
	if m.actions != nil {
		copy.actions = make(map[string]*Rpc, len(m.actions))
		for ident, action := range m.actions {
			copy.actions[ident] = action.clone(&copy).(*Rpc)
		}
	}
	
	if m.dataDefs != nil {
		copy.dataDefs = make([]Definition, len(m.dataDefs))
		copy.dataDefsIndex = make(map[string]Definition, len(m.dataDefs))
		for i, def := range m.dataDefs {
			copyDef := def.(cloneable).clone(&copy).(Definition)
			copy.dataDefs[i] = copyDef
			copy.dataDefsIndex[def.Ident()] = copyDef
		}
	}
	

	return &copy
}


// Ident is identity of Type
func (m *Type) Ident() string {
	return m.ident
}

// Description of Type
func (m *Type) Description() string {
	return m.desc
}

func (m *Type) setDescription(desc string) {
	m.desc = desc
}

func (m *Type) Reference() string {
	return m.ref
}

func (m *Type) setReference(ref string) {
	m.ref = ref
}

func (m *Type) Extensions() []*Extension {
	return m.extensions
}

func (m *Type) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Ident is identity of Identity
func (m *Identity) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Identity) Parent() Meta {
	return m.parent
}

// Description of Identity
func (m *Identity) Description() string {
	return m.desc
}

func (m *Identity) setDescription(desc string) {
	m.desc = desc
}

func (m *Identity) Reference() string {
	return m.ref
}

func (m *Identity) setReference(ref string) {
	m.ref = ref
}

func (m *Identity) Extensions() []*Extension {
	return m.extensions
}

func (m *Identity) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Identity) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Identity) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}


// Ident is identity of Feature
func (m *Feature) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Feature) Parent() Meta {
	return m.parent
}

// Description of Feature
func (m *Feature) Description() string {
	return m.desc
}

func (m *Feature) setDescription(desc string) {
	m.desc = desc
}

func (m *Feature) Reference() string {
	return m.ref
}

func (m *Feature) setReference(ref string) {
	m.ref = ref
}

func (m *Feature) Extensions() []*Extension {
	return m.extensions
}

func (m *Feature) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Feature) IfFeatures() []*IfFeature {
	return m.ifs
}

func (m *Feature) addIfFeature(i *IfFeature) {
    m.ifs = append(m.ifs, i)
}


// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *IfFeature) Parent() Meta {
	return m.parent
}

func (m *IfFeature) Extensions() []*Extension {
	return m.extensions
}

func (m *IfFeature) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}




// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *When) Parent() Meta {
	return m.parent
}

// Description of When
func (m *When) Description() string {
	return m.desc
}

func (m *When) setDescription(desc string) {
	m.desc = desc
}

func (m *When) Reference() string {
	return m.ref
}

func (m *When) setReference(ref string) {
	m.ref = ref
}

func (m *When) Extensions() []*Extension {
	return m.extensions
}

func (m *When) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Must) Parent() Meta {
	return m.parent
}

func (m *Must) Extensions() []*Extension {
	return m.extensions
}

func (m *Must) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}


func (m *Must) clone(parent Meta) interface{} {
	copy := *m
	copy.parent = parent

	return &copy
}


// Ident is identity of ExtensionDef
func (m *ExtensionDef) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *ExtensionDef) Parent() Meta {
	return m.parent
}

// Description of ExtensionDef
func (m *ExtensionDef) Description() string {
	return m.desc
}

func (m *ExtensionDef) setDescription(desc string) {
	m.desc = desc
}

func (m *ExtensionDef) Reference() string {
	return m.ref
}

func (m *ExtensionDef) setReference(ref string) {
	m.ref = ref
}

func (m *ExtensionDef) Status() Status {
	return m.status
}

func (m *ExtensionDef) setStatus(status Status) {
	m.status = status
}

func (m *ExtensionDef) Extensions() []*Extension {
	return m.extensions
}

func (m *ExtensionDef) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Ident is identity of ExtensionDefArg
func (m *ExtensionDefArg) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *ExtensionDefArg) Parent() Meta {
	return m.parent
}

// Description of ExtensionDefArg
func (m *ExtensionDefArg) Description() string {
	return m.desc
}

func (m *ExtensionDefArg) setDescription(desc string) {
	m.desc = desc
}

func (m *ExtensionDefArg) Reference() string {
	return m.ref
}

func (m *ExtensionDefArg) setReference(ref string) {
	m.ref = ref
}

func (m *ExtensionDefArg) Extensions() []*Extension {
	return m.extensions
}

func (m *ExtensionDefArg) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Ident is identity of Extension
func (m *Extension) Ident() string {
	return m.ident
}

// Parent is where this extension is define unless the extension is a
// secondary extension like a description and then this is the parent
// of that description
func (m *Extension) Parent() Meta {
	return m.parent
}

func (m *Extension) Extensions() []*Extension {
	return m.extensions
}

func (m *Extension) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



// Ident is identity of Enum
func (m *Enum) Ident() string {
	return m.ident
}

// Description of Enum
func (m *Enum) Description() string {
	return m.desc
}

func (m *Enum) setDescription(desc string) {
	m.desc = desc
}

func (m *Enum) Reference() string {
	return m.ref
}

func (m *Enum) setReference(ref string) {
	m.ref = ref
}

func (m *Enum) Extensions() []*Extension {
	return m.extensions
}

func (m *Enum) addExtension(extension *Extension) {
	m.extensions = append(m.extensions, extension)
}



