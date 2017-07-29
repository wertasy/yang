package yang

import (
	"testing"

	"github.com/c2stack/c2g/meta"
)

func TestTypeResolve(t *testing.T) {
	yang := `
module ff {
	namespace "ns";

	description "mod";

	revision 99-99-9999 {
	  description "bingo";
	}

	leaf x {
		type int32;
	}
	typedef q {
		type string;
	}
	list y {
		key "id";
		leaf id {
			type string;
		}
	}
	container z {
	  description "z1";
	  leaf z1 {
	    type leafref {
	    	path "../x";
	    }
	  }
	  leaf z2 {
	    type leafref {
	    	path "../y/id";
	    }
	  }
		leaf z3 {
			type q;
		}
	}
}
`
	m, err := LoadModuleCustomImport(yang, nil)
	if err != nil {
		t.Fatal(err)
	}
	z1, err := meta.FindByPath(m, "z/z1")
	if err != nil {
		t.Error(err)
	} else if z1 == nil {
		t.Errorf("No z1")
	}
	dt := z1.(meta.HasDataType).GetDataType()
	if f, err := dt.Format(); err != nil {
		t.Error(err)
	} else if f != meta.FMT_INT32 {
		t.Errorf("actual type %d", f)
	}
	z3, err := meta.FindByPath(m, "z/z3")
	if err != nil {
		t.Error(err)
	}
	dt = z3.(meta.HasDataType).GetDataType()
	if f, err := dt.Format(); err != nil {
		t.Error(err)
	} else if f != meta.FMT_STRING {
		t.Errorf("actual type %d", f)
	}
}
