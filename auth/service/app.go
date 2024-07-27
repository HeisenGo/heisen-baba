package service

import (
	"authservice/config"
	"authservice/internal/user"
	"authservice/pkg/adapters/storage"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg                 config.Config
	dbConn              *gorm.DB
	authService         *AuthService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()

	app.setAuthService()
	
	return app, nil
}

func (a *AppContainer) RawDBConnection() *gorm.DB {
	return a.dbConn
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db

	err = storage.AddExtension(a.dbConn)
	if err != nil {
		log.Fatal("Create extension failed: ", err)
	}

	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

func (a *AppContainer) AuthService() *AuthService {
	return a.authService
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn)), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}

// func (a *AppContainer) BoardService() *BoardService {
// 	return a.boardService
// }

// func (a *AppContainer) BoardServiceFromCtx(ctx context.Context) *BoardService {
// 	tx, ok := valuecontext.TryGetTxFromContext(ctx)
// 	if !ok {
// 		return a.boardService
// 	}

// 	gc, ok := tx.Tx().(*gorm.DB)
// 	if !ok {
// 		return a.boardService
// 	}

// 	return NewBoardService(
// 		user.NewOps(storage.NewUserRepo(gc)),
// 		board.NewOps(storage.NewBoardRepo(gc)),
// 		userboardrole.NewOps(storage.NewUserBoardRepo(gc)),
// 		column.NewOps(storage.NewColumnRepo(gc)),
// 		notification.NewOps(storage.NewNotificationRepo(gc)),
// 	)
// }


// func (a *AppContainer) setBoardService() {
// 	if a.boardService != nil {
// 		return
// 	}
// 	a.boardService = NewBoardService(user.NewOps(storage.NewUserRepo(a.dbConn)), board.NewOps(storage.NewBoardRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)), column.NewOps(storage.NewColumnRepo(a.dbConn)), notification.NewOps(storage.NewNotificationRepo(a.dbConn)))
// }

