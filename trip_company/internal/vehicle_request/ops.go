package vehiclerequest

import (
	"context"
	"tripcompanyservice/internal/company"
)

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

func (o *Ops) Delete(ctx context.Context, vRID uint) error {
	return o.repo.Delete(ctx, vRID)
}

func(o *Ops)GetVehicleReqByID(ctx context.Context, id uint) (*VehicleRequest, error) {
	p, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, ErrNotFound
	}
	return p, nil
}

func (o *Ops)GetTransportCompanyByVehicleRequestID(ctx context.Context, vehicleRequestID uint) (*company.TransportCompany, error) {
	return o.repo.GetTransportCompanyByVehicleRequestID(ctx, vehicleRequestID)
}