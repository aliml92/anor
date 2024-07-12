package anor

//go:generate moq -pkg mocks -out mocks/user_service.go . UserService
//go:generate moq -pkg mocks -out mocks/product_service.go . ProductService
//go:generate moq -pkg mocks -out mocks/cart_service.go . CartService
//go:generate moq -pkg mocks -out mocks/order_service.go . OrderService
//go:generate moq -pkg mocks -out mocks/auth_service.go . AuthService
