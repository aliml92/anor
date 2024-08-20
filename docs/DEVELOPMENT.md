# Development Guide

This guide will help you set up the Anor project for local development.

## Prerequisites

Ensure you have the following tools installed on your machine:
- Docker/Docker Compose.
- sqlc: For generating type-safe code from SQL.
- goose: For database migrations.
- task: A task runner/simpler Make alternative written in Go.
- air: For live reload.

## Getting Started
To get the application running locally, follow these steps:

1. **Start necessary services**

   Ensure you start Postgres, Typesense and Redis using Docker:
    ```bash
    task compose:up
    ```
2. **Run database migrations**:

   Apply database migrations with goose:
    ```
    task goose:up
    ```
3. **Configure the application**
   Copy the sample configuration file and fill in the necessary fields.
   ```
   cp ./config/config.sample.yaml ./config/config.dev.yaml
   ```
4. **Start the application**
   You can start the application using one of the following methods:
   a. Using air for live reload:
   ```
   export CONFIG=./config/config.dev.yaml
   air
   ```
   b. Or, run it directly with Go:
   ```
   CONFIG=./config/config.dev.yaml go run ./cmd/anor/*.go 
   ```
   The application will start on port 8008 by default.  
5. **Import sample data**

   Populate the database with the initial dataset:
    ```
    task import-dataset -- -ignore-imported
    ```
    The might be issues during the execution of the above command. 
    Even though the issue you will have enough data imported to both
    Postgres and Typesense.
6. **Setting up Popular Products**
   
   The "Popular Products" feature should ideally use the Intelligent Sort (I.S.) method, as described in this article:
   **[Intelligent Sort (I.S.): A New Method for Product Sorting in E-Commerce](https://medium.com/@khosravi.official/intelligent-sort-i-s-a-new-method-for-product-sorting-in-e-commerce-6d4f1d11c340)**.
  
   Current Implementation:
   1. A partial implementation of the I.S. method exists in [script.js](https://github.com/aliml92/anor/blob/main/web/static/js/script.js#L103-L214).
   2. This script tracks "seen" products on the product listings page and sends their IDs to the `/analytics/pl/views` route.
   3. However, the collected data is not yet processed or used for ranking.

   To fully implement the Popular Products feature:
   1. Create a ranking table that combines the "seen" product data with sales data.
   2. Implement the complete I.S. algorithm (or similar one) using this combined data.

   Temporary solution:
   Currently, the homepage displays a fixed set of "popular" products. To use this temporary solution.
   1. Add product IDs to your `.env` file in the following format:
    ```
    FAKE_POPULAR_PRODUCT_IDS=250973575448303,250973576497055
    ```
   2. The application will use these IDs to display "popular" products on the homepage.

7. **Add Featured Selections**

   To display featured selections on the homepage, you need to manually add records to the `featured_selections` table in the database. Follow these steps:

   a. Find the correct category IDs:
   - Navigate to the categories in your browser, e.g., http://localhost:8008/categories/furniture-149058613
   - The last part of the URL (e.g., furniture-149058613) contains the category ID

   b. Run the following SQL commands, replacing `<CATEGORY_ID>` with the actual IDs you found:

   ```sql
   -- Insert for Office Notebooks
   INSERT INTO featured_selections 
   (resource_path, banner_info, image_url, display_order, start_date, end_date)
   VALUES 
   (
       'http://localhost:8008/categories/note-pads-<CATEGORY_ID>',
       '{"title": "Office Notebooks"}',
       'static/images/banners/notepad.png',
       3,
       CURRENT_DATE - INTERVAL '1 day',
       CURRENT_DATE + INTERVAL '29 days'
   );

   -- Insert for Home Furniture
   INSERT INTO featured_selections 
   (resource_path, banner_info, image_url, display_order, start_date, end_date)
   VALUES 
   (
       'http://localhost:8008/categories/furniture-<CATEGORY_ID>',
       '{"title": "Home Furniture"}',
       'static/images/banners/interior-wide.jpg',
       2,
       CURRENT_DATE - INTERVAL '1 day',
       CURRENT_DATE + INTERVAL '29 days'
   );

   -- Insert for Athletic Shoes
   INSERT INTO featured_selections 
   (resource_path, banner_info, image_url, display_order, start_date, end_date)
   VALUES 
   (
       'http://localhost:8008/categories/athletic-shoes-<CATEGORY_ID>',
       '{"title": "Athletic Shoes"}',
       'static/images/banners/adidas.webp',
       1,
       CURRENT_DATE - INTERVAL '1 day',
       CURRENT_DATE + INTERVAL '29 days'
   );
   ```
   **Important:**
   * Replace `<CATEGORY_ID>` in each `resource_path` with the actual category IDs you found in step a.
   * If these data are not inserted, the homepage will not display the featured selections.
   * Ensure the image files (notepad.png, interior-wide.jpg, adidas.webp) exist in your `static/images/banners/` directory.
## Next Steps
   1. Familiarize yourself with the project structure.   
      (Maybe I will make a short video explaining the product setup and structure)
   2.  Check out the FEATURES.md file to see what features are implemented and what's not.
   3. Consider contributing to the project by picking up an open issue or proposing new features.

Happy coding!   