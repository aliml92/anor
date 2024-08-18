// htmx.logAll();

// enable tooltips everywhere
const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
})


document.addEventListener("anor:showToast", showToast);

function showToast(evt){
    console.log("showToast fired");
    const value = evt.detail.value;
    Toastify({
        text: value,
        duration: 3000,
        destination: "#",
        newWindow: false,
        close: true,
        gravity: "bottom", // `top` or `bottom`
        position: "right", // `left`, `center` or `right`
        stopOnFocus: true, // Prevents dismissing of toast on hover
        style: {
            background: "linear-gradient(to right, #00b09b, #96c93d)",
        },
        onClick: function(){} // Callback after click
    }).showToast();
}

function getPriceRange() {
    let priceRangeSlider = document.querySelector('.js-range-slider');
    let minPrice = priceRangeSlider.getAttribute("data-min");
    let maxPrice = priceRangeSlider.getAttribute("data-max");
    return minPrice + "-" + maxPrice;
}

function getPriceFrom() {
    let fromInput = document.querySelector('.js-input-from');
    let priceRangeSlider = document.querySelector('.js-range-slider');

    // Convert both values to numbers for comparison
    let minPrice = parseFloat(priceRangeSlider.getAttribute("data-min"));
    let fromValue = parseFloat(fromInput.value);

    // Return fromInput.value if it is greater than minPrice, otherwise return minPrice
    return fromValue > minPrice ? fromValue : null;
}

function getPriceTo() {
    let toInput = document.querySelector('.js-input-to');
    let priceRangeSlider = document.querySelector('.js-range-slider');

    // Convert both values to numbers for comparison
    let maxPrice = parseFloat(priceRangeSlider.getAttribute("data-max"));
    let toValue = parseFloat(toInput.value);

    // Return toInput.value if it is less than maxPrice, otherwise return maxPrice
    return toValue < maxPrice ? toValue : null;
}

function getCheckedBrands() {
    // Select all checked checkboxes inside the .side-filter-brands div
    const checkedBoxes = document.querySelectorAll('.side-filter-brands .form-check-input:checked');

    // Map the values of the checked checkboxes to an array
    const checkedValues = Array.from(checkedBoxes).map(box => box.value);

    // Check if there are any checked boxes and return a comma-separated string of values or null
    return checkedValues.length > 0 ? checkedValues.join(',') : null;
}

function getSort() {
    let selectElement = document.getElementById("sort-selector");
    if (selectElement) {
        return selectElement.value;
    } else {
        console.error("Element not found with selector:", '#sort-selector');
        return null; // or an appropriate default/fallback value
    }
}

function getQ() {
    let searchInputEl = document.getElementById("searchInput");
    if (searchInputEl) {
        let val = searchInputEl.value.trim();
        if (val) {
            return val;
        } else {
            return "*";
        }
    } else {
        return "*";
    }
}

document.querySelectorAll('.dropdown-menu').forEach(function (element) {
    element.addEventListener('click', function (e) {
        e.stopPropagation();
    });
});

// the following scripts are used to identify which product items are
// observed on product listings page and send product ids of these
// "seen" products as batch to server. The result of such requests
// are used to calculate popular products index
var viewedProducts = {};
var productTimeouts = {};
var activeObservers = {};

// Capture initial views upon page load
window.addEventListener('load', () => {
    if (window.location.pathname.includes('/categories')) {
        processProductViews();
    }
});

// URL change listener for non-SPA sites
window.addEventListener('htmx:afterSettle', function(evt) {
    if (evt.detail.requestConfig.path.includes("/categories/")) {
        processProductViews();
    }
});


function processProductViews() {
    for (let key in activeObservers) {
        if (activeObservers.hasOwnProperty(key)) {
            activeObservers[key].disconnect();
        }
    }
    activeObservers = {};

    const productCards = document.querySelectorAll('figure.card-product-grid');

    productCards.forEach(card => {
        let productId = card.getAttribute('data-product-id');
        if (!viewedProducts[productId] || (Date.now() - viewedProducts[productId].lastSent) >= 3600000) {
            observeProductCard(card, productId);
        } else {
            console.log('Product already observed within the last hour:', productId);
        }
    });
}

function observeProductCard(card, productId) {
    let observer = new IntersectionObserver(entries => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                // Start or reset the timer when the card is intersected
                clearTimeout(productTimeouts[productId]);
                productTimeouts[productId] = setTimeout(() => {
                    // After 5 seconds, check if the entry is still intersecting
                    if (entry.isIntersecting) {
                        // Assuming previous check in processProductViews is sufficient
                        viewedProducts[productId] = { lastSent: Date.now() };
                        sendProductViewData(productId);

                        // observer.unobserve(entry.target);
                    }
                }, 4000); // 4 seconds
            } else {
                // If the card is not intersecting, clear the timer
                clearTimeout(productTimeouts[productId]);
            }
        });
    }, { threshold: 1 });

    observer.observe(card);
    activeObservers[productId] = observer; // Track the observer by productId
}

var batchedProductIds = new Set();
var batchTimeout = null;

function sendProductViewData(productId) {
    console.log('Collecting view data for product:', productId);
    batchedProductIds.add(productId);

    clearTimeout(batchTimeout);

    batchTimeout = setTimeout(() => {
        sendBatchedProductViews();
    }, 3000000);

}

function sendBatchedProductViews() {
    if (batchedProductIds.size > 0) {
        console.log('Sending batched view data:', Array.from(batchedProductIds));
        fetch('/analytics/plp/views', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(
                { productIds: Array.from(batchedProductIds) }
            ),
        }).then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // or handle based on your API structure
        }).then(data => {
            console.log('Server response:', data);
        }).catch(error => {
            console.error('Failed to send batched product views:', error);
        });

        // Clear the batched product IDs set after sending
        batchedProductIds.clear();
        clearTimeout(batchTimeout);
    }
}

// selectActiveItem get selected text in autocomplete dropdown
function selectActiveItem(evt) {
    evt.preventDefault();
    const listItems = document.querySelectorAll('#search-dropdown .list-group-item');
    const activeItem = document.querySelector('#search-dropdown .list-group-item.active');
    let nextItem;

    if (!activeItem) {
        nextItem = listItems[0];
    } else {
        const activeIndex = Array.from(listItems).indexOf(activeItem);
        if (evt.key === 'ArrowDown') {
            nextItem = listItems[(activeIndex + 1) % listItems.length];
        } else if (evt.key === 'ArrowUp') {
            nextItem = listItems[(activeIndex - 1 + listItems.length) % listItems.length];
        }
    }

    listItems.forEach(item => item.classList.remove('active'));
    nextItem.classList.add('active');
    // nextItem.focus();
    evt.target.value = nextItem.textContent.trim();
}

document.addEventListener("htmx:configRequest", function(configEvent){
    const path = configEvent.detail.path;
    const hasTargetPath = ["/categories", "/search"].some(segment => path.includes(segment));
    if (hasTargetPath) {
        console.log("Before:", configEvent.detail.parameters);

        let filteredParameters = {};

        Object.entries(configEvent.detail.parameters).forEach(([key, value]) => {
            if (value !== null) {
                filteredParameters[key] = value;
            }
        });

        configEvent.detail.parameters = filteredParameters;

        console.log("After:", configEvent.detail.parameters);
    }
})

window.addEventListener('htmx:historyRestore', (evt) => {
    if (evt.detail.path.includes("/categories/") || evt.detail.path.includes("/search")) {
        // reinit price range slider
        initializePriceRangeSlider();
        let $range = $(".js-range-slider");
        $range.parents('div.range-slider').removeClass('invisible');

        // add event listeners to 'More'/'Less' button
        addEventListenersToBrandsBtn();
    }
});

window.addEventListener('htmx:beforeHistorySave', (evt) => {
    if (evt.detail.path.includes("/categories/") || evt.detail.path.includes("/search")) {
        window.priceRangeSlider.destroy();
    }
});

window.addEventListener('DOMContentLoaded', (evt) => {
    if (window.location.pathname.includes('/categories/') || window.location.pathname.includes("/search")) {
        // reinit price range slider
        initializePriceRangeSlider();
        let $range = $(".js-range-slider");
        $range.parents('div.range-slider').removeClass('invisible');

        // add event listeners to 'More'/'Less' button
        addEventListenersToBrandsBtn();
    }
})

document.addEventListener("htmx:afterSettle", (evt) => {
    if (evt.detail.requestConfig.path.includes("/categories/") || evt.detail.requestConfig.path.includes("/search")) {
        // reinit price range slider
        initializePriceRangeSlider();
        let $range = $(".js-range-slider");
        $range.parents('div.range-slider').removeClass('invisible');

        // add event listeners to 'More'/'Less' button
        addEventListenersToBrandsBtn();
    }
});

function initializePriceRangeSlider() {
    // setup price range slider
    let $range = $(".js-range-slider");
    let $inputFrom = $(".js-input-from");
    let $inputTo = $(".js-input-to");
    let min = 0;
    let max = 0;
    let from = 0;
    let to = 0;

    $range.ionRangeSlider({
        onStart: updateInputs,
        onChange: updateInputs,
        onFinish: triggerFilterRequest,
    });

    let instance = $range.data("ionRangeSlider");
    window.priceRangeSlider = instance;

    function updateInputs(data) {
        from = data.from;
        to = data.to;

        $inputFrom.prop("value", from);
        $inputTo.prop("value", to);
    }


    function triggerFilterRequest(data) {
        htmx.trigger(".js-range-slider", "filterByPriceRange")
    }

    $inputFrom.on("input", function () {
        let val = $(this).prop("value");

        // validate
        if (val < min) {
            val = min;
        } else if (val > to) {
            val = to;
        }

        instance.update({
            from: val
        });
    });
    //
    $inputTo.on("input", function () {
        let val = $(this).prop("value");

        // validate
        if (val < from) {
            val = from;
        } else if (val > max) {
            val = max;
        }

        instance.update({
            to: val
        });
    });
}

function addEventListenersToBrandsBtn() {
    document.querySelectorAll('.side-filter-collapsible').forEach(function (element) {
        element.addEventListener('show.bs.collapse', function (e) {
            let siblingListItem = e.currentTarget.parentElement.nextElementSibling;
            let icon = siblingListItem.querySelector("i");
            let buttonText = siblingListItem.querySelector("span");
            icon.classList.remove('bi-chevron-down');
            icon.classList.add('bi-chevron-up');
            buttonText.textContent = ' Less';
        });

        element.addEventListener('hide.bs.collapse', function (e) {
            let siblingListItem = e.currentTarget.parentElement.nextElementSibling;
            let icon = siblingListItem.querySelector("i");
            let buttonText = siblingListItem.querySelector("span");
            icon.classList.remove('bi-chevron-up');
            icon.classList.add('bi-chevron-down');
            buttonText.textContent = ' More';
        });
    });
}

// search input
function addSearchQuery(query) {
    const maxSize = 5;
    const localStorageKey = 'searchQueries';

    // Check if the query is a non-empty string
    if (typeof query !== 'string' || query.trim() === '') {
        console.warn('Query is empty or not a string.');
        return;
    }

    // Retrieve the existing search queries from local storage
    let searchQueries = JSON.parse(localStorage.getItem(localStorageKey)) || [];

    // Remove the query if it already exists to avoid duplicates
    searchQueries = searchQueries.filter(existingQuery => existingQuery !== query);

    // Add the new query to the beginning of the array
    searchQueries.unshift(query);

    // Ensure the array length does not exceed the maximum size
    if (searchQueries.length > maxSize) {
        searchQueries.pop();
    }

    // Save the updated search queries back to local storage
    localStorage.setItem(localStorageKey, JSON.stringify(searchQueries));
}

function populateSearchDropdownListWithRecentSearches() {
    const localStorageKey = 'searchQueries';
    const searchDropdownList = document.getElementById('search-dropdown-list');

    searchDropdownList.innerHTML = '';

    const searchQueries = JSON.parse(localStorage.getItem(localStorageKey)) || [];

    const trendingProductsCache = JSON.parse(localStorage.getItem('trendingProductsCache')) || [];
    const cacheExpiration = new Date(trendingProductsCache.expiration);
    if (trendingProductsCache.products && new Date() < cacheExpiration) {
        // Use cached trending products
        console.log("rendered from cache");
        renderRecentSearches(searchQueries);
        renderTrendingProducts(trendingProductsCache.products);
    } else {
        console.log("made request")
        // Fetch trending products from "/trending-products" API
        fetch('/trending-products')
          .then(response => response.json())
          .then(data => {
              if (Array.isArray(data) && data.length > 0) {
                  renderRecentSearches(searchQueries);
                  renderTrendingProducts(data);
                  // Cache trending products
                  localStorage.setItem('trendingProductsCache', JSON.stringify({ products: data, expiration: Date.now() + 900000 }));
              } else {
                  // Populate with recent searches if the response is empty
                  renderRecentSearches(searchQueries);
              }
          })
          .catch(error => {
              console.error('Error fetching trending products:', error);
              // Populate with recent searches if the request fails
              renderRecentSearches(searchQueries);
          });
    }

    function renderTrendingProducts(products) {
        products.forEach((query, index) => {
            // Create the list item element
            const listItem = document.createElement('a');
            listItem.href = '#';
            listItem.classList.add('list-group-item', 'list-group-item-action', 'd-flex');
            if (index !== products.length - 1) {
                listItem.classList.add('border-bottom-0');
            }

            // Add the clock icon, query text, and close icon
            listItem.innerHTML = `
                <i class="me-3">
                    <svg xmlns="http://www.w3.org/2000/svg" height="20px" viewBox="0 0 24 24" width="20px" fill="#5f6368"><path d="M0 0h24v24H0z" fill="none"/><path d="M16 6l2.29 2.29-4.88 4.88-4-4L2 16.59 3.41 18l6-6 4 4 6.3-6.29L22 12V6z"/></svg>
                </i>
                <span>${query}</span>
            `;

            listItem.addEventListener('click', function(event) {
                event.preventDefault();
                const searchInput = document.getElementById('searchInput');
                searchInput.value = event.currentTarget.querySelector("span").textContent;
                htmx.trigger(searchInput, "searchTrigger");
            });


            // Append the list item to the dropdown list
            searchDropdownList.appendChild(listItem);
        });
    }

    function renderRecentSearches(searchQueries) {
        searchQueries.forEach((query, index) => {
            // Create the list item element
            const listItem = document.createElement('a');
            listItem.href = '#';
            listItem.classList.add('list-group-item', 'list-group-item-action', 'd-flex', 'history');
            if (index !== searchQueries.length - 1) {
                listItem.classList.add('border-bottom-0');
            }

            // Add the clock icon, query text, and close icon
            listItem.innerHTML = `
            <i class="bi bi-clock-history me-3"></i>
                <span>${query}</span>
                <i class="bi bi-x ms-auto" data-index="${index}"></i>
            `;

            listItem.addEventListener('click', function(event) {
                event.preventDefault();
                console.log("click fired 1")
                const searchInput = document.getElementById('searchInput');
                searchInput.value = event.target.querySelector("span").textContent;
                htmx.trigger(searchInput, "searchTrigger");
            });

            // Attach an event listener to the close icon
            listItem.querySelector('.bi-x').addEventListener('click', function(event) {
                event.preventDefault();
                console.log("click fired 2")
                const queryIndex = parseInt(this.getAttribute('data-index'));
                removeSearchQuery(queryIndex);
                populateSearchDropdownListWithRecentSearches(); // Refresh the list after removal
                document.getElementById("searchInput").focus();
                event.stopPropagation();
            });

            // Append the list item to the dropdown list
            searchDropdownList.appendChild(listItem);
        });
    }

}

function removeSearchQuery(index) {
    const localStorageKey = 'searchQueries';
    let searchQueries = JSON.parse(localStorage.getItem(localStorageKey)) || [];
    searchQueries.splice(index, 1);
    localStorage.setItem(localStorageKey, JSON.stringify(searchQueries));
}


// ============= Product details page ============== 
function setupProductVariantSelectOptions() {
    let productVariantScriptEl = document.getElementById('product-variant-matrix');
    let productVariantDim = parseInt(productVariantScriptEl.getAttribute('data-dim'));
    let productVariantMatrix = JSON.parse(productVariantScriptEl.textContent);

    let quantitySpan = document.getElementById('quantity-in-stock');
    let quantityInput = document.querySelector('input[name="quantity"]');

    // Get all select elements
    let selects = document.querySelectorAll('.attr-select');

    // Get all button elements
    let minusButton = document.querySelector('.button-minus');
    let plusButton = document.querySelector('.button-plus');

    switch (productVariantDim) {
        case 0:
            // TODO: handle case where is only one productVariant (default productVariant)
            let defaultProductVariantID = productVariantMatrix[0].id;
            productVariantScriptEl.setAttribute('data-product-variant-id', defaultProductVariantID);
            if (quantitySpan.textContent !== "None left") {
                enableQuantityToggle(minusButton, plusButton, quantityInput);
            }
            break;
        case 1:
            selects.forEach(function(select) {
                updateOptions1(select);
                select.addEventListener('change', function() {
                    let i = null;
                    selects.forEach(function(sel) {
                        const selectedIndex = sel.options[sel.selectedIndex].getAttribute("data-attr-val-index");
                        if (sel.getAttribute('data-attr-index') === '0') {
                            i = selectedIndex;
                        }
                    });

                    console.log('i:', i);

                    if (i !== null) {
                        // update productVariantID
                        let selectedProductVariantID = productVariantMatrix[i].id;
                        productVariantScriptEl.setAttribute('data-product-variant-id', selectedProductVariantID);

                        // enable quantity increment/decrements
                        enableQuantityToggle(minusButton, plusButton, quantityInput);

                        // hide select error message
                        const invalidFeedback = select.parentElement.querySelector(".invalid-feedback");
                        invalidFeedback.classList.remove('show');

                        // construct qty
                        let qty = productVariantMatrix[i].quantity;
                        if ( qty === 0) {
                            disableQuantityToggle(minusButton, plusButton, quantityInput);
                            quantitySpan.textContent = "None left";
                        } else if (qty === 1) {
                            quantitySpan.textContent = "Only one left";
                        } else if (qty <= 10) {
                            quantitySpan.textContent = qty + " left";
                        } else {
                            quantitySpan.textContent = "More than 10 available";
                        }
                        quantityInput.max = productVariantMatrix[i].quantity;
                    } else {
                        // disable quantity increment/decrement since i === null means it is default option selected
                        disableQuantityToggle(minusButton, plusButton, quantityInput);
                    }
                });
            });

            function updateOptions1(select) {
                const options = select.options;
                for (let i = 1; i < options.length; i++) {
                    const attrValIndex = options[i].getAttribute("data-attr-val-index");
                    if (productVariantMatrix[parseInt(attrValIndex)].quantity === 0) {
                        if (!options[i].textContent.includes("(Out of stock)")) {
                            options[i].textContent += " (Out of stock)";
                        }
                        options[i].disabled = true;
                    }
                }
            }
            break;

       // handle productVariant variation options when there are two product attributes
        case 2:
            function updateAvailability(selects, productVariantMatrix) {
                let selectedIndexes = Array.from(selects).map(select =>
                  select.options[select.selectedIndex].getAttribute("data-attr-val-index"));

                let [i, j] = selectedIndexes.map(index => index !== null ? parseInt(index) : null);

                if (i !== null && j !== null) {
                    let productVariant = productVariantMatrix[i][j];
                    let qty = productVariant.quantity;

                    productVariantScriptEl.setAttribute('data-product-variant-id', productVariant.id);

                    if (qty === 0) {
                        disableQuantityToggle(minusButton, plusButton, quantityInput);
                        quantitySpan.textContent = "None left";
                    } else {
                        enableQuantityToggle(minusButton, plusButton, quantityInput);
                        quantitySpan.textContent = qty === 1 ? "Only one left" :
                          qty < 10 ? qty + " left" :
                            "More than 10 available";
                    }
                    quantityInput.max = qty;
                } else {
                    disableQuantityToggle(minusButton, plusButton, quantityInput);
                    quantitySpan.textContent = "More than 10 available";
                }
            }

            function updateOptions2(selects, productVariantMatrix) {
                Array.from(selects).forEach((currentSelect, currentSelectIndex) => {
                    const otherSelectIndex = 1 - currentSelectIndex;
                    const otherSelect = selects[otherSelectIndex];
                    const selectedOtherIndex = otherSelect.selectedIndex !== 0 ?
                      parseInt(otherSelect.options[otherSelect.selectedIndex].getAttribute("data-attr-val-index")) : null;

                    Array.from(currentSelect.options).forEach((option, optionIndex) => {
                        if (optionIndex === 0) return; // Skip the default option

                        const currentAttrValIndex = parseInt(option.getAttribute("data-attr-val-index"));

                        let outOfStock = selectedOtherIndex !== null ?
                          (currentSelectIndex === 0 ?
                            productVariantMatrix[currentAttrValIndex][selectedOtherIndex].quantity === 0 :
                            productVariantMatrix[selectedOtherIndex][currentAttrValIndex].quantity === 0) :
                          Array.from(otherSelect.options).every((otherOption, otherOptionIndex) => {
                              if (otherOptionIndex === 0) return true;
                              const otherAttrValIndex = parseInt(otherOption.getAttribute("data-attr-val-index"));
                              return currentSelectIndex === 0 ?
                                productVariantMatrix[currentAttrValIndex][otherAttrValIndex].quantity === 0 :
                                productVariantMatrix[otherAttrValIndex][currentAttrValIndex].quantity === 0;
                          });

                        option.disabled = outOfStock;
                        option.textContent = option.textContent.replace(" (Out of stock)", "");
                        if (outOfStock) {
                            option.textContent += " (Out of stock)";
                        }
                    });
                });
            }

            // Initial update
            updateOptions2(selects, productVariantMatrix);

            selects.forEach((select, selectIndex) => {
                select.addEventListener('change', function(evt) {
                    // Hide select error message
                    if (evt.target.selectedIndex !== 0) {
                        const invalidFeedback = evt.target.nextElementSibling;
                        if (invalidFeedback && invalidFeedback.classList.contains('invalid-feedback')) {
                            invalidFeedback.classList.remove('show');
                        }
                    }

                    updateAvailability(selects, productVariantMatrix);

                    // Update options for the other select
                    const otherSelectIndex = 1 - selectIndex;
                    const otherSelect = selects[otherSelectIndex];
                    const selectedIndex = evt.target.options[evt.target.selectedIndex].getAttribute("data-attr-val-index");

                    if (selectedIndex !== null) {
                        Array.from(otherSelect.options).forEach((option, optionIndex) => {
                            if (optionIndex === 0) return; // Skip the default option

                            const attrValIndex = parseInt(option.getAttribute("data-attr-val-index"));
                            const quantity = selectIndex === 0 ?
                              productVariantMatrix[parseInt(selectedIndex)][attrValIndex].quantity :
                              productVariantMatrix[attrValIndex][parseInt(selectedIndex)].quantity;

                            option.disabled = quantity === 0;
                            option.textContent = option.textContent.replace(" (Out of stock)", "");
                            if (quantity === 0) {
                                option.textContent += " (Out of stock)";
                            }
                        });
                    } else {
                        // If default option is selected, reset the other select
                        updateOptions2(selects, productVariantMatrix);
                    }
                });
            });

            break;
        case 3:
            function updateAvailability3(selects, productVariantMatrix) {
                let selectedIndexes = Array.from(selects).map(select =>
                  select.options[select.selectedIndex].getAttribute("data-attr-val-index"));

                let [i, j, k] = selectedIndexes.map(index => index !== null ? parseInt(index) : null);

                if (i !== null && j !== null && k !== null) {
                    let productVariant = productVariantMatrix[i]?.[j]?.[k];
                    if (productVariant) {
                        let qty = productVariant.quantity;

                        productVariantScriptEl.setAttribute('data-product-variant-id', productVariant.id);

                        if (qty === 0) {
                            disableQuantityToggle(minusButton, plusButton, quantityInput);
                            quantitySpan.textContent = "None left";
                        } else {
                            enableQuantityToggle(minusButton, plusButton, quantityInput);
                            quantitySpan.textContent = qty === 1 ? "Only one left" :
                              qty < 10 ? qty + " left" :
                                "More than 10 available";
                        }
                        quantityInput.max = qty;
                    } else {
                        disableQuantityToggle(minusButton, plusButton, quantityInput);
                        quantitySpan.textContent = "Unavailable";
                    }
                } else {
                    disableQuantityToggle(minusButton, plusButton, quantityInput);
                    quantitySpan.textContent = "";
                }
            }

            function updateOptions3(selects, productVariantMatrix) {
                Array.from(selects).forEach((currentSelect, currentSelectIndex) => {
                    const otherSelectIndexes = [0, 1, 2].filter(index => index !== currentSelectIndex);
                    const otherSelects = otherSelectIndexes.map(index => selects[index]);

                    const selectedOtherIndexes = otherSelects.map(select =>
                      select.selectedIndex !== 0 ? parseInt(select.options[select.selectedIndex].getAttribute("data-attr-val-index")) : null);

                    Array.from(currentSelect.options).forEach((option, optionIndex) => {
                        if (optionIndex === 0) return; // Skip the default option

                        const currentAttrValIndex = parseInt(option.getAttribute("data-attr-val-index"));

                        let outOfStock = true;

                        if (selectedOtherIndexes.every(index => index !== null)) {
                            // If both other selects have a selection
                            let [otherIndex1, otherIndex2] = selectedOtherIndexes;
                            let variant = getVariant(currentSelectIndex, currentAttrValIndex, otherIndex1, otherIndex2, productVariantMatrix);
                            outOfStock = !variant || variant.quantity === 0;
                        } else {
                            // Check all combinations
                            for (let i = 1; i < otherSelects[0].options.length; i++) {
                                for (let j = 1; j < otherSelects[1].options.length; j++) {
                                    const index1 = parseInt(otherSelects[0].options[i].getAttribute("data-attr-val-index"));
                                    const index2 = parseInt(otherSelects[1].options[j].getAttribute("data-attr-val-index"));

                                    let variant = getVariant(currentSelectIndex, currentAttrValIndex, index1, index2, productVariantMatrix);

                                    if (variant && variant.quantity > 0) {
                                        outOfStock = false;
                                        break;
                                    }
                                }
                                if (!outOfStock) break;
                            }
                        }

                        option.disabled = outOfStock;
                        option.textContent = option.textContent.replace(" (Out of stock)", "");
                        if (outOfStock) {
                            option.textContent += " (Out of stock)";
                        }
                    });
                });
            }

            function getVariant(currentSelectIndex, currentAttrValIndex, index1, index2, productVariantMatrix) {
                switch(currentSelectIndex) {
                    case 0:
                        return productVariantMatrix[currentAttrValIndex]?.[index1]?.[index2];
                    case 1:
                        return productVariantMatrix[index1]?.[currentAttrValIndex]?.[index2];
                    case 2:
                        return productVariantMatrix[index1]?.[index2]?.[currentAttrValIndex];
                }
            }

                // Initial update
                updateOptions3(selects, productVariantMatrix);

                selects.forEach((select, selectIndex) => {
                    select.addEventListener('change', function(evt) {
                        // Hide select error message
                        if (evt.target.selectedIndex !== 0) {
                            const invalidFeedback = evt.target.nextElementSibling;
                            if (invalidFeedback && invalidFeedback.classList.contains('invalid-feedback')) {
                                invalidFeedback.classList.remove('show');
                            }
                        }

                        updateAvailability3(selects, productVariantMatrix);

                        // Update options for the other selects
                        updateOptions3(selects, productVariantMatrix);
                    });
                });

                break;
        default:
            console.log("no attributes of this product")
    }
}

function incrementValue(e) {
    e.preventDefault();
    const target = e.target;
    const fieldName = target.getAttribute("data-field");
    const parent = target.closest(".input-group");
    const inputField = parent.querySelector("input[name=\"" + fieldName + "\"]");
    let currentVal = parseInt(inputField.value, 10);
    let max = parseInt(inputField.getAttribute('max'));

    if (!isNaN(currentVal) && currentVal < max) {
        inputField.value = currentVal + 1;
    } else {
        inputField.value = max;
    }
}

function decrementValue(e) {
    e.preventDefault();
    const target = e.target;
    const fieldName = target.getAttribute("data-field");
    const parent = target.closest(".input-group");
    let inputField = parent.querySelector("input[name=\"" + fieldName + "\"]");
    let currentVal = parseInt(inputField.value, 10);

    if (!isNaN(currentVal) && currentVal > 1) {
        inputField.value = currentVal - 1;
    } else {
        inputField.value = 1;
    }
}

function disableQuantityToggle(minusButton, plusButton, quantityInput) {
    // Check if any of the elements is already disabled
    if (!minusButton.disabled || !plusButton.disabled || !quantityInput.disabled) {
        // If any of the elements is not disabled, disable all elements
        minusButton.disabled = true;
        plusButton.disabled = true;
        quantityInput.disabled = true;

        console.log("qty toggle disabled")
    }
}

function enableQuantityToggle(minusButton, plusButton, quantityInput) {
    // Check if any of the elements is already enabled
    if (minusButton.disabled && plusButton.disabled && quantityInput.disabled) {
        // If all elements are disabled, enable all elements
        minusButton.disabled = false;
        plusButton.disabled = false;
        quantityInput.disabled = false;

        console.log("qty toggle enabled")
    }
}

function validAttributeSelects() {
    let selects = document.querySelectorAll('.attr-select');
    let isSelected = false;
    selects.forEach(function(select) {
        // Check if the first option is selected
        if (select.selectedIndex === 0) {
            isSelected = true;
            // Get the corresponding invalid feedback div
            const invalidFeedback = select.parentElement.querySelector(".invalid-feedback");
            // Add the d-block class to show the invalid feedback with transition
            invalidFeedback.classList.add('show');
            // Remove the d-block class after 10 seconds if not already removed
            setTimeout(function() {
                invalidFeedback.classList.remove('show');
            }, 10000); // 10 seconds in milliseconds
        }
    });
    return !isSelected;
}

function showStackedAlert(evt) {
    const message = "Item added to cart successfully!";
    const pv = evt.detail;
    // Create alert element
    let alertElement = document.createElement('div');
    alertElement.className = "alert alert-cart alert-dismissible fade show position-fixed start-50 translate-middle-x mb-0 shadow rounded"
    alertElement.setAttribute('role', 'alert');

    // Add alert content
    let currency;
    if (pv.currency_code === "USD") {
        currency = "$"
    }

    let variantAttributesHtml = '';
    if (pv.variant_attributes) {
        if (typeof pv.variant_attributes === 'object') {
            variantAttributesHtml = Object.entries(pv.variant_attributes)
              .map(([key, value]) => `<p class="text-slightly-muted mb-2" style="font-size: 0.9em;">${key}: ${value}</p>`)
              .join('');
        } else {
            variantAttributesHtml = `<p class="text-slightly-muted mb-2" style="font-size: 0.9em;">${pv.variant_attributes}</p>`;
        }
    }


    alertElement.innerHTML = `
      <p class="text-center" style="color: black;">${message}</p>
      <div class="d-flex align-items-start">
        <div class="me-3 flex-shrink-0">
          <img src="${pv.thumbnail}" class="cart-item-thumb img-thumbnail">
        </div>
        <div>
          <p class="text-slightly-muted mb-2" style="font-size: 0.9em;">${pv.product_name}</p>
          ${variantAttributesHtml}
          <p class="fs-small text-slightly-muted mb-0">${currency}${pv.price}</p>
          <p class="fs-small text-slightly-muted mb-0">Qty: ${pv.qty}</p>
        </div>
      </div>
      <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;

    // Append alert to body
    document.body.appendChild(alertElement);

    // Automatically dismiss after 5 seconds (adjust as needed)
    setTimeout(function() {
        alertElement.remove();
    }, 5000);
}

function validProductQty() {
    let quantitySpan = document.getElementById('quantity-in-stock');
    let quantityInput = document.querySelector('input[name="quantity"]');
    // Check if the quantity span doesn't include "None" and the quantity input value is not zero
    return !quantitySpan.textContent.includes("None") && quantityInput.value !== "0";
}

function getProductVariantId() {
    return parseInt(document.getElementById('product-variant-matrix').getAttribute('data-product-variant-id'));
}

function getQty() {
    return parseInt(document.querySelector('input[name="quantity"]').value)
}

function getProductVariantName() {
    let selects = document.querySelectorAll('.attr-select');
    let productName = document.getElementById("product-name").textContent;
    let selectedOptions = [];

    selects.forEach(function(elem) {
        let val = elem.options[elem.selectedIndex].value;
        selectedOptions.push(val);
    })

    let productVariantName = productName;
    if (selectedOptions.length > 0){
        productVariantName += ' - ' + selectedOptions.join(' - ');
    }

    return productVariantName;

}

function validateQuantity(evt) {
    let target = evt.target;
    let value = parseInt(target.value, 10);

    if (isNaN(value) || value < min) {
        target.value = min;
    } else if (value > max) {
        target.value = max;
    }
}

function validaQuantityOnKeydown(evt) {
    if (evt.key === 'Enter') {
        evt.preventDefault();
        validateQuantity();
        evt.target.blur();
    }
}

function handleMinusBtn(evt) {
    let quantityInput = document.querySelector('input[name="quantity"]');
    let value = parseInt(quantityInput.value, 10);
    let min = parseInt(quantityInput.getAttribute('min'), 10);
    if (!isNaN(value) && value > min) {
        quantityInput.value = value - 1;
    } else {
        quantityInput.value = min;
    }
}

function handlePlusBtn(evt) {
    let quantityInput = document.querySelector('input[name="quantity"]');
    let value = parseInt(quantityInput.value, 10);
    let max = parseInt(quantityInput.getAttribute('max'), 10);
    if (!isNaN(value) && value < max) {
        quantityInput.value = value + 1;
    } else {
        quantityInput.value = max;
    }
}


window.addEventListener('DOMContentLoaded', (evt) => {
    if (window.location.pathname.includes("/products/")) {
        loadImageSlider(evt, "image slider loaded after DOM content loaded");

        document.addEventListener("anor:incrementCartItemCount", incrementCartItemCount);
        console.log("anor:incrementCartItemCount listener added after DOM content loaded");

        document.addEventListener("anor:showCartItemAlert", showStackedAlert);
        console.log("anor:showCartItemAlert listener added after DOM content loaded");

        setupProductVariantSelectOptions();
        console.log("product variant select options setup after DOM content loaded");

        let quantityInput = document.querySelector('input[name="quantity"]');
        quantityInput.addEventListener('blur', validateQuantity);
        quantityInput.addEventListener('keydown', validaQuantityOnKeydown);

        let minusBtn = document.querySelector('.button-minus');
        let plusBtn = document.querySelector('.button-plus');
        minusBtn.addEventListener("click", handleMinusBtn);
        plusBtn.addEventListener("click", handlePlusBtn);
    }
});

document.addEventListener("htmx:afterSettle", (evt) => {
    if (evt.detail.requestConfig.path.includes("/products/")) {
        loadImageSlider(evt, "image slider loaded after after htmx settled");

        document.addEventListener("anor:incrementCartItemCount", incrementCartItemCount);
        console.log("anor:incrementCartItemCount listener loaded after after htmx settled");

        document.addEventListener("anor:showCartItemAlert", showStackedAlert);
        console.log("anor:showCartItemAlert listener loaded after after htmx settled");

        setupProductVariantSelectOptions();
        console.log("product variant select options setup after after htmx settled");

        let quantityInput = document.querySelector('input[name="quantity"]');
        quantityInput.addEventListener('blur', validateQuantity);
        quantityInput.addEventListener('keydown', validaQuantityOnKeydown);

        let minusBtn = document.querySelector('.button-minus');
        let plusBtn = document.querySelector('.button-plus');
        minusBtn.addEventListener("click", handleMinusBtn);
        plusBtn.addEventListener("click", handlePlusBtn);
    }
});

document.addEventListener("anor:showCartItemAlert", showStackedAlert);

document.addEventListener("anor:incrementCartItemCount", incrementCartItemCount)

document.addEventListener("anor:decrementCartItemCount", decrementCartItemCount)

document.addEventListener("anor:refreshCart", refreshCart);

async function refreshCart() {
    // htmx.trigger("#cart" , "refreshCart");
    // console.log("refreshCart fired");
}

function decrementCartItemCount() {
    const cartIcon = document.querySelector('.cart-icon');
    if (cartIcon) {
        const iElement = cartIcon.querySelector('i.bi-cart3');
        let notifySpan = cartIcon.querySelector('span.notify');

        if (notifySpan) {
            // Increment the number inside the notify span
            const currentCount = parseInt(notifySpan.textContent, 10);
            if (currentCount === 1) {
                notifySpan.remove();
                addNoItemsMessage();
            } else {
                notifySpan.textContent = (currentCount - 1 ).toString();
            }
        }
    }
}

function addNoItemsMessage() {
    const cartContainer = document.getElementById('cart');
    const h4Element = cartContainer.querySelector('h4.card-title');

    if (h4Element) {
        // Create the <p> element
        const pElement = document.createElement('p');
        pElement.className = 'text-center top-50 start-50 fs-3';
        pElement.textContent = 'No items in your cart';

        // Insert the <p> element after the <h4> element
        h4Element.insertAdjacentElement('afterend', pElement);
    }
}


function incrementCartItemCount() {
    console.log("update cart item count")
    const cartIcon = document.querySelector('.cart-icon');
    if (cartIcon) {
        const iElement = cartIcon.querySelector('i.bi-cart3');
        let notifySpan = cartIcon.querySelector('span.notify');

        if (notifySpan) {
            // Increment the number inside the notify span
            const currentCount = parseInt(notifySpan.textContent, 10);
            notifySpan.textContent = (currentCount + 1).toString();
        } else {
            // Create a new notify span and set its value to 1
            notifySpan = document.createElement('span');
            notifySpan.className = 'notify';
            notifySpan.textContent = '1';

            // Insert the notify span after the svg element
            iElement.parentNode.insertBefore(notifySpan, iElement.nextSibling);
        }
    }
}

window.addEventListener('htmx:historyRestore', (evt) => {
    if (evt.detail.path.includes('/products/')) {
        loadImageSlider(evt, "image slider loaded after history restored");

        document.addEventListener("anor:incrementCartItemCount", incrementCartItemCount);
        console.log("anor:incrementCartItemCount listener added after history restored");

        document.addEventListener("anor:showCartItemAlert", showStackedAlert);
        console.log("anor:showCartItemAlert listener added after history restored");

        setupProductVariantSelectOptions();
        console.log("product variant select options setup after history restored");

        let quantityInput = document.querySelector('input[name="quantity"]');
        quantityInput.addEventListener('blur', validateQuantity);
        quantityInput.addEventListener('keydown', validaQuantityOnKeydown);

        let minusBtn = document.querySelector('.button-minus');
        let plusBtn = document.querySelector('.button-plus');
        minusBtn.addEventListener("click", handleMinusBtn);
        plusBtn.addEventListener("click", handlePlusBtn);
    }
});

window.addEventListener('htmx:beforeHistorySave', (evt) => {
    if (evt.detail.path.includes('/products/')) {
        if (window.pdpMainImageSlider ){ window.pdpMainImageSlider.destroy() };
        if (window.pdpVerticalImageSlider) { window.pdpVerticalImageSlider.destroy() };
        console.log("image sliders destroyed before saving history");

        document.body.removeEventListener("anor:incrementCartItemCount", incrementCartItemCount);
        console.log("anor:incrementCartItemCount listener removed before saving history");

        document.body.removeEventListener("anor:showCartItemAlert", showStackedAlert);
        console.log("anor:showCartItemAlert listener removed before saving history");

        let quantityInput = document.querySelector('input[name="quantity"]');
        quantityInput.removeEventListener('blur', validateQuantity);
        quantityInput.removeEventListener('keydown', validaQuantityOnKeydown);

        let minusBtn = document.querySelector('.button-minus');
        let plusBtn = document.querySelector('.button-plus');
        minusBtn.removeEventListener("click", handleMinusBtn);
        plusBtn.removeEventListener("click", handlePlusBtn);
    }
});

function loadImageSlider(evt, msg) {
    const mainSlider = document.getElementById("main-slider");
    if (mainSlider){
        console.log(mainSlider)
        window.pdpMainImageSlider =  new Splide("#main-slider", {
            type: 'loop',
            height: 580,
            cover: true,
            pagination: false,
            rewind: true,
            breakpoints: {
                767: {
                    height: 290,
                },
            }
        });

        window.pdpVerticalImageSlider = new Splide("#vertical-slider", {
            height: 580,
            direction: "ttb",
            isNavigation: true,
            // fixedHeight: 90,
            autoHeight: true,
            gap        : 10,
            rewind     : true,
            pagination : false,
            arrows: false,
            wheel: true,
        });
        window.pdpMainImageSlider .sync(window.pdpVerticalImageSlider).mount();
        window.pdpVerticalImageSlider.mount();
        console.log(msg);
    }
}

function showFullAddress(id) {
    var address = document.getElementById('address' + id);
    var btn = document.getElementById('viewFullBtn' + id);

    // Simply add the 'full' class to show the address as-is
    address.classList.add('full');

    // Remove the button
    btn.remove();
}


