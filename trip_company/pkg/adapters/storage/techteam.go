package storage

import (
	"context"
	"fmt"
	"strings"
	"tripcompanyservice/internal/techteam"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type techTeamRepo struct {
	db *gorm.DB
}

func NewTechTeamRepo(db *gorm.DB) techteam.Repo {
	return &techTeamRepo{db}
}

func (r *techTeamRepo) GetTechTeamByID(ctx context.Context, id uint) (*techteam.TechTeam, error) {
	var t entities.TechTeam
	if err := r.db.WithContext(ctx).
		Preload("Members").
		First(&t, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	dT := mappers.TechTeamEntityToDomain(t)
	return &dT, nil
}

func (r *techTeamRepo) Insert(ctx context.Context, t *techteam.TechTeam) error {
	teamEntity := mappers.TechTeamDomainToEntity(t)
	if err := r.db.WithContext(ctx).Save(&teamEntity).Error; err != nil {
		return err
	}

	t.ID = teamEntity.ID

	return nil
}

func (r *techTeamRepo) InsertMember(ctx context.Context, t *techteam.TechTeamMember) error {
	memberEntity := mappers.MemberDomainToEntity(t)
	if err := r.db.WithContext(ctx).Save(&memberEntity).Error; err != nil {
		return err
	}

	t.ID = memberEntity.ID
	return nil
}

func (r *techTeamRepo) GetTechTeamsOfCompany(ctx context.Context, companyId uint, limit, offset uint) ([]techteam.TechTeam, uint, error) {
	query := r.db.WithContext(ctx).
		Preload("Members").
		Model(&entities.TechTeam{}).
		Where("transport_company_id = ?", companyId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tech teams: %w", err)
	}

	if offset > 0 {
		query = query.Offset(int(offset))
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	var teams []entities.TechTeam
	if err := query.Find(&teams).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, 0, techteam.ErrFailedToFetchRecords
		}
		return nil, 0, fmt.Errorf("failed to fetch tech teams: %w", err)
	}

	return mappers.BatchTecTeamEnToDo(teams), uint(total), nil
}

func (r *techTeamRepo) GetTechTeamMemberByUserIDAndTechTeamID(ctx context.Context, userID uint, techTeamID uint) (*techteam.TechTeamMember, error) {
	var techTeamMember entities.TechTeamMember

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND tech_team_id = ?", userID, techTeamID).
		First(&techTeamMember).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving TechTeamMember: %w", err)
	}

	m := mappers.MemberEntityToMemberDomain(techTeamMember)
	return &m, nil
}

func (r *techTeamRepo) IsUserTechnicianInCompany(ctx context.Context, companyID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.TechTeamMember{}).
		Joins("JOIN tech_teams ON tech_team_members.tech_team_id = tech_teams.id").
		Where("tech_teams.transport_company_id = ? AND tech_team_members.user_id = ?", companyID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *techTeamRepo) Delete(ctx context.Context, tID uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.TechTeam{}, tID).Error; err != nil {
		return fmt.Errorf("%w %w", techteam.ErrDeleteTeam, err)
	} else {
		return nil
	}

}

func (r *techTeamRepo) GetFullTechTeamByID(ctx context.Context, id uint) (*techteam.TechTeam, error) {
	var t entities.TechTeam
	if err := r.db.WithContext(ctx).
		Preload("TransportCompany").
		Preload("Members").
		First(&t, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	teamD := mappers.FullTechTeamEntityToDomain(t)
	return &teamD, nil
}
