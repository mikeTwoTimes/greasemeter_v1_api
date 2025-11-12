package utility

import (
	"errors"
	"net/mail"
	"regexp"
	"strconv"

	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
)

var validASCII = regexp.MustCompile(`^[\x21-\x7E]+$`)

func ParsePagination(c *gin.Context) (types.Pagination, error) {
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil || page <= 0 {
		return types.Pagination{}, errors.New("Invalid page")
	}

	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil || limit <= 0 || limit > 20 {
		return types.Pagination{}, errors.New("Invalid limit")
	}

	return types.Pagination{
		Offset: page,
		Limit:  limit,
	}, nil
}

func ParseCoordinates(c *gin.Context) (float64, float64, error) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)

	if err != nil || lat < -90 || lat > 90 {
		return 0.0, 0.0, errors.New("Invalid latitude")
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)

	if err != nil || lng < -180 || lng > 180 {
		return 0.0, 0.0, errors.New("Invalid longitude")
	}

	return lat, lng, nil
}

func ParseBoundingBox(c *gin.Context) (types.Bounds, error) {
	lat, lng, err := ParseCoordinates(c)

	if err != nil {
		return types.Bounds{}, err
	}

	latDelta, err := strconv.ParseFloat(c.Query("latDelta"), 64)

	if err != nil || latDelta <= 0 || latDelta > 0.25 {
		return types.Bounds{}, errors.New("Invalid latitude delta")
	}

	lngDelta, err := strconv.ParseFloat(c.Query("lngDelta"), 64)

	if err != nil || lngDelta <= 0 || lngDelta > 0.45 {
		return types.Bounds{}, errors.New("Invalid longitude delta")
	}

	return types.Bounds{
		LatMin: lat - latDelta/2,
		LatMax: lat + latDelta/2,
		LngMin: lng - lngDelta/2,
		LngMax: lng + lngDelta/2,
	}, nil
}

func ParseReview(c *gin.Context) (types.ReviewPayload, error) {
	var req types.ReviewPayload

	if c.ShouldBindJSON(&req) != nil {
		return types.ReviewPayload{}, errors.New("Failed to bind review")
	} else if len([]rune(req.Text)) < 1 ||
		len([]rune(req.Text)) > 200 {
		return types.ReviewPayload{}, errors.New(
			"Reviews must be between 1 and 200 characters",
		)
	} else if req.Rating < 1 || req.Rating > 5 {
		return types.ReviewPayload{}, errors.New(
			"Ratings must be between 1 and 5",
		)
	} else if goaway.IsProfane(req.Text) {
		return types.ReviewPayload{}, errors.New(
			"Reviews contains profanity",
		)
	}

	return req, nil
}

func ParseRecommendation(c *gin.Context) (types.RecommendationPayload, error) {
	var req types.RecommendationPayload

	if c.ShouldBindJSON(&req) != nil {
		return types.RecommendationPayload{}, errors.New(
			"Failed to bind recommendation",
		)
	} else if len(req.Name) < 1 ||
		len(req.Name) > 255 {
		return types.RecommendationPayload{}, errors.New(
			"Place names must be between 1 and 255 characters",
		)
	} else if len(req.Address) < 1 ||
		len(req.Address) > 255 {
		return types.RecommendationPayload{}, errors.New(
			"Place addresses must be between 1 and 255 characters",
		)
	}

	return req, nil
}

func ParseReport(c *gin.Context) (string, error) {
	var req types.ReportPayload

	if c.ShouldBindJSON(&req) != nil {
		return "", errors.New("Failed to bind report")
	} else if len(req.Reason) < 1 ||
		len(req.Reason) > 255 {
		return "", errors.New(
			"Report reasons must be between 1 and 255 characters",
		)
	}

	return req.Reason, nil
}

func ParseForgotPassword(c *gin.Context) (string, error) {
	var req types.ForgotPasswordPayload
	err := c.ShouldBindJSON(&req)

	if err != nil {
		return "", errors.New("Failed to bind email")
	} else if _, err = mail.ParseAddress(req.Email); err != nil {
		return "", errors.New("Invalid email address")
	}

	return req.Email, nil
}

func ParseResetPassword(c *gin.Context) (string, error) {
	var req types.ResetPasswordPayload
	err := c.ShouldBindJSON(&req)

	if err != nil {
		return "", errors.New("Failed to bind password")
	} else if len(req.Password) < 12 || len(req.Password) > 255 {
		return "", errors.New(
			"Passwords must be between 12 and 255 characters",
		)
	} else if !validASCII.MatchString(req.Password) {
		return "", errors.New("Password contains illegal character")
	}

	return req.Password, nil
}

func ParseRegister(c *gin.Context) (types.RegisterPayload, error) {
	var req types.RegisterPayload

	if c.ShouldBindJSON(&req) != nil {
		return types.RegisterPayload{}, errors.New("Failed to bind user")
	} else if len(req.Name) < 6 || len(req.Name) > 30 {
		return types.RegisterPayload{}, errors.New(
			"Usernames must be between 6 and 30 characters",
		)
	} else if !validASCII.MatchString(req.Name) {
		return types.RegisterPayload{}, errors.New(
			"Username contains illegal character",
		)
	} else if goaway.IsProfane(req.Name) {
		return types.RegisterPayload{}, errors.New(
			"Username contains profanity",
		)
	} else if len(req.Password) < 12 || len(req.Password) > 255 {
		return types.RegisterPayload{}, errors.New(
			"Passwords must be between 12 and 255 characters",
		)
	} else if !validASCII.MatchString(req.Password) {
		return types.RegisterPayload{}, errors.New(
			"Password contains illegal character",
		)
	} else if _, err := mail.ParseAddress(req.Email); err != nil {
		return types.RegisterPayload{}, errors.New(
			"Invalid email address",
		)
	}

	return req, nil
}

func ParseLogin(c *gin.Context) (types.LoginPayload, error) {
	var req types.LoginPayload

	if c.ShouldBindJSON(&req) != nil {
		return req, errors.New("Failed to bind user")
	} else if len(req.Name) < 6 || len(req.Name) > 30 ||
		!validASCII.MatchString(req.Name) ||
		goaway.IsProfane(req.Name) ||
		len(req.Password) < 12 || len(req.Password) > 255 ||
		!validASCII.MatchString(req.Password) {
		return types.LoginPayload{}, errors.New(
			"Invalid username or password",
		)
	}

	return req, nil
}
