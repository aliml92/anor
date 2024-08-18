// document.addEventListener('DOMContentLoaded', function() {
//   const editBillingAddressIcon = document.getElementById('edit-billing-address');
//   const billingAddressItem = document.querySelector('.billing-address-item');
//   const billingAddressOption = document.querySelector('.billing-address-option');
//   const editBillingAddressSection = document.getElementById('edit-billing-address-section');
//   const addressNextBtn = document.getElementById('address-next-btn');
//   const addNewBillingRadio = document.getElementById('addNewBilling');
//   const newBillingForm = document.getElementById('newBillingForm');
//   const billingCollapseIcon = document.querySelector('.billing-collapse-icon');
//
//   editBillingAddressIcon.addEventListener('click', function(event) {
//     event.preventDefault();
//     addressNextBtn.style.display = 'none';
//     editBillingAddressIcon.style.display = 'none';
//
//     if (billingAddressItem && billingAddressOption) {
//       billingAddressOption.innerHTML = billingAddressItem.innerHTML
//       billingAddressItem.style.display = 'none';
//     }
//
//     editBillingAddressSection.style.display = 'block';
//
//   });
//
//   addNewBillingRadio.addEventListener('change', function() {
//     if (this.checked) {
//       newBillingForm.style.display = 'block';
//       billingCollapseIcon.style.display = 'inline-block';
//     } else {
//       newBillingForm.style.display = 'none';
//       billingCollapseIcon.style.display = 'none';
//     }
//   });
//
//   billingCollapseIcon.addEventListener('click', function(event) {
//     event.preventDefault();
//     event.stopPropagation();
//     newBillingForm.style.display = 'none';
//     billingCollapseIcon.style.display = 'none';
//     addNewBillingRadio.checked = false;
//   });
//
//
// })
//
// // handle add or edit shipping address
// document.addEventListener('DOMContentLoaded', function() {
//   const addNewShippingRadio = document.getElementById('addNewShipping');
//   const newShippingForm = document.getElementById('newShippingForm');
//   const editShippingForm = document.getElementById('editShippingForm');
//   const shippingCollapseIcon = document.querySelector('.shipping-collapse-icon');
//   const editShippingIcon = document.querySelector('.edit-shipping-icon');
//   const defaultShippingCollapseIcon = document.querySelector('.default-shipping-collapse-icon');
//   const shippingDisplay = document.getElementById('shippingDisplay');
//   const editShippingText = document.getElementById('editShippingText');
//
//   addNewShippingRadio.addEventListener('change', function() {
//     if (this.checked) {
//       newShippingForm.style.display = 'block';
//       shippingCollapseIcon.style.display = 'inline-block';
//       editShippingForm.style.display = 'none';
//       resetDefaultAddressDisplay();
//     } else {
//       newShippingForm.style.display = 'none';
//       shippingCollapseIcon.style.display = 'none';
//     }
//   });
//
//   shippingCollapseIcon.addEventListener('click', function(event) {
//     event.preventDefault();
//     event.stopPropagation();
//     newShippingForm.style.display = 'none';
//     shippingCollapseIcon.style.display = 'none';
//     addNewShippingRadio.checked = false;
//     document.getElementById('defaultShipping').checked = true;
//   });
//
//   editShippingIcon.addEventListener('click', function(event) {
//     event.preventDefault();
//     editShippingForm.style.display = 'block';
//     newShippingForm.style.display = 'none';
//     shippingCollapseIcon.style.display = 'none';
//     addNewShippingRadio.checked = false;
//     document.getElementById('defaultShipping').checked = true;
//
//     // Show "Edit address" text and hide address details
//     shippingDisplay.style.display = 'none';
//     editShippingText.style.display = 'inline';
//     editShippingIcon.style.display = 'none';
//     defaultShippingCollapseIcon.style.display = 'inline-block';
//
//     // Populate edit form with current address data
//     document.getElementById('editFullName').value = document.getElementById('addressName').textContent;
//     document.getElementById('editAddressLine1').value = document.getElementById('addressLine1').textContent;
//     document.getElementById('editAddressLine2').value = document.getElementById('addressLine2').textContent;
//     document.getElementById('editCity').value = document.getElementById('addressCity').textContent;
//     document.getElementById('editState').value = document.getElementById('addressState').textContent;
//     document.getElementById('editZipCode').value = document.getElementById('addressZip').textContent;
//     document.getElementById('editCountry').value = document.getElementById('addressCountry').textContent;
//   });
//
//   defaultShippingCollapseIcon.addEventListener('click', function(event) {
//     event.preventDefault();
//     event.stopPropagation();
//     resetDefaultAddressDisplay();
//     editShippingForm.style.display = 'none';
//   });
//
//   document.getElementById('saveEditButton').addEventListener('click', function() {
//     // Update the address with edited data
//     document.getElementById('addressName').textContent = document.getElementById('editFullName').value;
//     document.getElementById('addressLine1').textContent = document.getElementById('editAddressLine1').value;
//     document.getElementById('addressLine2').textContent = document.getElementById('editAddressLine2').value;
//     document.getElementById('addressCity').textContent = document.getElementById('editCity').value;
//     document.getElementById('addressState').textContent = document.getElementById('editState').value;
//     document.getElementById('addressZip').textContent = document.getElementById('editZipCode').value;
//     document.getElementById('addressCountry').textContent = document.getElementById('editCountry').value;
//
//     // Hide the edit form and reset display
//     editShippingForm.style.display = 'none';
//     resetDefaultAddressDisplay();
//   });
//
//   function resetDefaultAddressDisplay() {
//     shippingDisplay.style.display = 'inline';
//     editShippingText.style.display = 'none';
//     editShippingIcon.style.display = 'inline-block';
//     defaultShippingCollapseIcon.style.display = 'none';
//   }
// });