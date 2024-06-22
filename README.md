<img src="./assets/logo.svg" alt="anor" style="width: 150px; height: auto;">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/aliml92/anor/blob/master/LICENSE)

**Anor** is a full-stack e-commerce application inspired by the functionality and design of [Uzum](https://uzum.uz/uz), [eBay](https://www.ebay.com/), and [Zalando](https://en.zalando.de/).  It uses Golang for the backend, HTMX and _hyperscript with Bootstrap for the frontend, and PostgreSQL, Redis, and Typesense for data storage, caching, and search capabilities.


## Table of Contents
- [Features](#features)
- [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Getting Started](#getting-started)
    - [Additional Configuration](#additional-configuration)
- [Testing](#testing)
- [Learning Resources](#learning-resources)
    - [Resources Used](#resources-used)
    - [Resources to Explore](#resources-to-explore)
- [Contributing](#contributing)
- [License](#license)


## Features

- **User Authentication:** 
  - [x] Signup 
  - [x] Signin 
  - [x] Forgot Password 
  - [x] Logout 
  - [ ] Google Sign-in 
- **Product listing page with categories:**
  - Sort by:
    - [ ] Popularity
    - [x] Price: Lowest
    - [x] Price: Highest
    - [ ] Highest rated
    - [ ] New arrivals
    - [ ] Best sellers
  - Filter by:
    - [x] Price range
    - [x] Brand
    - [ ] Rating
    - [ ] Custom attributes (size, color, material, etc.)
  - Pagination:
    - [x] Page-based navigation
    - [x] "Show More" functionality (dynamically adds new product item cards)
- **Product Details Page:**
  - [x] Product details display
  - [x] Additional product images carousel
  - [ ] Product reviews and ratings
  - [ ] Related products section
  - [x] Add to cart button

- Customer Profile Page:
  - [ ] Display user information
  - [ ] Order history
  - [ ] Edit profile functionality
  - [ ] Change password option
  - [ ] Address book management
- [ ] Shopping cart functionality
- [ ] Checkout and payment integration
- [ ] Order tracking
- [ ] Recommendation engine for products
- [ ] Coupon and discount functionality
- [ ] User review and rating system
- [ ] Advanced search filters
- [ ] Mobile app version

## Installation

### Prerequisites
Before installing the application, ensure you have the following tools installed on your machine:
- Docker/Docker Compose.
- sqlc: For generating type-safe code from SQL.
- goose: For database migrations.
- task: A task runner/simpler Make alternative written in Go.

### Getting Started
To get the application running locally, follow these steps:

1. **Start necessary services**

   Ensure you start Postgres, Typesense and Redis using Docker:
    ```bash
    task compose-up
    ```
2. **Run database migrations**:

    Apply database migrations with goose::
    ```
    task goose-up
    ``` 
3. **Import sample data**

    Populate the database with the initial dataset:
    ```
    task import-dataset
    ```
4. **Start the application**

   Populate the database with the initial dataset:
   ```
   export CONFIG_FILEPATH=./config.dev.yaml
   go run cmd/anor/*.go
   ```
Project starts on port 8008 by default.

### Additional Configuration  
Adjust additional settings and configurations as needed:

- **Configuration Files**:  Customize `config.dev.yaml` as needed to tweak database connections, service endpoints, and other critical settings. 

## Testing
Coming soon

## Learning Resources

### Resources Used
Here are some of the resources that have been instrumental in the development of this eCommerce platform:

- **Official documentations of [Go](https://tour.golang.org/), [htmx](https://htmx.org/docs/), [_hyperscript](https://hyperscript.org), [Bootstrap](https://getbootstrap.com/docs/5.1) , [PostgreSQL](https://www.postgresql.org/docs/current/index.html), [Typesense](https://typesense.org/docs/26.0/api/).**
- **[Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)** (Medium article): A helpful guide to structuring Go projects effectively.
- **[WTF Dial](https://github.com/benbjohnson/wtf)**: A repository providing insights into organizing project structure, useful for project structuring reference.
- **[Modelling Products and Variants for E-Commerce](https://martinbean.dev/blog/2023/01/27/product-variants-laravel/)**: A blog post on database modeling for "products," "product attributes," "SKUs," and "product attribute SKUs" tables in e-commerce applications.
- **[Sharding & IDs at Instagram](https://instagram-engineering.com/sharding-ids-at-instagram-1cf5a71e5a5c)**: An insightful article discussing the implementation of sharding and ID generation at Instagram.
- **[Sharding and IDs at Instagram](https://news.ycombinator.com/item?id=3058327)**: A Hacker News discussion thread related to the Medium article "Sharding & IDs at Instagram," providing additional insights and perspectives on the topic.
- **[Optimistic UI Patterns for Improved Perceived Performance](https://simonhearne.com/2021/optimistic-ui-patterns/)**. 
- **[Commerce for devs](https://commercefordevs.org/)**: A publication for developers who want to become experts in ecommerce.
- **[Optimistic checkouts](https://commercefordevs.org/optimistic-checkouts/)**: Blog post about optimistic checkouts in Ecommerce websites.

### Resources to Explore
Future resources to deepen knowledge and enhance the application:

- **[Intelligent Sort (I.S.): A New Method for Product Sorting in E-Commerce](https://medium.com/@khosravi.official/intelligent-sort-i-s-a-new-method-for-product-sorting-in-e-commerce-6d4f1d11c340)**: (Medium article): A helpful guide to implement the `Popular` sorting algorithm, commonly used as the default sorting method in product listing and search results pages for enhanced relevance.


## Contributing

Contributions to the project are welcome! If you have suggestions or improvements, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

We appreciate your input and look forward to improving the application together!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
