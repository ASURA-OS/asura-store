package main

import (
	"strings"
	"testing"
)

func TestParseWindowManager(t *testing.T) {
	cases := map[string]string{
		"Provider: hyprland\nLive: true\n": "hyprland",
		"Provider: (none)\n":               "",
		"Capabilities: unavailable\n":      "",
		"Provider: niri\n":                 "niri",
	}
	for in, want := range cases {
		if got := parseWindowManager(in); got != want {
			t.Fatalf("parseWindowManager(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestWindowManagerGate(t *testing.T) {
	cases := []struct {
		required, running string
		ok                bool
	}{
		{"", "hyprland", true},           // declares nothing: runs anywhere
		{"hyprland", "", true},           // provider unreadable: gate nothing
		{"hyprland", "hyprland", true},   // same manager
		{"Hyprland", "hyprland", true},   // names compare case-insensitively
		{" hyprland ", "hyprland", true}, // whitespace in a manifest
		{"niri", "hyprland", false},      // made for another manager
	}
	for _, c := range cases {
		ok, running := windowManagerGate(c.required, c.running)
		if ok != c.ok {
			t.Errorf("windowManagerGate(%q, %q) ok = %v, want %v", c.required, c.running, ok, c.ok)
		}
		if !ok && running != c.running {
			t.Errorf("windowManagerGate(%q, %q) running = %q, want %q", c.required, c.running, running, c.running)
		}
	}
}

// install refuses an item made for another window manager, and lets through
// one made for this one or for none.
func TestInstallRespectsWindowManagerGate(t *testing.T) {
	saved := runningWindowManager
	runningWindowManager = func() string { return "hyprland" }
	defer func() { runningWindowManager = saved }()

	p := fakeProvider{
		category: Category{ID: "themes", Name: "Themes", Group: "wear"},
		items: []Item{
			{ID: "niri-only", Name: "Niri Theme", Category: "themes", RequiredWindowManager: "niri"},
			{ID: "for-hypr", Name: "Hypr Theme", Category: "themes", RequiredWindowManager: "Hyprland"},
			{ID: "anywhere", Name: "Any Theme", Category: "themes"},
		},
	}
	if err := runInstall([]Provider{p}, []string{"themes", "niri-only"}); err == nil {
		t.Fatal("installing an item made for niri on hyprland was allowed")
	} else if !strings.Contains(err.Error(), "made for niri") {
		t.Fatalf("unexpected refusal: %v", err)
	}
	for _, id := range []string{"for-hypr", "anywhere"} {
		if err := runInstall([]Provider{p}, []string{"themes", id}); err != nil {
			t.Fatalf("installing %s: %v", id, err)
		}
	}
}
