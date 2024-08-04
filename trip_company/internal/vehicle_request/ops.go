package vehiclerequest

import "context"

type Ops struct {
	repo Repo
	//penaltyRepo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Create(ctx context.Context, vR *VehicleRequest) error {
	return o.repo.Insert(ctx, vR)
}

func (o *Ops) UpdateVehicleReq(ctx context.Context, id uint, updates map[string]interface{}) error {
	return o.repo.UpdateVehicleReq(ctx, id, updates)
}
