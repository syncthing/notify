// Copyright (c) 2014-2015 The Notify Authors. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

//go:build linux
// +build linux

package notify

// recursiveTreeEvents adds the platform-independent events needed to maintain
// the emulated recursive tree when a caller requests only inotify event masks.
func recursiveTreeEvents(e Event) Event {
	events := Create
	if e&(InDelete|InDeleteSelf) != 0 {
		events |= Remove
	}
	if e&(InMovedFrom|InMoveSelf) != 0 {
		events |= Rename
	}
	return events
}
