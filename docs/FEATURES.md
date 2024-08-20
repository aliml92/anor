# Anor Features

This document outlines the current and planned features for the Anor e-commerce platform. Features are categorized and marked as implemented (✅), in progress (🚧), or planned (📅).

## User Authentication
- ✅ Signup
- ✅ Signin
- ✅ Forgot Password
- ✅ Logout
- ✅ Google Sign-in

## Home Page
- ✅ Featured products carousel
- 🚧 Popular products*

## Product Listing Page
### Sort by:
- 📅 Popular
- ✅ Price: Lowest
- ✅ Price: Highest
- 📅 Highest rated
- ✅ New arrivals
- 📅 Best sellers

### Filter by:
- ✅ Price range
- ✅ Brand
- 📅 Rating
- 📅 Custom attributes (size, color, material, etc.)

### Pagination:
- ✅ Page-based navigation
- ✅ "Show More" functionality (dynamically adds new product item cards)

## Search Functionality
- 🚧 Recent search queries
- 🚧 Trending search queries

### Search Autocomplete:
- ✅ Product suggestions
- 📅 Category suggestions
- 📅 Store suggestions

## Search Listing Page
### Sort by:
- 📅 Popular
- ✅ Price: Lowest
- ✅ Price: Highest
- 📅 Highest rated
- 📅 New arrivals
- 📅 Best sellers

### Filter by:
- ✅ Price range
- ✅ Brand
- 📅 Rating
- 📅 Custom attributes (size, color, material, etc.)

### Pagination:
- ✅ Page-based navigation
- ✅ "Show More" functionality (dynamically adds new product item cards)

## Product Details Page
- ✅ Product details display
- ✅ Additional product images carousel
- 📅 Product reviews and ratings
- 📅 Related products section
- ✅ "Add To Cart" by authenticated user
- ✅ "Add To Cart" by guest user
- 📅 "Buy It Now" button

## User Profile Page
- 📅 Display user information
- ✅ Order history
- 📅 Edit profile functionality
- 📅 Change password option
- 📅 Address book management

## User Cart Page
- ✅ Remove cart item for authenticated user
- ✅ Remove cart item for guest user
- ✅ Update cart item quantity for authenticated user
- ✅ Update cart item quantity for guest user
- 📅 Apply coupon
- 📅 Related products section

## Checkout Page (for authenticated users only)
- ✅ Create stripe payment intent
- ✅ Set delivery/billing addresses
- ✅ Set payment method
- ✅ Pay

## Product View Tracking
- ✅ Count product view on product listings page
- 📅 Count product view on product details page

## Additional Planned Features
- 📅 Wishlist management
- ✅ Merging guest cart with user cart when logged in
- 📅 Order tracking
- 📅 Recommendation engine for products
- 📅 Coupon and discount functionality
- 📅 User review and rating system
- 📅 Advanced search filters
- 📅 Mobile app version
- 📅 Admin/Seller pages (the whole other project)

*Note: Popular products feature is partially implemented. See [DEVELOPMENT.md](./DEVELOPMENT.md) for more details.

## Legend
- ✅ Implemented
- 🚧 In Progress
- 📅 Planned

We welcome contributions to help implement planned features or improve existing ones. Please see our [CONTRIBUTING.md](./CONTRIBUTING.md) for more information on how to get involved.