package node

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/meta"
)

func TestJsonWriterListInList(t *testing.T) {
	moduleStr := `
module m {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	typedef td {
		type string;
	}
	list l1 {
		list l2 {
		    key "a";
			leaf a {
				type td;
			}
			leaf b {
			    type string;
			}
		}
	}
}
	`
	m := YangFromString(moduleStr)
	root := map[string]interface{}{
		"l1": []map[string]interface{}{
			map[string]interface{}{"l2": []map[string]interface{}{
				map[string]interface{}{
					"a": "hi",
					"b": "bye",
				},
			},
			},
		},
	}
	b := MapNode(root)
	var json bytes.Buffer
	sel := NewBrowser(m, b).Root()
	if err := sel.UpsertInto(NewJsonWriter(&json).Node()).LastErr; err != nil {
		t.Fatal(err)
	}
	actual := json.String()
	expected := `{"l1":[{"l2":[{"a":"hi","b":"bye"}]}]}`
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}

func TestJsonAnyData(t *testing.T) {
	tests := []struct {
		anything interface{}
		expected string
	}{
		{
			anything: map[string]interface{}{
				"a": "A",
				"b": "B",
			},
			expected: `"x":{"a":"A","b":"B"}`,
		},
		{
			anything: []interface{}{
				map[string]interface{}{
					"a": "A",
				},
				map[string]interface{}{
					"b": "B",
				},
			},
			expected: `"x":[{"a":"A"},{"b":"B"}]`,
		},
	}
	for _, test := range tests {
		var actual bytes.Buffer
		buf := bufio.NewWriter(&actual)
		w := &JsonWriter{
			out: buf,
		}
		m := meta.NewLeaf("x", "na")
		v := &Value{Type: meta.NewDataType(nil, "any"), AnyData: test.anything}
		w.writeValue(m, v)
		buf.Flush()
		if err := c2.CheckEqual(test.expected, actual.String()); err != nil {
			t.Error(err)
		}
	}
}
