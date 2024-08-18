function togglePasswordVisibility(passwordId, showEyeId, hideEyeId) {
  const passwordInput = document.getElementById(passwordId);
  const showEye = document.getElementById(showEyeId);
  const hideEye = document.getElementById(hideEyeId);

  hideEye.classList.remove("d-none");

  if (passwordInput.type === "password") {
    passwordInput.type = "text";
    showEye.style.display = "none";
    hideEye.style.display = "block";
  } else {
    passwordInput.type = "password";
    showEye.style.display = "block";
    hideEye.style.display = "none";
  }
}

function passwordShowHide() {
  togglePasswordVisibility("password", "password-show-eye", "password-hide-eye");
}

function confirmPasswordShowHide() {
  togglePasswordVisibility("confirm-password", "confirm-password-show-eye", "confirm-password-hide-eye");
}

window.addEventListener('htmx:afterSettle', function(evt) {
    showAlert();
});

function showAlert() {
  let errMsgEl = document.getElementById("alert-msg");
  if (errMsgEl.children.length > 0 || errMsgEl.textContent !== '') {
    let errWrapperEl = errMsgEl.parentElement;
    errWrapperEl.classList.remove("invisible");

    // Fade in
    errWrapperEl.style.opacity = '0';
    let opacity = 0;
    const fadeInInterval = setInterval(function() {
      opacity += 0.2; // Increase opacity faster
      errWrapperEl.style.opacity = opacity;
      if (opacity >= 1) {
        clearInterval(fadeInInterval);
      }
    }, 40); // Decrease interval for faster fade in

    const alertText = errMsgEl.textContent.trim();
    const wordCount = alertText.split(/\s+/).length;
    const displayDuration = Math.min(Math.max(wordCount * 200, 5000), 10000);


    setTimeout(function() {
      // Fade out
      let opacity = 1;
      const fadeOutInterval = setInterval(function() {
        opacity -= 0.2; // Decrease opacity faster
        errWrapperEl.style.opacity = opacity;
        if (opacity <= 0) {
          clearInterval(fadeOutInterval);
          errWrapperEl.classList.add("invisible");
        }
      }, 40); // Decrease interval for faster fade out

      // Remove child elements after fading out
      setTimeout(function() {
        while (errMsgEl.firstChild) {
          errMsgEl.removeChild(errMsgEl.firstChild);
        }
      }, 500); // Adjust timing as needed
    }, displayDuration);
  }
}

// getResetToken
function getResetToken() {
  const searchParams = new URLSearchParams(window.location.search);
  if (searchParams.has('token')) {
    return searchParams.get('token')
  }

  return ''
}

document.addEventListener("htmx:configRequest", function(configEvent){
  console.log("auth: before:", configEvent.detail.parameters);
  let filteredParameters = {};

  Object.entries(configEvent.detail.parameters).forEach(([key, value]) => {
    if (value !== null) {
      filteredParameters[key] = value;
    }
  });

  configEvent.detail.parameters = filteredParameters;
  console.log("auth: after:", configEvent.detail.parameters);
})