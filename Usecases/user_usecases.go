package usecases

// imports
import (
	"errors";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Domain";
)

type UserUseCase struct {
	UserRepository     domain.UserRepository
	JWTService         domain.JWTService
	PwdService         domain.PasswordService
}

// creates new UserUseCase instance
func NewUserUseCase(userRepo domain.UserRepository, jwtServ domain.JWTService, pwdServ domain.PasswordService,) domain.UserUseCase {
	return &UserUseCase{ UserRepository:userRepo, JWTService:jwtServ, PwdService:pwdServ}
}

// register user
func (userUsc *UserUseCase) Register(user *domain.User) error {
	
	// validate input
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// check if user already exists
	existing, err := userUsc.UserRepository.GetByUsername(user.Username)
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}
	if existing != nil {
		return domain.ErrUserExists
	}

	// hash password securely 
	hashed, err := userUsc.PwdService.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed       // set user password to hashed password

	// set default role
	user.Role = "user"

	// first user becomes admin
	count, err := userUsc.UserRepository.GetUserCount()
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = "admin"
	}

	return userUsc.UserRepository.CreateUser(user)
}

// authenticate user
func (userUsc *UserUseCase) Login(credentials *domain.Credentials) (string, *domain.User, error) {
	
	// validate input
	if credentials.Username == "" || credentials.Password == "" {
		return "", nil, errors.New("username and password are required")
	}

	// get user from repository
	user, err := userUsc.UserRepository.GetByUsername(credentials.Username)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, err
	}

	// verify password
	if !userUsc.PwdService.CheckPassword(user.Password, credentials.Password) {
		return "", nil, domain.ErrInvalidCredentials
	}

	// generate jwt token
	token, err := userUsc.JWTService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	// return token and user (without sensitive data)
	returnUser := &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	return token, returnUser, nil
}

// promote a user to admin role (only admin can do this)
func (userUsc *UserUseCase) PromoteToAdmin(userID string) error {
	
	// validate input
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	// check if user exists
	_, err := userUsc.UserRepository.GetUserById(userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrUserNotFound
		}
		return err
	}

	// update role
	return userUsc.UserRepository.UpdateRole(userID, "admin")
}