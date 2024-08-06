package techteam

import (
	"context"
)

type Ops struct {
	repo Repo
	//penaltyRepo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CreateTechTeam(ctx context.Context, t *TechTeam) error {
	return o.repo.Insert(ctx, t)
}

func (o *Ops) CreateTechTeamMember(ctx context.Context, m *TechTeamMember) error {
	return o.repo.InsertMember(ctx, m)
}

func (o *Ops) GetTechTeamByID(ctx context.Context, id uint) (*TechTeam, error) {
	t, err := o.repo.GetTechTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrTeamNotFound
	}

	return t, nil
}

func (o *Ops) GetTechTeamsOfCompany(ctx context.Context, companyId uint, page, pageSize uint) ([]TechTeam, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	return o.repo.GetTechTeamsOfCompany(ctx, companyId, limit, offset)
}

func (o *Ops) GetTechTeamMemberByUserIDAndTechTeamID(ctx context.Context, userID uint, techTeamID uint) (*TechTeamMember, error) {
	m, err:= o.repo.GetTechTeamMemberByUserIDAndTechTeamID(ctx, userID, techTeamID)
	if err!=nil{
		return nil, err
	}
	if m == nil {
		return nil, ErrMemberNotFound
	}
	return m, nil
}

func (o *Ops)IsUserTechnicianInCompany(ctx context.Context, companyID uint, userID uint) (bool, error) {
	return o.repo.IsUserTechnicianInCompany(ctx, companyID, userID)
}