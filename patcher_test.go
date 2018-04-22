package json_patcher

import (
	"testing"
	"reflect"
	"errors"
	"github.com/gabriel-araujjo/json-patcher/mock"
)

func TestNewPatch(t *testing.T) {
	cases := []struct{
		name string
		patchBuf []byte
		expectErr bool
		expectPatch Patch
	}{{
		name: "Add",
		patchBuf: []byte(`[{"op": "add", "path": "/foo", "value": "bar"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"add", Path:"/foo", Value: "bar"}),
	},{
		name: "Replace",
		patchBuf: []byte(`[{"op": "replace", "path": "/foo", "value": "bar"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"replace", Path:"/foo", Value: "bar"}),
	},{
		name: "Remove",
		patchBuf: []byte(`[{"op": "remove", "path": "/foo"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"remove", Path:"/foo"}),
	},{
		name: "Move",
		patchBuf: []byte(`[{"op": "move", "from": "/foo/1", "path": "/foo/2"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"move", From:"/foo/1", Path:"/foo/2"}),
	},{
		name: "Copy",
		patchBuf: []byte(`[{"op": "copy", "from": "/baz", "path": "/foo"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"copy", From:"/baz", Path:"/foo"}),
	},{
		name: "Test",
		patchBuf: []byte(`[{"op": "test", "from": "/baz", "value": "baz"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"test", From:"/baz", Value:"baz"}),
	},{
		name: "ExtraKeys",
		patchBuf: []byte(`[{"op": "test", "from": "/baz", "value": "baz", "foo": "foo"}]`),
		expectErr: false,
		expectPatch: append(Patch{}, Operation{Op:"test", From:"/baz", Value:"baz"}),
	},{
		name: "OnlyObj",
		patchBuf: []byte(`{"op": "test", "from": "/baz", "value": "baz"}`),
		expectErr: true,
		expectPatch: nil,
	},{
		name: "MalFormatted",
		patchBuf: []byte(`{"op": "test",, "from": "/baz", "value": "baz"}`),
		expectErr: true,
		expectPatch: nil,
	}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			patch, err := NewPatch(tt.patchBuf)
			if tt.expectErr {
				if err == nil || patch != nil {
					t.Errorf("expecting err, but nil was returned instead")
				}
			} else {
				if err != nil || patch == nil {
					t.Errorf("not expecting err %q", err)
				}
				if !reflect.DeepEqual(tt.expectPatch, patch) {
					t.Errorf("Expecting %#v instead of %#v", tt.expectPatch, patch)
				}
			}
		})
	}
}

var someErr = errors.New("some error")

func TestMend(t *testing.T) {
	cases := []struct{
		name  string
		patch Patch
		input mock.Tailor
		want  mock.Tailor
		expectErr bool
	}{{
		name: "Add",
		patch: Patch{{
			Op:    "add",
		}},
		want: mock.Tailor{
			AddCalled: true,
		},
		expectErr: false,
	},{
		name: "Replace",
		patch: Patch{{
			Op:    "replace",
		}},
		want: mock.Tailor{
			ReplaceCalled: true,
		},
		expectErr: false,
	},{
		name: "Remove",
		patch: Patch{{
			Op:    "remove",
		}},
		want: mock.Tailor{
			RemoveCalled: true,
		},
		expectErr: false,
	},{
		name: "Move",
		patch: Patch{{
			Op:    "move",
			From: "/foo/1",
			Path: "/foo/2",
		}},
		want: mock.Tailor{
			MoveCalled: true,
		},
		expectErr: false,
	},{
		name: "AddError",
		patch: Patch{{
			Op:    "add",
		}},
		input: mock.Tailor{AddReturn:someErr},
		want: mock.Tailor{
			AddCalled: true,
			AddReturn: someErr,
		},
		expectErr: true,
	},{
		name: "ReplaceError",
		patch: Patch{{
			Op:    "replace",
		}},
		input: mock.Tailor{ReplaceReturn:someErr},
		want: mock.Tailor{
			ReplaceCalled: true,
			ReplaceReturn: someErr,
		},
		expectErr: true,
	},{
		name: "RemoveError",
		patch: Patch{{
			Op:    "remove",
		}},
		input: mock.Tailor{RemoveReturn:someErr},
		want: mock.Tailor{
			RemoveCalled: true,
			RemoveReturn: someErr,
		},
		expectErr: true,
	},{
		name: "MoveError",
		patch: Patch{{
			Op:    "move",
			From: "/foo/1",
			Path: "/foo/2",
		}},
		input: mock.Tailor{MoveReturn:someErr},
		want: mock.Tailor{
			MoveCalled: true,
			MoveReturn: someErr,
		},
		expectErr: true,
	},{
		name: "InvalidMovePaths",
		patch: Patch{{
			Op:    "move",
			From: "/foo/1",
			Path: "/bar/2",
		}},
		expectErr: true,
	},{
		name: "InvalidMovePaths#1",
		patch: Patch{{
			Op:    "move",
			From: "foo",
			Path: "bar/2",
		}},
		expectErr: true,
	},{
		name: "InvalidMovePaths#2",
		patch: Patch{{
			Op:    "move",
			From: "/foo/",
			Path: "bar/2",
		}},
		expectErr: true,
	},{
		name: "InvalidMovePaths#3",
		patch: Patch{{
			Op:    "move",
			From: "/foo/",
			Path: "b",
		}},
		expectErr: true,
	},{
		name: "InvalidMovePaths#4",
		patch: Patch{{
			Op:    "move",
			From: "/foo/x",
			Path: "/foo/1",
		}},
		expectErr: true,
	},{
		name: "InvalidMovePaths#5",
		patch: Patch{{
			Op:    "move",
			From: "/foo/1",
			Path: "/foo/x",
		}},
		expectErr: true,
	},{
		name: "InvalidOperation",
		patch: Patch{{
			Op: "foo",
		}},
		expectErr: true,
	}}

	dummy := "dummy"

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := Mend(&tt.input, tt.patch, &dummy)
			if tt.expectErr {
				if err == nil {
					t.Errorf("test %q should return error. %#v", tt.name, tt.input)
				}
			}

			if !reflect.DeepEqual(tt.want, tt.input) {
				t.Errorf("expecting %#v instead of %#v", tt.want, tt.input)
			}
		})
	}

	t.Run("NilTailor", func(t *testing.T) {
		var p Patch
		err := Mend(nil, p, dummy)
		if err == nil {
			t.Errorf("error  %e was not expected", err)
		}
	})

	t.Run("NilObj", func(t *testing.T) {
		var p Patch
		err := Mend(&mock.Tailor{}, p, nil)
		if err == nil {
			t.Errorf("error  %e was not expected", err)
		}
	})
}
