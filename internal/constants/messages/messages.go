package messages

const (
	InvalidToken = "incorrect_token"

	ErrorUnauthorized   = "unauthorized"
	ErrorInvalidRequest = "invalid_request"
	ErrorNotFound       = "not_found"
	ErrorForbidden      = "forbidden"
	ErrorAlreadyExists = "already_exists"

	MessageAuthorizationRequired = "Authorization header required"
	MessageInvalidAuthorization  = "Invalid or expired token"

	MessageTitleRequired         = "title is required"
	MessageTitleTooLong          = "title must not exceed %d characters provided %d"
	MessageDescriptionRequired   = "description is required"
	MessageDescriptionLength     = "description must be between 1 and %d characters"
	MessagePriceRange            = "price must be between 0 and 1000000 inclusive"
	MessagePriceNotNumber        = "cannot convert string value %s to float"
	MessageQuantityRange         = "quantity must be between 1 and 100 inclusive"
	MessagePhotoRequired         = "at least 1 photo is required"
	MessageTooManyPhotos         = "maximum %d photos allowed"
	MessageAdvertisementNotFound = "Advertisement not found"
	MessageFavoriteNotFound      = "Favorite entry not found"
	MessageFavoriteAlreadyExists = "Advertisement is already in favorites"
	MessageFavoriteOwnForbidden = "Cannot add own advertisement to favorites"
	MessageInvalidUUIDLength    = "invalid UUID length: %d"
)
