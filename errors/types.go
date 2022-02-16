package errors

import "github.com/gofiber/fiber/v2"

// ERROR CODE CONVENTION, BUT JUST CONVENTION
// 400|[SERVICE_ID]|[ERROR_CODE]
// 3 DIGITS HTTP RESP STATUS CODE | 2 DIGITS SERVICE ID | 2 DIGITS ERROR CODE
//
// 10 -> authentication service
// 20 -> social authentication service
// 30 -> oauth service
// 40 -> verification service
// 50 -> token service
// 60 -> account service
// 70 -> region service
// 80 -> template service
// 90 -> websso service
// 11 -> recovery info service
// 12 -> account verification service
// 13 -> consent service
// 14 -> multi factor auth service
// 15 -> repo service
// 16 -> authentication gateway
// 17 -> anomaly detection service
// 18 -> ethereum auth service

var NONCE_EXPIRED_ERROR *fiber.Error = fiber.NewError(
	4001801,
	"nonce is expired",
)

var INVALID_SIGNATURE_ERROR *fiber.Error = fiber.NewError(
	4001802,
	"invalid input signature",
)

var INVALID_PADDR_ERROR *fiber.Error = fiber.NewError(
	4001803,
	"invalid input paddr",
)

var INTERNAL_SERVER_ERROR *fiber.Error = fiber.NewError(
	fiber.StatusInternalServerError,
	"internal server error",
)

var SERVICE_FAILED_DEPENDENCY_ERROR *fiber.Error = fiber.NewError(
	fiber.StatusFailedDependency,
	"failed dependency",
)
