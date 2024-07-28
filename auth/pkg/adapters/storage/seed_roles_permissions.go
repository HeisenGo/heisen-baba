package storage

import (
    "authservice/pkg/adapters/storage/entities"
    "gorm.io/gorm"
)

func seedRolesAndPermissions(db *gorm.DB) error {
    // Define initial roles
    roles := []entities.Role{
        {Name: "admin", Description: "System administrator"},
        {Name: "employee", Description: "Regular employee"},
        {Name: "manager", Description: "Company manager"},
    }

    // Define initial permissions
    permissions := []entities.Permission{
        {Name: "view_company", Description: "Can view company details"},
        {Name: "edit_company", Description: "Can edit company details"},
        {Name: "view_trip", Description: "Can view trip details"},
        {Name: "edit_trip", Description: "Can edit trip details"},
    }

    // Upsert roles
    for _, role := range roles {
        err := db.Where(entities.Role{Name: role.Name}).FirstOrCreate(&role).Error
        if err != nil {
            return err
        }
    }

    // Upsert permissions
    for _, permission := range permissions {
        err := db.Where(entities.Permission{Name: permission.Name}).FirstOrCreate(&permission).Error
        if err != nil {
            return err
        }
    }

    // Assign permissions to roles
    adminRole := entities.Role{}
    db.Where(entities.Role{Name: "admin"}).First(&adminRole)

    employeeRole := entities.Role{}
    db.Where(entities.Role{Name: "employee"}).First(&employeeRole)

    managerRole := entities.Role{}
    db.Where(entities.Role{Name: "manager"}).First(&managerRole)

    for _, permission := range permissions {
        db.Where(entities.Permission{Name: permission.Name}).First(&permission)

        // Assign all permissions to admin
        db.FirstOrCreate(&entities.RolePermission{RoleID: adminRole.ID, PermissionID: permission.ID})

        // Assign view permissions to employee
        if permission.Name == "view_company" || permission.Name == "view_trip" {
            db.FirstOrCreate(&entities.RolePermission{RoleID: employeeRole.ID, PermissionID: permission.ID})
        }

        // Assign all permissions except edit_company to manager
        if permission.Name != "edit_company" {
            db.FirstOrCreate(&entities.RolePermission{RoleID: managerRole.ID, PermissionID: permission.ID})
        }
    }

    return nil
}