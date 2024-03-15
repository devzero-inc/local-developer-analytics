package collector

import "testing"

// TestParseCommand tests the ParseCommand function with various command inputs.
func TestParseCommand(t *testing.T) {

	testCases := []struct {
		command  string
		expected string
	}{
		{"./myapp", "myapp"},
		{"/usr/bin/myapp", "myapp"},
		{"sudo myapp", "myapp"},
		{"nohup ./myapp", "myapp"},
		{"kubectl describe pod", "kubectl"},
		{"docker run image", "docker"},
		{"git pull origin master", "git"},
		{"/bin/ls -l", "ls"},
		{"nohup /usr/local/bin/myapp", "myapp"},
	}

	for _, tc := range testCases {
		result := ParseCommand(tc.command)

		if result != tc.expected {
			t.Errorf("ParseCommand(%q) = %q, expected %q", tc.command, result, tc.expected)
		}
	}
}
