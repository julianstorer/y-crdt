package y_crdt

import (
	"testing"
)

// newTestUndoManager creates an UndoManager scoped to arr with a 500ms capture timeout.
// nil-origin transactions are always tracked (the guard only skips non-nil origins that
// aren't in the tracked-origins set).
func newTestUndoManager(arr *YArray) *UndoManager {
	return NewUndoManager(arr, 500, func(_ *Item) bool { return true }, NewSet())
}

// TestUndoStackShrinksAfterUndo verifies that UndoStack is empty after undoing the only
// tracked change. The bug: PopStackItem took stack []*StackItem by value, so the
// re-slice inside only mutated the local copy and u.UndoStack never shrank.
func TestUndoStackShrinksAfterUndo(t *testing.T) {
	doc := NewDoc("guid", false, nil, nil, false)
	arr := doc.GetArray("test")
	um := newTestUndoManager(arr)

	arr.Push(ArrayAny{"a"})

	if len(um.UndoStack) != 1 {
		t.Fatalf("expected 1 undo stack item after push, got %d", len(um.UndoStack))
	}

	um.Undo()

	if len(um.UndoStack) != 0 {
		t.Errorf("UndoStack should be empty after undo, got %d items", len(um.UndoStack))
	}
}

// TestUndoTwoIndependentChanges verifies that two changes captured in separate StackItems
// can each be independently undone. The bug: without the stack pointer fix, the second
// Undo() re-attempted the already-tombstoned top item and returned nil, so the first
// change was never reverted.
func TestUndoTwoIndependentChanges(t *testing.T) {
	doc := NewDoc("guid", false, nil, nil, false)
	ymap := doc.GetMap("test").(*YMap)
	um := NewUndoManager(ymap, 500, func(_ *Item) bool { return true }, NewSet())

	doc.Transact(func(_ *Transaction) {
		ymap.Set("a", "1")
	}, nil)
	um.StopCapturing()
	doc.Transact(func(_ *Transaction) {
		ymap.Set("b", "2")
	}, nil)

	if len(um.UndoStack) != 2 {
		t.Fatalf("expected 2 undo stack items, got %d", len(um.UndoStack))
	}

	// First undo: reverts the "b" set.
	um.Undo()
	if ymap.Get("b") != nil {
		t.Errorf("first undo: expected b=nil, got %v", ymap.Get("b"))
	}
	if ymap.Get("a") != "1" {
		t.Errorf("first undo: expected a=1, got %v", ymap.Get("a"))
	}
	if len(um.UndoStack) != 1 {
		t.Errorf("first undo: expected 1 remaining stack item, got %d", len(um.UndoStack))
	}

	// Second undo: reverts the "a" set.
	um.Undo()
	if ymap.Get("a") != nil {
		t.Errorf("second undo: expected a=nil, got %v", ymap.Get("a"))
	}
	if len(um.UndoStack) != 0 {
		t.Errorf("second undo: expected empty stack, got %d items", len(um.UndoStack))
	}
}
