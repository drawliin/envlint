package scanner

import (
	"testing"
)

func TestScanDetectsPythonEnvironmentIndex(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config.py", "import os\nfirst = os.environ[\"API_KEY\"]\nsecond = os.environ['OTHER_KEY']\n")
	result, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"API_KEY", "OTHER_KEY"} {
		if len(result.Referenced[key]) != 1 {
			t.Fatalf("expected %s once, got %v", key, result.Referenced)
		}
	}
}
