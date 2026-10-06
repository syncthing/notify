// Copyright (c) 2014-2015 The Notify Authors. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

//go:build linux
// +build linux

package notify

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNotifySystemAndGlobalMix(t *testing.T) {
	n := NewNotifyTest(t, "testdata/vfs.txt")

	ch := NewChans(2)

	n.Watch("src/github.com/rjeczalik/fs", ch[0], Create)
	n.Watch("src/github.com/rjeczalik/fs", ch[1], InCreate)

	cases := []NCase{
		{
			Event:    icreate(n.W(), "src/github.com/rjeczalik/fs/.main.cc.swr"),
			Receiver: Chans{ch[0], ch[1]},
		},
	}

	n.ExpectNotifyEvents(cases, ch)
}

func TestUnknownEvent(t *testing.T) {
	n := NewNotifyTest(t, "testdata/vfs.txt")

	ch := NewChans(1)

	n.WatchErr("src/github.com/rjeczalik/fs", ch[0], nil, inExclUnlink)
}

func TestRenameDirectoryRoundTripWithInotifyEvents(t *testing.T) {
	tmpDir := t.TempDir()
	original := filepath.Join(tmpDir, "original")
	temporary := filepath.Join(tmpDir, "temporary")
	file := filepath.Join(original, "file")

	eventChan := make(chan EventInfo, 1000)
	mustT(t, Watch(tmpDir+"/...", eventChan,
		InCreate, InMovedTo, InDelete, InDeleteSelf, InModify,
		InMovedFrom, InMoveSelf, InAttrib,
	))
	defer Stop(eventChan)

	mustT(t, os.Mkdir(original, 0777))
	waitForPath(t, eventChan, original)
	time.Sleep(100 * time.Millisecond)

	mustT(t, os.Rename(original, temporary))
	waitForPath(t, eventChan, temporary)
	time.Sleep(100 * time.Millisecond)

	mustT(t, os.Rename(temporary, original))
	waitForPath(t, eventChan, original)
	time.Sleep(100 * time.Millisecond)

	mustT(t, os.WriteFile(file, []byte("content"), 0666))
	timeout := time.After(5 * time.Second)
	for {
		select {
		case event := <-eventChan:
			if filepath.Base(event.Path()) != filepath.Base(file) {
				continue
			}
			if event.Path() != file {
				t.Fatalf("event path is stale: got %q, want %q", event.Path(), file)
			}
			return
		case <-timeout:
			t.Fatal("timed out before receiving file event")
		}
	}
}
