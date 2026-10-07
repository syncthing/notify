// Copyright (c) 2014-2015 The Notify Authors. All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

//go:build !linux
// +build !linux

package notify

func recursiveTreeEvents(Event) Event {
	return Create
}
