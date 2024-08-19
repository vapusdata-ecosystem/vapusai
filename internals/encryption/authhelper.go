package encrytion

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/vapusdata-ecosystem/vapusai-studio/internals/models"
)

const (
	ClaimScopeKey = "scope"
	ClaimRoleKey  = "role"
	ClaimUserKey  = "userID"
)

func ParseUnValidatedJWT(tokenString string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	newParser := jwt.NewParser()
	_, _, err := newParser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func FlatJwtClaims(claims *models.AIStudioSpaceClaims, separator string) map[string]string {
	if err := claims.Scope.Validate(); err != nil {
		encryptLogger.Err(err).Msg("error while validating vapusdata platform access claims")
		return nil
	}
	if claims != nil {
		if claims.Scope.Space == "" || claims.Scope.Role == "" || claims.Scope.UserId == "" {
			return nil
		}
		return map[string]string{ClaimScopeKey: claims.Scope.Space, ClaimRoleKey: claims.Scope.Role, ClaimUserKey: claims.Scope.UserId}
	}
	return nil
}
