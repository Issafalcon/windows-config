package module

import "testing"

func TestGetInstallOrderIncludesTransitiveDependencies(t *testing.T) {
	r := NewRegistry()
	for _, m := range []Module{
		{Name: "git"},
		{Name: "neovim", Dependencies: []string{"git"}},
		{Name: "lazygit", Dependencies: []string{"neovim"}},
	} {
		if err := r.Register(m); err != nil {
			t.Fatal(err)
		}
	}

	order, err := r.GetInstallOrder("lazygit")
	if err != nil {
		t.Fatal(err)
	}
	got := []string{order[0].Name, order[1].Name, order[2].Name}
	want := []string{"git", "neovim", "lazygit"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("install order = %v, want %v", got, want)
		}
	}
}
