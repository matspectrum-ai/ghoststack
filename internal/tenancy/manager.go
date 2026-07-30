package tenancy

import "context"

type Tenant struct {
	ID   string
	Name string
}

type TenantManager struct {
	tenants map[string]Tenant
}

func NewTenantManager() *TenantManager {
	return &TenantManager{tenants: make(map[string]Tenant)}
}

func (m *TenantManager) Register(ctx context.Context, tenant Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *TenantManager) Get(ctx context.Context, id string) (Tenant, error) {
	tenant, ok := m.tenants[id]
	if !ok {
		return Tenant{}, ErrTenantNotFound
	}
	return tenant, nil
}

var ErrTenantNotFound = &tenancyError{"tenant not found"}

type tenancyError struct{ msg string }

func (e *tenancyError) Error() string { return e.msg }
