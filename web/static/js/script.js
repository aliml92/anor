
/////// Enable tooltip of Bootstrap5
var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
  return new bootstrap.Tooltip(tooltipTriggerEl)
})

/////// Prevent closing from click inside dropdown
document.querySelectorAll('.dropdown-menu').forEach(function(element){
    element.addEventListener('click', function (e) {
      e.stopPropagation();
    });
});
// end querySelectorAll

document.querySelectorAll('.side-filter-collapsible').forEach(function(element) {
    element.addEventListener('show.bs.collapse', function(e) {
        var siblingListItem = e.currentTarget.nextElementSibling;
        var icon = siblingListItem.querySelector('i');
        var buttonText = siblingListItem.querySelector('span');
        icon.classList.remove('fa-chevron-down');
        icon.classList.add('fa-chevron-up');
        buttonText.textContent = ' Less';
    });

    element.addEventListener('hide.bs.collapse', function(e) {
        var siblingListItem = e.currentTarget.nextElementSibling;
        var icon = siblingListItem.querySelector('i');
        var buttonText = siblingListItem.querySelector('span');
        icon.classList.remove('fa-chevron-up');
        icon.classList.add('fa-chevron-down');
        buttonText.textContent = ' More';
    });
});

var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
})
