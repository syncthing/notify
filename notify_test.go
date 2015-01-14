package notify

import "testing"

func TestNotifyExample(t *testing.T) {
	n := NewNotifyTest(t, "testdata/gopath.txt")
	defer n.Close()

	ch := NewChans(3)

	// Watch-points can be set explicitely via Watch/Stop calls...
	n.Watch("src/github.com/rjeczalik/fs", ch[0], Write)
	n.Watch("src/github.com/rjeczalik/which", ch[0], Write)
	n.Watch("src/github.com/rjeczalik/which/...", ch[1], Create)
	n.Watch("src/github.com/rjeczalik/fs/cmd/...", ch[2], Delete)

	cases := []NCase{
		{
			Event:    write(n.W(), "src/github.com/rjeczalik/fs/fs.go", []byte("XD")),
			Receiver: Chans{ch[0]},
		},
		{
			Event:    write(n.W(), "src/github.com/rjeczalik/which/README.md", []byte("XD")),
			Receiver: Chans{ch[0]},
		},
		{
			Event:    write(n.W(), "src/github.com/rjeczalik/fs/cmd/gotree/go.go", []byte("XD")),
			Receiver: nil,
		},
		{
			Event:    create(n.W(), "src/github.com/rjeczalik/which/.which.go.swp"),
			Receiver: Chans{ch[1]},
		},
		{
			Event:    create(n.W(), "src/github.com/rjeczalik/which/.which.go.swo"),
			Receiver: Chans{ch[1]},
		},
		{
			Event:    remove(n.W(), "src/github.com/rjeczalik/fs/cmd/gotree/go.go"),
			Receiver: Chans{ch[2]},
		},
	}

	n.ExpectNotifyEvents(cases)
	n.ExpectDry(ch)

	// ...or using Call structures.
	stops := [...]Call{
		{
			F: FuncStop,
			C: ch[0],
		},
		{
			F: FuncStop,
			C: ch[1],
		},
	}

	n.Call(stops[:]...)

	cases = []NCase{
		{
			Event:    write(n.W(), "src/github.com/rjeczalik/fs/fs.go", []byte("XD")),
			Receiver: nil,
		},
		{
			Event:    write(n.W(), "src/github.com/rjeczalik/which/README.md", []byte("XD")),
			Receiver: nil,
		},
		{
			Event:    create(n.W(), "src/github.com/rjeczalik/which/.which.go.swr"),
			Receiver: nil,
		},
		{
			Event:    remove(n.W(), "src/github.com/rjeczalik/fs/cmd/gotree/main.go"),
			Receiver: Chans{ch[2]},
		},
	}

	n.ExpectNotifyEvents(cases)
	n.ExpectDry(ch)
}

func TestStop(t *testing.T) {
	t.Skip("TODO(rjeczalik)")
}
