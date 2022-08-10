package gcode

import (
	"strings"
	"sync"
)

// Fragment is a collection of codes, intended to be part of a final gcode
// document
type Fragment []code

// Append another Fragment onto the current fragment
func (f Fragment) Append(others Fragment) Fragment {
	return append(f, others...)
}

// String implements the stringer interface, creating the near-final gcode
// output of the fragment. Note, this does not apply any hooks or other code
// modifications.
func (f Fragment) String() string {
	var out strings.Builder

	for _, line := range f {
		out.WriteString(line.String())
		out.WriteString("\n")
	}
	return out.String()

}

// HookFn is a function that accepts a CodePoint and returns a Fragment that
// contains a replacement definition. The original codepoint will not be part
// of the final document *unless it is returned as part of the Fragment*.
type HookFn func(CodePoint) Fragment

// Document is a gcode document or program that controls a machine of some
// kind. It supports the addition of hooks to alter the behavior of gcode on a
// per-document level. This allows the document to be changed to suit different
// machines or limitations. For example, let's say that at a general level, you
// want to change from tool 1 to tool 2. Not all machines have automatic tool
// changers and some require a different procedure to change the tool.
//
// Having a hook on the M600 command, for example, allows the original
// generator of the gcode to be unconcerned with _how_ the "filament" (in this
// case) changes, just that it changes. A machine can override the M600 command
// and insert the correct procedure (move to a specific point, pause, re-probe,
// etc.).
//
// *Any* gcode can be intercepted with a hook, with the addition of a couple
// extra ones. There are Initializer and Finalizer hooks available, to do
// exactly what their names are -- add a fragment to the beginning or end of
// the document, respectively.
type Document struct {
	lines     Fragment
	hooks     map[mapKey]HookFn
	hooksLock sync.Mutex
}

/// AddLines adds one or more lines of Codes to the document
func (d *Document) AddLines(l ...code) {
	d.lines = append(d.lines, l...)
}

// AddFragment appends the contents of a fragment to the document
func (d *Document) AddFragment(f Fragment) {
	d.AddLines(f...)
}

// String implements the Stringer interface. All hooks are applied at the
// generation of the string output.
func (d *Document) String() string {
	var out strings.Builder

	execHook := func(p CodePoint) bool {
		d.hooksLock.Lock()
		defer d.hooksLock.Unlock()

		fn, ok := d.hooks[p.Key()]
		if ok {
			if frag := fn(p); frag != nil {
				out.WriteString(frag.String())
			}
		}
		return ok
	}

	execHook(Initializer)

	for _, c := range d.lines {
		if execHook(c) {
			continue
		} else {
			out.WriteString(c.String())
		}
		out.WriteString("\n")
	}

	execHook(Finalizer)

	return out.String()
}

func (d *Document) addHook(k mapKey, fn HookFn) {
	d.hooksLock.Lock()
	defer d.hooksLock.Unlock()

	if d.hooks == nil {
		d.hooks = make(map[mapKey]HookFn)
	}

	d.hooks[k] = fn

}

type Keyable interface {
	Key() mapKey
}

// AddHook sets the hook function for the given command. There can only be one
// hook per command.
func (d *Document) AddHook(cf Keyable, fn HookFn) { d.addHook(cf.Key(), fn) }
func (d *Document) Initalizer(fn HookFn)          { d.addHook(Initializer.Key(), fn) }
func (d *Document) Finalizer(fn HookFn)           { d.addHook(Finalizer.Key(), fn) }

type CodePoint = code
