package community

import "context"

type Template struct {
	ID    string
	Name  string
	YAML  string
}

type Gallery struct {
	templates map[string]Template
}

func NewGallery() *Gallery {
	return &Gallery{templates: make(map[string]Template)}
}

func (g *Gallery) Publish(ctx context.Context, tpl Template) error {
	g.templates[tpl.ID] = tpl
	return nil
}

func (g *Gallery) List(ctx context.Context) ([]Template, error) {
	out := make([]Template, 0, len(g.templates))
	for _, tpl := range g.templates {
		out = append(out, tpl)
	}
	return out, nil
}
