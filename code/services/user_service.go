package services

import (
	"database/sql"
	"errors"
	"login/database"
	"login/models"
	"login/utils"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(req *models.UserRegisterRequest) (*models.UserResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.checkUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 对密码进行哈希
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 插入用户数据
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := database.DB.Exec(query, req.Username, hashedPassword)
	if err != nil {
		return nil, err
	}

	// 获取插入的用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	user, err := s.GetUserByID(int(userID))
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *models.UserLoginRequest) (*models.LoginResponse, error) {
	// 根据用户名查询用户
	user, err := s.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, "your-super-secret-jwt-key-change-in-production", 3600)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID int) (*models.User, error) {
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE id = ?`

	var user models.User
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username = ?`

	var user models.User
	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// checkUsernameExists 检查用户名是否存在
func (s *UserService) checkUsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	var count int
	err := database.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
