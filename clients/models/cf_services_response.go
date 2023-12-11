package models

type CloudFoundryServiceInstance struct {
	Guid          string        `json:"guid"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Tags          []string      `json:"tags"`
	LastOperation LastOperation `json:"last_operation,omitempty"`
	PlanGuid      string        `jsonry:"relationships.service_plan.data.guid"`
	SpaceGuid     string        `jsonry:"relationships.space.data.guid"`
	Metadata      Metadata      `json:"metadata"`

	Plan     ServicePlan     `json:"-"`
	Offering ServiceOffering `json:"-"`
}

type LastOperation struct {
	Type        string `json:"type"`
	State       string `json:"state"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ServicePlan struct {
	Guid         string `json:"guid"`
	Name         string `json:"name"`
	OfferingGuid string `jsonry:"relationships.service_offering.data.guid,omitempty"`
}

type ServiceOffering struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type ServiceBinding struct {
	Guid    string `json:"guid"`
	Name    string `json:"name,omitempty"`
	AppGuid string `jsonry:"relationships.app.data.guid,omitempty"`

	AppName string `json:"-"`
}

type ServiceInstanceAuxiliaryContent struct {
	ServicePlans     []ServicePlan     `json:"service_plans"`
	ServiceOfferings []ServiceOffering `json:"service_offerings"`
}

type ServiceBindingAuxiliaryContent struct {
	Apps []CloudFoundryApplication `json:"apps"`
}
