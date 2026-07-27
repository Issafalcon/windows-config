package module

import (
	"fmt"
	"sort"
)

type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry {
	return &Registry{modules: make(map[string]Module)}
}

var DefaultRegistry = NewRegistry()

func (r *Registry) Register(m Module) error {
	if m.Name == "" {
		return fmt.Errorf("module name is required")
	}
	if _, exists := r.modules[m.Name]; exists {
		return fmt.Errorf("module %q is already registered", m.Name)
	}
	r.modules[m.Name] = m
	return nil
}

func (r *Registry) Get(name string) (Module, bool) {
	m, ok := r.modules[name]
	return m, ok
}

func (r *Registry) All() []Module {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	sort.Strings(names)

	modules := make([]Module, 0, len(names))
	for _, name := range names {
		modules = append(modules, r.modules[name])
	}
	return modules
}

func (r *Registry) Categories() []string {
	set := make(map[string]struct{})
	for _, m := range r.modules {
		if m.Category != "" {
			set[m.Category] = struct{}{}
		}
	}

	categories := make([]string, 0, len(set))
	for category := range set {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

func (r *Registry) GetInstallOrder(names ...string) ([]Module, error) {
	if len(names) == 0 {
		for name := range r.modules {
			names = append(names, name)
		}
	}

	selected := make(map[string]struct{})
	var visit func(string) error
	visit = func(name string) error {
		if _, seen := selected[name]; seen {
			return nil
		}
		m, ok := r.modules[name]
		if !ok {
			return fmt.Errorf("module %q is not registered", name)
		}
		selected[name] = struct{}{}
		for _, dep := range m.Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		return nil
	}
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	indegree := make(map[string]int, len(selected))
	dependents := make(map[string][]string, len(selected))
	for name := range selected {
		indegree[name] = 0
	}
	for name := range selected {
		for _, dep := range r.modules[name].Dependencies {
			if _, included := selected[dep]; included {
				indegree[name]++
				dependents[dep] = append(dependents[dep], name)
			}
		}
	}

	ready := make([]string, 0, len(selected))
	for name, degree := range indegree {
		if degree == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)

	order := make([]Module, 0, len(selected))
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		order = append(order, r.modules[name])
		for _, dependent := range dependents[name] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.Strings(ready)
	}
	if len(order) != len(selected) {
		return nil, fmt.Errorf("module dependency cycle detected")
	}
	return order, nil
}
