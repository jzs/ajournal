package profile

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/jzs/ajournal/blob"
	"github.com/jzs/ajournal/user"
	"github.com/jzs/ajournal/utils"
)

// Service describes the methods on a profile service
type Service interface {
	Create(ctx context.Context, p *Profile) (*Profile, error)
	Profile(ctx context.Context) (*Profile, error)
	ProfileByShortName(ctx context.Context, sn string) (*Profile, error)
	UserProfile(ctx context.Context, userid int64) (*Profile, error)
	UpdateProfile(ctx context.Context, p *Profile) (*Profile, error)
	Subscribe(ctx context.Context, sub *Subscription) error
	GenerateShortName(sn string) string
	ValidateShortName(ctx context.Context, userID int64, sn string) bool
	ChangePicture(ctx context.Context, userID int64, img io.Reader, filetype string) error
}

// NewService returns a new implementation of the Service interface
func NewService(pr Repository, sr SubscriptionRepository, bs blob.Service) Service {
	return &service{pr: pr, sr: sr, bs: bs}
}

type service struct {
	pr Repository
	sr SubscriptionRepository
	bs blob.Service
}

func (s *service) Create(ctx context.Context, p *Profile) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create profile without a user context")
	}
	if usr.ID != p.ID {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "Cannot create profile for another user")
	}

	if p.Plan == 0 {
		p.Plan = PlanFree
	}

	if p.ShortName == "" {
		p.ShortName = s.GenerateShortName(p.Email)
	}

	if !s.ValidateShortName(ctx, p.ID, p.ShortName) {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, fmt.Sprintf("Shortname '%v' is not a valid short name", p.ShortName))
	}

	prof, err := s.pr.Create(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "CreateProfile")
	}
	return prof, nil

}

func (s *service) UserProfile(ctx context.Context, userid int64) (*Profile, error) {
	pro, err := s.pr.FindByID(ctx, userid)
	if err == ErrProfileNotExist {
		return nil, errors.Wrap(err, "Profile doesn't exist")
	}
	pro.Email = ""
	return pro, nil
}

func (s *service) ProfileByShortName(ctx context.Context, sn string) (*Profile, error) {
	pro, err := s.pr.FindByShortName(ctx, sn)
	if err == ErrProfileNotExist {
		return nil, errors.Wrap(err, "Profile doesn't exist")
	}
	pro.Email = ""
	return pro, nil
}

func (s *service) Profile(ctx context.Context) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create a journal without a user context")
	}

	pro, err := s.pr.FindByID(ctx, usr.ID)
	if err != nil && err != ErrProfileNotExist {
		return nil, err
	}
	if err == ErrProfileNotExist {
		// Create profile and return that.
		pro, err = s.Create(ctx, &Profile{ID: usr.ID, Email: usr.Username, Plan: PlanFree})
		if err != nil {
			return nil, errors.Wrap(err, "Could not create profile for user")
		}
		return pro, nil
	}
	return pro, nil
}

func (s *service) UpdateProfile(ctx context.Context, p *Profile) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create a journal without a user context")
	}
	if usr.ID != p.ID {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "Cannot update another users profile")
	}

	if !s.ValidateShortName(ctx, p.ID, p.ShortName) {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, fmt.Sprintf("Shortname '%v' is not a valid short name", p.ShortName))
	}

	prof, err := s.pr.Update(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateProfile")
	}
	return prof, nil
}

// SubscriptionArgs args for signing up for a subscription
type SubscriptionArgs struct {
	CardName string
	Number   string
	Month    string
	Year     string
	CVC      string
}

func (s *service) Subscribe(ctx context.Context, sub *Subscription) error {
	_, err := s.sr.Create(ctx, sub)
	return err
}

var reg = regexp.MustCompile("[^a-z0-9]")
var regmatch = regexp.MustCompile("^[a-z0-9\\-]+$")

// GenerateShortName takes an input string, shortens it down and tries to generate a valid short name.
func (s *service) GenerateShortName(sn string) string {
	// Strip @ to end.
	sn = strings.Split(sn, "@")[0]
	sn = strings.ToLower(sn)
	sn = reg.ReplaceAllString(sn, "-")
	return sn
}

// ValidateShortName validates whether a short name is valid for storage.
func (s *service) ValidateShortName(ctx context.Context, userID int64, sn string) bool {
	match := regmatch.MatchString(sn)
	if !match {
		return false
	}
	// look up if it already exists in db...
	prof, err := s.pr.FindByShortName(ctx, sn)
	if err == ErrProfileNotExist {
		return true
	}
	if err != nil {
		return false
	}
	return prof.ID == userID
}

func (s *service) ChangePicture(ctx context.Context, userID int64, file io.Reader, header string) error {
	f, err := s.bs.Create(path.Join("users", fmt.Sprint(userID), "profile-pic"), header, file)
	if err != nil {
		return errors.Wrap(err, "ProfileService could not create blob")
	}

	pr, err := s.pr.FindByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "ProfileService could not find profile for user")
	}

	pr.Picture = blob.File{Key: f.Key}
	_, err = s.pr.Update(ctx, pr)
	if err != nil {
		return errors.Wrap(err, "ProfileService could not update profile for user")
	}
	return nil
}
