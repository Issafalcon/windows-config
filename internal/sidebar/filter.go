package sidebar

import "strings"

type ModuleItem struct {
	Name, Icon, Description, Category string
	Installed                         bool
}

func FuzzyMatch(query, target string) bool {
	return strings.Contains(strings.ToLower(target), strings.ToLower(query))
}

func FilterModules(items []ModuleItem, query string) []ModuleItem {
	if strings.TrimSpace(query) == "" {
		return items
	}
	var filtered []ModuleItem
	for _, item := range items {
		if FuzzyMatch(query, item.Name) || FuzzyMatch(query, item.Description) || FuzzyMatch(query, item.Category) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func FilterByCategory(items []ModuleItem, category string) []ModuleItem {
	if category == "" {
		return items
	}
	var filtered []ModuleItem
	for _, item := range items {
		if item.Category == category {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
